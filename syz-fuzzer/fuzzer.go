// Copyright 2015 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/syzkaller/pkg/hash"
	"github.com/google/syzkaller/pkg/host"
	"github.com/google/syzkaller/pkg/ipc"
	"github.com/google/syzkaller/pkg/ipc/ipcconfig"
	"github.com/google/syzkaller/pkg/log"
	"github.com/google/syzkaller/pkg/osutil"
	"github.com/google/syzkaller/pkg/rpctype"
	"github.com/google/syzkaller/pkg/signal"
	"github.com/google/syzkaller/prog"
	_ "github.com/google/syzkaller/sys"

	pb "github.com/google/syzkaller/pkg/dra"
)

type Fuzzer struct {
	name        string
	outputType  OutputType
	config      *ipc.Config
	execOpts    *ipc.ExecOpts
	procs       []*Proc
	gate        *ipc.Gate
	workQueue   *WorkQueue
	needPoll    chan struct{}
	choiceTable *prog.ChoiceTable
	stats       [StatCount]uint64
	manager     *rpctype.RPCClient
	target      *prog.Target

	faultInjectionEnabled    bool
	comparisonTracingEnabled bool

	corpusMu     sync.RWMutex
	corpus       []*prog.Prog
	corpusHashes map[hash.Sig]struct{}

	signalMu     sync.RWMutex
	corpusSignal signal.Signal // signal of inputs in corpus
	maxSignal    signal.Signal // max signal ever observed including flakes
	newSignal    signal.Signal // diff of maxSignal since last sync with master

	logMu sync.Mutex

	dManager *pb.DRPCClient
	//corpusDMu        sync.RWMutex
	//corpusSig        []string
	//corpusDependency map[string]*prog.Prog

	coverMu sync.RWMutex
	cover   map[int]*pb.Call
}

type Stat int

const (
	StatGenerate Stat = iota
	StatFuzz
	StatCandidate
	StatTriage
	StatMinimize
	StatSmash
	StatHint
	StatSeed
	StatDependency
	StatCount
)

var statNames = [StatCount]string{
	StatGenerate:   "exec gen",
	StatFuzz:       "exec fuzz",
	StatCandidate:  "exec candidate",
	StatTriage:     "exec triage",
	StatMinimize:   "exec minimize",
	StatSmash:      "exec smash",
	StatHint:       "exec hints",
	StatSeed:       "exec seeds",
	StatDependency: "exec dependency",
}

type OutputType int

const (
	OutputNone OutputType = iota
	OutputStdout
	OutputDmesg
	OutputFile
)

func main() {
	debug.SetGCPercent(50)

	var (
		flagName     = flag.String("name", "test", "unique name for manager")
		flagOS       = flag.String("os", runtime.GOOS, "target OS")
		flagArch     = flag.String("arch", runtime.GOARCH, "target arch")
		flagManager  = flag.String("manager", "", "manager rpc address")
		flagDManager = flag.String("dmanager", "", "dependency manager rpc address")
		flagProcs    = flag.Int("procs", 1, "number of parallel test processes")
		flagOutput   = flag.String("output", "stdout", "write programs to none/stdout/dmesg/file")
		flagPprof    = flag.String("pprof", "", "address to serve pprof profiles")
		flagTest     = flag.Bool("test", false, "enable image testing mode")      // used by syz-ci
		flagRunTest  = flag.Bool("runtest", false, "enable program testing mode") // used by pkg/runtest

	)
	flag.Parse()
	outputType := parseOutputType(*flagOutput)
	log.Logf(0, "fuzzer started")

	target, err := prog.GetTarget(*flagOS, *flagArch)
	if err != nil {
		log.Fatalf("%v", err)
	}

	config, execOpts, err := ipcconfig.Default(target)
	if err != nil {
		log.Fatalf("failed to create default ipc config: %v", err)
	}
	sandbox := ipc.FlagsToSandbox(config.Flags)
	shutdown := make(chan struct{})
	osutil.HandleInterrupts(shutdown)
	go func() {
		// Handles graceful preemption on GCE.
		<-shutdown
		log.Logf(0, "SYZ-FUZZER: PREEMPTED")
		os.Exit(1)
	}()

	checkArgs := &checkArgs{
		target:      target,
		sandbox:     sandbox,
		ipcConfig:   config,
		ipcExecOpts: execOpts,
	}
	if *flagTest {
		testImage(*flagManager, checkArgs)
		return
	}

	if *flagPprof != "" {
		go func() {
			err := http.ListenAndServe(*flagPprof, nil)
			log.Fatalf("failed to serve pprof profiles: %v", err)
		}()
	} else {
		runtime.MemProfileRate = 0
	}

	log.Logf(0, "dialing manager at %v", *flagManager)
	manager, err := rpctype.NewRPCClient(*flagManager)
	if err != nil {
		log.Fatalf("failed to connect to manager: %v ", err)
	}
	a := &rpctype.ConnectArgs{Name: *flagName}
	r := &rpctype.ConnectRes{}
	if err := manager.Call("Manager.Connect", a, r); err != nil {
		log.Fatalf("failed to connect to manager: %v ", err)
	}
	if r.CheckResult == nil {
		checkArgs.gitRevision = r.GitRevision
		checkArgs.targetRevision = r.TargetRevision
		checkArgs.enabledCalls = r.EnabledCalls
		checkArgs.allSandboxes = r.AllSandboxes
		r.CheckResult, err = checkMachine(checkArgs)
		if err != nil {
			if r.CheckResult == nil {
				r.CheckResult = new(rpctype.CheckArgs)
			}
			r.CheckResult.Error = err.Error()
		}
		r.CheckResult.Name = *flagName
		if err := manager.Call("Manager.Check", r.CheckResult, nil); err != nil {
			log.Fatalf("Manager.Check call failed: %v", err)
		}
		if r.CheckResult.Error != "" {
			log.Fatalf("%v", r.CheckResult.Error)
		}
	}
	log.Logf(0, "syscalls: %v", len(r.CheckResult.EnabledCalls[sandbox]))
	for _, feat := range r.CheckResult.Features {
		log.Logf(0, "%v: %v", feat.Name, feat.Reason)
	}
	periodicCallback, err := host.Setup(target, r.CheckResult.Features)
	if err != nil {
		log.Fatalf("BUG: %v", err)
	}
	var gateCallback func()
	if periodicCallback != nil {
		gateCallback = func() { periodicCallback(r.MemoryLeakFrames) }
	}
	if r.CheckResult.Features[host.FeatureExtraCoverage].Enabled {
		config.Flags |= ipc.FlagExtraCover
	}
	if r.CheckResult.Features[host.FeatureFaultInjection].Enabled {
		config.Flags |= ipc.FlagEnableFault
	}
	if r.CheckResult.Features[host.FeatureNetworkInjection].Enabled {
		config.Flags |= ipc.FlagEnableTun
	}
	if r.CheckResult.Features[host.FeatureNetworkDevices].Enabled {
		config.Flags |= ipc.FlagEnableNetDev
	}
	config.Flags |= ipc.FlagEnableNetReset
	config.Flags |= ipc.FlagEnableCgroups
	config.Flags |= ipc.FlagEnableBinfmtMisc
	config.Flags |= ipc.FlagEnableCloseFds

	if *flagRunTest {
		runTest(target, manager, *flagName, config.Executor)
		return
	}

	log.Logf(0, "dialing dManager at %v", *flagDManager)
	dManager := &pb.DRPCClient{
		I: []*pb.Input{},
	}
	dManager.RunDependencyRPCClient(flagDManager, flagName)

	needPoll := make(chan struct{}, 1)
	needPoll <- struct{}{}
	fuzzer := &Fuzzer{
		name:                     *flagName,
		outputType:               outputType,
		config:                   config,
		execOpts:                 execOpts,
		gate:                     ipc.NewGate(2**flagProcs, gateCallback),
		workQueue:                newWorkQueue(*flagProcs, needPoll),
		needPoll:                 needPoll,
		manager:                  manager,
		target:                   target,
		faultInjectionEnabled:    r.CheckResult.Features[host.FeatureFaultInjection].Enabled,
		comparisonTracingEnabled: r.CheckResult.Features[host.FeatureComparisons].Enabled,
		corpusHashes:             make(map[hash.Sig]struct{}),
		dManager:                 dManager,
		//corpusSig:                []string{},
		//corpusDependency:         make(map[string]*prog.Prog),
		cover: make(map[int]*pb.Call),
	}
	for i := 0; fuzzer.poll(i == 0, nil); i++ {
	}
	calls := make(map[*prog.Syscall]bool)
	for _, id := range r.CheckResult.EnabledCalls[sandbox] {
		calls[target.Syscalls[id]] = true
	}
	prios := target.CalculatePriorities(fuzzer.corpus)
	fuzzer.choiceTable = target.BuildChoiceTable(prios, calls)

	for pid := 0; pid < *flagProcs; pid++ {
		proc, err := newProc(fuzzer, pid)
		if err != nil {
			log.Fatalf("failed to create proc: %v", err)
		}
		fuzzer.procs = append(fuzzer.procs, proc)
		go proc.loop()
	}

	fuzzer.pollLoop()
}

func (fuzzer *Fuzzer) pollLoop() {
	var execTotal uint64
	var lastPoll time.Time
	var lastPrint time.Time
	ticker := time.NewTicker(3 * time.Second).C
	for {

		//data := fuzzer.corpus[0].Serialize()
		//sig := hash.Hash(data)
		//fuzzer.dManager.SendDependencyInput(sig.String())
		if len(fuzzer.workQueue.dependency) == 0 {
			newDependencyInput := fuzzer.dManager.GetDependencyInput(fuzzer.name)
			for _, dependencyInput := range newDependencyInput.GetDependencyInput() {
				fuzzer.addDInputFromAnotherFuzzer(dependencyInput)
			}
		}
		fuzzer.dManager.SSendLog()

		poll := false
		select {
		case <-ticker:
		case <-fuzzer.needPoll:
			poll = true
		}
		if fuzzer.outputType != OutputStdout && time.Since(lastPrint) > 10*time.Second {
			// Keep-alive for manager.
			log.Logf(0, "alive, executed %v", execTotal)
			lastPrint = time.Now()
		}
		if poll || time.Since(lastPoll) > 10*time.Second {
			needCandidates := fuzzer.workQueue.wantCandidates()
			if poll && !needCandidates {
				continue
			}
			stats := make(map[string]uint64)
			for _, proc := range fuzzer.procs {
				stats["exec total"] += atomic.SwapUint64(&proc.env.StatExecs, 0)
				stats["executor restarts"] += atomic.SwapUint64(&proc.env.StatRestarts, 0)
			}
			for stat := Stat(0); stat < StatCount; stat++ {
				v := atomic.SwapUint64(&fuzzer.stats[stat], 0)
				stats[statNames[stat]] = v
				execTotal += v
			}
			if !fuzzer.poll(needCandidates, stats) {
				lastPoll = time.Now()
			}
		}

	}
}

func (fuzzer *Fuzzer) poll(needCandidates bool, stats map[string]uint64) bool {
	a := &rpctype.PollArgs{
		Name:           fuzzer.name,
		NeedCandidates: needCandidates,
		MaxSignal:      fuzzer.grabNewSignal().Serialize(),
		Stats:          stats,
	}
	r := &rpctype.PollRes{}
	if err := fuzzer.manager.Call("Manager.Poll", a, r); err != nil {
		log.Fatalf("Manager.Poll call failed: %v", err)
	}
	maxSignal := r.MaxSignal.Deserialize()
	log.Logf(1, "poll: candidates=%v inputs=%v signal=%v",
		len(r.Candidates), len(r.NewInputs), maxSignal.Len())
	fuzzer.addMaxSignal(maxSignal)
	for _, inp := range r.NewInputs {
		fuzzer.addInputFromAnotherFuzzer(inp)
	}
	for _, candidate := range r.Candidates {
		p, err := fuzzer.target.Deserialize(candidate.Prog, prog.NonStrict)
		if err != nil {
			log.Fatalf("failed to parse program from manager: %v", err)
		}
		flags := ProgCandidate
		if candidate.Minimized {
			flags |= ProgMinimized
		}
		if candidate.Smashed {
			flags |= ProgSmashed
		}
		fuzzer.workQueue.enqueue(&WorkCandidate{
			p:     p,
			flags: flags,
		})
	}
	return len(r.NewInputs) != 0 || len(r.Candidates) != 0 || maxSignal.Len() != 0
}

func (fuzzer *Fuzzer) sendInputToManager(inp rpctype.RPCInput) {
	a := &rpctype.NewInputArgs{
		Name:     fuzzer.name,
		RPCInput: inp,
	}
	if err := fuzzer.manager.Call("Manager.NewInput", a, nil); err != nil {
		log.Fatalf("Manager.NewInput call failed: %v", err)
	}
}

func (fuzzer *Fuzzer) addInputFromAnotherFuzzer(inp rpctype.RPCInput) {
	p, err := fuzzer.target.Deserialize(inp.Prog, prog.NonStrict)
	if err != nil {
		log.Fatalf("failed to deserialize prog from another fuzzer: %v", err)
	}
	sig := hash.Hash(inp.Prog)
	sign := inp.Signal.Deserialize()
	fuzzer.addInputToCorpus(p, sign, sig)
}

func (fuzzer *Fuzzer) addInputToCorpus(p *prog.Prog, sign signal.Signal, sig hash.Sig) {
	fuzzer.corpusMu.Lock()
	if _, ok := fuzzer.corpusHashes[sig]; !ok {
		fuzzer.corpus = append(fuzzer.corpus, p)
		fuzzer.corpusHashes[sig] = struct{}{}
	}
	fuzzer.corpusMu.Unlock()

	if !sign.Empty() {
		fuzzer.signalMu.Lock()
		fuzzer.corpusSignal.Merge(sign)
		fuzzer.maxSignal.Merge(sign)
		fuzzer.signalMu.Unlock()
	}
}

func (fuzzer *Fuzzer) addDInputFromAnotherFuzzer(dependencyInput *pb.DependencyInput) {
	log.Logf(1, "dependencyInput : %v", dependencyInput)
	//fuzzer.dManager.SendLog(fmt.Sprintf("dependencyInput : %v", dependencyInput))
	//sig := dependencyInput.GetSig()
	p, err := fuzzer.target.Deserialize(dependencyInput.GetProg(), prog.NonStrict)
	p.Uncover = make(map[int]*prog.Uncover)
	if err != nil {
		log.Fatalf("failed to deserialize prog from another fuzzer: %v", err)
	}

	for idx, u := range dependencyInput.GetUncoveredAddress() {
		u1 := new(prog.Uncover)
		u1.UncoveredAddress = u.GetAddress()
		u1.Idx = u.GetIdx()
		for _, a := range u.GetRelatedAddress() {

			a1 := &prog.RelatedAddresses{
				RelatedAddress: a.GetAddress(),
				Prio:           a.GetPrio(),
				Repeat:         a.GetRepeat(),
			}

			for _, i := range a.GetRelatedInput() {
				rp, err := fuzzer.target.Deserialize(i.GetProg(), prog.NonStrict)
				if err != nil {
					panic(err)
				}
				a1.RelatedProgs = append(a1.RelatedProgs, rp)
			}

			for _, i := range a.GetRelatedSyscall() {
				c1 := &prog.Call{
					Meta:    nil,
					Ret:     nil,
					Comment: "dependency",
				}

				log.Logf(1, "cmd value : %x", i.Number)
				//fuzzer.dManager.SendLog(fmt.Sprintf("cmd value : %x", i.Number))
				// only work for ioctl
				for n, c := range fuzzer.target.SyscallMap {
					if strings.HasPrefix(n, i.Name) {
						for _, a := range c.Args {
							if a.FieldName() == "cmd" {
								switch t := a.DefaultArg().(type) {
								case *prog.ConstArg:
									val, _ := t.Value()
									if val == i.Number {
										c1.Meta = c
										log.Logf(1, "ioctl name : %v", c.Name)
										//fuzzer.dManager.SendLog(fmt.Sprintf("ioctl name : %v", c.Name))
										c1.Ret = prog.MakeReturnArg(c.Ret)
										for _, typ := range c.Args {
											arg := typ.DefaultArg()
											c1.Args = append(c1.Args, arg)
										}
										a1.RelatedCalls = append(a1.RelatedCalls, c1)
									}
								default:

								}
							}
						}
					}
				}
			}
			u1.RelatedAddress = append(u1.RelatedAddress, a1)
		}
		p.Uncover[idx] = u1
	}

	fuzzer.workQueue.enqueue(&WorkDependency{
		p: p.CloneWithUncover(),
	})

	for _, u := range p.Uncover {
		log.Logf(1, "fuzzer.addDInputFromAnotherFuzzer Uncover: %v", u)
		//fuzzer.dManager.SendLog(fmt.Sprintf("fuzzer.addDInputFromAnotherFuzzer Uncover: %v", u))
	}

}

//func (fuzzer *Fuzzer) corpusSigSnapshot() []string {
//	fuzzer.corpusDMu.RLock()
//	defer fuzzer.corpusDMu.RUnlock()
//	return fuzzer.corpusSig
//}
//
//func (fuzzer *Fuzzer) corpusDependencySnapshot() map[string]*prog.Prog {
//	fuzzer.corpusDMu.RLock()
//	defer fuzzer.corpusDMu.RUnlock()
//	return fuzzer.corpusDependency
//}

func (fuzzer *Fuzzer) corpusSnapshot() []*prog.Prog {
	fuzzer.corpusMu.RLock()
	defer fuzzer.corpusMu.RUnlock()
	return fuzzer.corpus
}

func (fuzzer *Fuzzer) addMaxSignal(sign signal.Signal) {
	if sign.Len() == 0 {
		return
	}
	fuzzer.signalMu.Lock()
	defer fuzzer.signalMu.Unlock()
	fuzzer.maxSignal.Merge(sign)
}

func (fuzzer *Fuzzer) grabNewSignal() signal.Signal {
	fuzzer.signalMu.Lock()
	defer fuzzer.signalMu.Unlock()
	sign := fuzzer.newSignal
	if sign.Empty() {
		return nil
	}
	fuzzer.newSignal = nil
	return sign
}

func (fuzzer *Fuzzer) corpusSignalDiff(sign signal.Signal) signal.Signal {
	fuzzer.signalMu.RLock()
	defer fuzzer.signalMu.RUnlock()
	return fuzzer.corpusSignal.Diff(sign)
}

func (fuzzer *Fuzzer) checkNewSignal(p *prog.Prog, info *ipc.ProgInfo) (calls []int, extra bool) {
	fuzzer.signalMu.RLock()
	defer fuzzer.signalMu.RUnlock()
	for i, inf := range info.Calls {
		if fuzzer.checkNewCallSignal(p, &inf, i) {
			calls = append(calls, i)
		}
	}
	extra = fuzzer.checkNewCallSignal(p, &info.Extra, -1)
	return
}

func (fuzzer *Fuzzer) checkNewCallSignal(p *prog.Prog, info *ipc.CallInfo, call int) bool {
	diff := fuzzer.maxSignal.DiffRaw(info.Signal, signalPrio(p, info, call))
	if diff.Empty() {
		return false
	}
	fuzzer.signalMu.RUnlock()
	fuzzer.signalMu.Lock()
	fuzzer.maxSignal.Merge(diff)
	fuzzer.newSignal.Merge(diff)
	fuzzer.signalMu.Unlock()
	fuzzer.signalMu.RLock()
	return true
}

func signalPrio(p *prog.Prog, info *ipc.CallInfo, call int) (prio uint8) {
	if call == -1 {
		return 0
	}
	if info.Errno == 0 {
		prio |= 1 << 1
	}
	if !p.Target.CallContainsAny(p.Calls[call]) {
		prio |= 1 << 0
	}
	return
}

func parseOutputType(str string) OutputType {
	switch str {
	case "none":
		return OutputNone
	case "stdout":
		return OutputStdout
	case "dmesg":
		return OutputDmesg
	case "file":
		return OutputFile
	default:
		log.Fatalf("-output flag must be one of none/stdout/dmesg/file")
		return OutputNone
	}
}

func (fuzzer *Fuzzer) checkNewCoverage(p *prog.Prog, info *ipc.ProgInfo) (calls []int) {
	fuzzer.coverMu.Lock()

	input := &pb.Input{
		Call: make(map[uint32]*pb.Call),
	}
	data := p.Serialize()
	sig := hash.Hash(data)
	input.Sig = sig.String()
	tflags := false
	for i, inf := range info.Calls {
		input.Call[uint32(i)] = &pb.Call{
			Idx:     uint32(i),
			Address: map[uint32]uint32{},
		}
		newCall := input.Call[uint32(i)]

		id := p.Calls[i].Meta.ID
		if _, ok := fuzzer.cover[id]; !ok {
			fuzzer.cover[id] = &pb.Call{
				Idx:     0,
				Address: make(map[uint32]uint32),
			}
		}
		call := fuzzer.cover[id].Address
		flags := false
		for _, address := range inf.Cover {
			if _, ok := call[address]; !ok {
				call[address] = 0
				flags = true
				newCall.Address[address] = 0
			}
		}
		if flags == true {
			calls = append(calls, i)
			tflags = true
		}
	}

	if tflags {
		//fuzzer.dManager.SendInput(input)
	}

	//for _, cc := range info.Calls {
	//	log.Logf(1, "Dependency gRPC checkNewCoverage address : %v", cc.Cover)
	//}

	fuzzer.coverMu.Unlock()
	return
}

func (fuzzer *Fuzzer) checkIsCovered(id int, address uint32) (res bool) {
	fuzzer.coverMu.RLock()
	call := fuzzer.cover[id].Address
	if _, ok := call[address]; !ok {
		return false
	} else {
		return true
	}
}
