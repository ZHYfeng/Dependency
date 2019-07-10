// Copyright 2017 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/syzkaller/pkg/cover"
	pb "github.com/google/syzkaller/pkg/dra"
	"github.com/google/syzkaller/pkg/hash"
	"github.com/google/syzkaller/pkg/ipc"
	"github.com/google/syzkaller/pkg/log"
	"github.com/google/syzkaller/pkg/rpctype"
	"github.com/google/syzkaller/pkg/signal"
	"github.com/google/syzkaller/prog"
)

const (
	programLength = 30
)

// Proc represents a single fuzzing process (executor).
type Proc struct {
	fuzzer            *Fuzzer
	pid               int
	env               *ipc.Env
	rnd               *rand.Rand
	execOpts          *ipc.ExecOpts
	execOptsCover     *ipc.ExecOpts
	execOptsComps     *ipc.ExecOpts
	execOptsNoCollide *ipc.ExecOpts
}

func newProc(fuzzer *Fuzzer, pid int) (*Proc, error) {
	env, err := ipc.MakeEnv(fuzzer.config, pid)
	if err != nil {
		return nil, err
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(pid)*1e12))
	execOptsNoCollide := *fuzzer.execOpts
	execOptsNoCollide.Flags &= ^ipc.FlagCollide
	execOptsCover := execOptsNoCollide
	execOptsCover.Flags |= ipc.FlagCollectCover
	execOptsComps := execOptsNoCollide
	execOptsComps.Flags |= ipc.FlagCollectComps
	proc := &Proc{
		fuzzer:            fuzzer,
		pid:               pid,
		env:               env,
		rnd:               rnd,
		execOpts:          fuzzer.execOpts,
		execOptsCover:     &execOptsCover,
		execOptsComps:     &execOptsComps,
		execOptsNoCollide: &execOptsNoCollide,
	}
	return proc, nil
}

func (proc *Proc) loop() {
	generatePeriod := 100
	if proc.fuzzer.config.Flags&ipc.FlagSignal == 0 {
		// If we don't have real coverage signal, generate programs more frequently
		// because fallback signal is weak.
		generatePeriod = 2
	}
	for i := 0; ; i++ {
		log.Logf(1, "loop : %v", i)
		item := proc.fuzzer.workQueue.dequeue()
		log.Logf(1, "item := proc.fuzzer.workQueue.dequeue()")
		if item != nil {
			switch item := item.(type) {
			case *WorkTriage:
				proc.triageInput(item)
			case *WorkCandidate:
				proc.execute(proc.execOpts, item.p, item.flags, StatCandidate)
			case *WorkSmash:
				proc.smashInput(item)
			case *WorkDependency:
				proc.dependencyMutate(item)
			default:
				log.Fatalf("unknown work type: %#v", item)
			}
			continue
		}

		ct := proc.fuzzer.choiceTable
		corpus := proc.fuzzer.corpusSnapshot()
		if len(corpus) == 0 || i%generatePeriod == 0 {
			// Generate a new prog.
			p := proc.fuzzer.target.Generate(proc.rnd, programLength, ct)
			log.Logf(1, "#%v: generated", proc.pid)
			proc.execute(proc.execOpts, p, ProgNormal, StatGenerate)
		} else {
			// Mutate an existing prog.
			log.Logf(1, "#%v: mutated", proc.pid)
			p := corpus[proc.rnd.Intn(len(corpus))].Clone()
			p.Mutate(proc.rnd, programLength, ct, corpus)
			proc.execute(proc.execOpts, p, ProgNormal, StatFuzz)
			//info := proc.execute(proc.execOpts, p, ProgNormal, StatFuzz)
			//proc.fuzzer.checkNewCoverage(p, info)
		}
	}
}

func (proc *Proc) triageInput(item *WorkTriage) {

	log.Logf(1, "#%v: triaging type=%x", proc.pid, item.flags)
	prio := signalPrio(item.p, &item.info, item.call)
	inputSignal := signal.FromRaw(item.info.Signal, prio)
	newSignal := proc.fuzzer.corpusSignalDiff(inputSignal)
	if newSignal.Empty() {
		return
	}
	callName := ".extra"
	logCallName := "extra"
	if item.call != -1 {
		callName = item.p.Calls[item.call].Meta.CallName
		logCallName = fmt.Sprintf("call #%v %v", item.call, callName)
	}
	log.Logf(3, "triaging input for %v (new signal=%v)", logCallName, newSignal.Len())
	var inputCover cover.Cover
	const (
		signalRuns       = 3
		minimizeAttempts = 3
	)

	// Compute input coverage and non-flaky signal for minimization.
	notexecuted := 0
	for i := 0; i < signalRuns; i++ {
		log.Logf(3, "triaging input signalRuns")
		info := proc.executeRaw(proc.execOptsCover, item.p, StatTriage)
		if !reexecutionSuccess(info, &item.info, item.call) {
			// The call was not executed or failed.
			notexecuted++
			if notexecuted > signalRuns/2+1 {
				return // if happens too often, give up
			}
			continue
		}
		thisSignal, thisCover := getSignalAndCover(item.p, info, item.call)
		newSignal = newSignal.Intersection(thisSignal)
		// Without !minimized check manager starts losing some considerable amount
		// of coverage after each restart. Mechanics of this are not completely clear.
		if newSignal.Empty() && item.flags&ProgMinimized == 0 {
			return
		}
		inputCover.Merge(thisCover)

		proc.fuzzer.checkNewCoverage(item.p, info)
		//for i, c := range info.Calls {
		//	ii := uint32(i)
		//	if cc, ok := input.Call[ii]; !ok {
		//		cc = &pb.Call{
		//			Idx:     ii,
		//			Address: make(map[uint32]uint32),
		//		}
		//		input.Call[ii] = cc
		//	}
		//	cc := input.Call[ii]
		//	for _, a := range c.Cover {
		//		cc.Address[a] = 0
		//	}
		//
		//}
	}
	if item.flags&ProgMinimized == 0 {
		item.p, item.call = prog.Minimize(item.p, item.call, false,
			func(p1 *prog.Prog, call1 int) bool {
				for i := 0; i < minimizeAttempts; i++ {
					log.Logf(3, "minimizeAttempts")
					info := proc.execute(proc.execOptsNoCollide, p1, ProgNormal, StatMinimize)
					if !reexecutionSuccess(info, &item.info, call1) {
						// The call was not executed or failed.
						continue
					}
					thisSignal, _ := getSignalAndCover(p1, info, call1)
					if newSignal.Intersection(thisSignal).Len() == newSignal.Len() {
						return true
					}
				}
				return false
			})
	}

	data := item.p.Serialize()
	sig := hash.Hash(data)

	log.Logf(2, "added new input sig %s for %v to corpus:\n%s", sig.String(), logCallName, data)
	proc.fuzzer.sendInputToManager(rpctype.RPCInput{
		Call:   callName,
		Prog:   data,
		Signal: inputSignal.Serialize(),
		Cover:  inputCover.Serialize(),
	})

	proc.fuzzer.addInputToCorpus(item.p, inputSignal, sig)

	if item.flags&ProgSmashed == 0 {
		proc.fuzzer.workQueue.enqueue(&WorkSmash{item.p, item.call})
	}

	input := pb.Input{
		Call:       make(map[uint32]*pb.Call),
		Dependency: false,
	}
	if item.call != -1 {
		cc := &pb.Call{
			Idx:     uint32(item.call),
			Address: make(map[uint32]uint32),
		}
		input.Call[uint32(item.call)] = cc
		for a, _ := range inputCover {
			cc.Address[a] = 0
		}
	}
	for _, c := range item.p.Comments {
		if c == "StatDependency" {
			input.Dependency = true
			proc.fuzzer.dManager.SendLog(fmt.Sprintf("real new input from StatDependency : %v", item.p))
		}
	}

	input.Sig = sig.String()
	for _, c := range data {
		input.Program = append(input.Program, c)
	}
	proc.fuzzer.dManager.SendInput(&input)
}

func (proc *Proc) getCall(sc *pb.Syscall) (res *prog.Syscall) {
	// only work for ioctl
	for n, c := range proc.fuzzer.target.SyscallMap {
		if strings.HasPrefix(n, sc.Name) {
			for _, a := range c.Args {
				if a.FieldName() == "cmd" {
					switch t := a.DefaultArg().(type) {
					case *prog.ConstArg:
						val, _ := t.Value()
						if val == sc.Cmd {
							res = c
							return
						}
					default:
					}
				}
			}
		}
	}
	return
}

func (proc *Proc) dependencyMutate(item *WorkDependency) (result bool) {

	log.Logf(1, "#%v: DependencyMutate", proc.pid)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("#%v: DependencyMutate", proc.pid))

	//corpus := proc.fuzzer.corpusSnapshot()
	//corpusSigSnapshot := proc.fuzzer.corpusSigSnapshot()
	//log.Logf(3, "corpusSigSnapshot size : %v", len(corpusSigSnapshot))
	//corpusDependencySnapshot := proc.fuzzer.corpusDependencySnapshot()
	//log.Logf(3, "corpusDependencySnapshot size : %v", len(corpusDependencySnapshot))

	dependencyInput := item.dependencyInput
	p, err := proc.fuzzer.target.Deserialize(dependencyInput.Program, prog.NonStrict)
	if err != nil {
		log.Fatalf("failed to deserialize program from dependencyInput: %v", err)
	}
	log.Logf(1, "DependencyMutate program : \n%s", dependencyInput.Program)

	for _, u := range dependencyInput.UncoveredAddress {

		if u.RunTimeDate.CheckAddress == false {

			info := proc.execute(proc.execOptsCover, p, ProgNormal, StatDependency)

			if checkAddress(u.ConditionAddress, info.Calls[u.Idx].Cover) {
				u.RunTimeDate.CheckCondition = true
				for _, wa := range u.WriteAddress {
					proc.dependencyWriteAddress(wa)
				}
			} else {

			}
		} else {

		}

	}
	return
}

func (proc *Proc) dependencyWriteAddress(wa *pb.WriteAddress) (res bool, info *ipc.ProgInfo) {
	if wa.RunTimeDate.TaskStatus == pb.RunTimeData_untested {
		for _, wc := range wa.WriteSyscall {
			if wc.RunTimeDate.TaskStatus == pb.RunTimeData_untested {

				ct := proc.fuzzer.choiceTable
				p, err := proc.fuzzer.target.Deserialize(wc.RunTimeDate.Program, prog.NonStrict)
				if err != nil {
					log.Fatalf("failed to deserialize program from dependencyWriteAddress: %v", err)
				}

				p0 := p.Clone()
				c0c := p0.GetCall(proc.rnd, proc.getCall(wc), wc.RunTimeDate.Parent.Idx, ct)

				p0.InsertCall(c0c, wc.RunTimeDate.Parent.Idx, programLength)

				size := uint32(len(c0c))
				wc.RunTimeDate.Idx = wc.RunTimeDate.Parent.Idx + size - 1

				data := p0.Serialize()
				copy(wc.RunTimeDate.Program, data)

				var info *ipc.ProgInfo
				for i := 0; i < 100; i++ {
					p0.MutateIoctl3Arg(proc.rnd, wc.RunTimeDate.Idx, ct)
					info = proc.execute(proc.execOptsCover, p0, ProgNormal, StatDependency)

					if checkAddress(wc.RunTimeDate.Address, info.Calls[wc.RunTimeDate.Idx].Cover) {
						wc.RunTimeDate.CheckAddress = true
						if len(wc.RunTimeDate.Program) == 0 {
							data := p0.Serialize()
							copy(wc.RunTimeDate.Program, data)
						}
					}

					if checkAddress(wc.RunTimeDate.Parent.ConditionAddress, info.Calls[wc.RunTimeDate.Parent.Idx].Cover) {
						wc.RunTimeDate.Parent.CheckCondition = true
						if checkAddress(wc.RunTimeDate.Parent.Address, info.Calls[wc.RunTimeDate.Parent.Idx].Cover) {

							wc.RunTimeDate.Parent.CheckAddress = true

							wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested

							data := p0.Serialize()
							copy(wc.RunTimeDate.Program, data)

							return true, info
						}
					}

				}

				if wc.RunTimeDate.Parent.CheckCondition == true && wc.RunTimeDate.CheckAddress == false {
					// recursive for getting write address
					wc.RunTimeDate.TaskStatus = pb.RunTimeData_recursive
					var cover cover.Cover
					cover.Merge(info.Calls[wc.RunTimeDate.Idx].Cover)
					for _, c := range wc.CriticalCondition {
						if checkCondition(c, cover) {

						}
					}
				} else {
					wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested
				}

			} else if wc.RunTimeDate.TaskStatus == pb.RunTimeData_recursive && len(wc.WriteAddress) != 0 {
				for _, wwa := range wc.WriteAddress {
					if ok, info := proc.dependencyWriteAddress(wwa); ok {

						copy(wc.RunTimeDate.Program, wwa.RunTimeDate.Program)

						if wwa.RunTimeDate.Parent.Address == wc.RunTimeDate.Address {

							wc.RunTimeDate.CheckAddress = true
							wa.RunTimeDate.CheckAddress = true

							if checkAddress(wc.RunTimeDate.Parent.ConditionAddress, info.Calls[wc.RunTimeDate.Parent.Idx].Cover) {
								wc.RunTimeDate.Parent.CheckCondition = true
								if checkAddress(wc.RunTimeDate.Parent.Address, info.Calls[wc.RunTimeDate.Parent.Idx].Cover) {

									wc.RunTimeDate.Parent.CheckAddress = true

									wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested

									copy(wc.RunTimeDate.Program, wwa.RunTimeDate.Program)

									return true, info
								}
							}
						} else {
							// recursive for getting next critical condition
							var cover cover.Cover
							cover.Merge(info.Calls[wc.RunTimeDate.Idx].Cover)
							for _, condition := range wc.CriticalCondition {
								if checkCondition(condition, cover) {

								} else {
									wc.RunTimeDate.ConditionAddress = condition.ConditionAddress
								}
							}
						}
					}
				}

				wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested
				for _, wwa := range wc.WriteAddress {
					if wwa.RunTimeDate.TaskStatus != pb.RunTimeData_tested {
						wc.RunTimeDate.TaskStatus = wwa.RunTimeDate.TaskStatus
					}
				}

			} else if wc.RunTimeDate.TaskStatus == pb.RunTimeData_tested {

			}
		}

		wa.RunTimeDate.TaskStatus = pb.RunTimeData_tested
		for _, wc := range wa.WriteSyscall {
			if wc.RunTimeDate.TaskStatus != pb.RunTimeData_tested {
				wa.RunTimeDate.TaskStatus = wc.RunTimeDate.TaskStatus
			}
		}
	}

	return false, nil
}

func testSyscall(wc *pb.Syscall) (res bool) {
	return false
}

func forprogam() {

	// for repeat
	//if wa.Repeat == 0 {
	//	mini := 1
	//	wa.Repeat = uint32(proc.rnd.Int31n(int32(programLength-len(p.Calls))-int32(mini)) + int32(mini))
	//}
	// log.Logf(1, "repeat : %v", wa.Repeat)

	//	for _, wi := range wa.WriteInput {
	//
	//		log.Logf(1, "write program : \n%s", wi.Program)
	//		proc.fuzzer.dManager.SendLog(fmt.Sprintf("write program : \n%s", wi.Program))
	//
	//		wp, err := proc.fuzzer.target.Deserialize(wi.Program, prog.NonStrict)
	//		if err != nil {
	//			log.Fatalf("failed to deserialize program from write program: %v", err)
	//		}
	//		wpInfo := proc.execute(proc.execOptsCover, wp, ProgNormal, StatDependency)
	//		u.RunTimeDate.CheckAddress = checkAddress(wi.WriteAddress, wpInfo.Calls[wi.Idx].Cover)
	//
	//		p0 := p.Clone()
	//		p0.Splice(wp, u.Idx, programLength)
	//
	//		data := p0.Serialize()
	//		log.Logf(1, "test case with write program : \n%s", data)
	//		proc.fuzzer.dManager.SendLog(fmt.Sprintf("test case with write program : \n%s", data))
	//
	//		info := proc.execute(proc.execOptsCover, p0, ProgNormal, StatDependency)
	//		u.RunTimeDate.CheckAddress = checkAddress(wi.WriteAddress, info.Calls[wi.Idx].Cover)
	//
	//		ok1, ok2, ok3 := proc.checkCoverage(p, inputCover)
	//		if ok1 {
	//			proc.fuzzer.dManager.SendLog(fmt.Sprintf("checkWriteAddress : %x", p.WriteAddress))
	//		} else {
	//			proc.fuzzer.dManager.SendLog(fmt.Sprintf("not checkWriteAddress : %x", p.WriteAddress))
	//		}
	//		if ok2 {
	//			proc.fuzzer.dManager.SendLog(fmt.Sprintf("checkConditionAddress : %x", p.Uncover[p.UncoverIdx].ConditionAddress))
	//		} else {
	//			proc.fuzzer.dManager.SendLog(fmt.Sprintf("not checkConditionAddress : %x", p.Uncover[p.UncoverIdx].ConditionAddress))
	//		}
	//		if ok3 {
	//			u.RunTimeDate.CheckAddress = true
	//			goto cover
	//		} else {
	//
	//		}
	//	}
}

func checkAddress(conditionAddress uint32, cover []uint32) (res bool) {
	res = false
	for _, c := range cover {
		if c == conditionAddress {
			res = true
			return
		}
	}
	return
}

func checkCondition(condition *pb.Condition, cover cover.Cover) (res bool) {
	res = false
	if _, ok := cover[condition.ConditionAddress]; ok {
		for _, a := range condition.RightBranchAddress {
			if _, ok := cover[a]; ok {
				res = true
				return
			}
		}
	}
	return
}

func checkCriticalCondition() (condition *pb.Condition) {

}

func reexecutionSuccess(info *ipc.ProgInfo, oldInfo *ipc.CallInfo, call int) bool {
	if info == nil || len(info.Calls) == 0 {
		return false
	}
	if call != -1 {
		// Don't minimize calls from successful to unsuccessful.
		// Successful calls are much more valuable.
		if oldInfo.Errno == 0 && info.Calls[call].Errno != 0 {
			return false
		}
		return len(info.Calls[call].Signal) != 0
	}
	return len(info.Extra.Signal) != 0
}

func getSignalAndCover(p *prog.Prog, info *ipc.ProgInfo, call int) (signal.Signal, []uint32) {
	inf := &info.Extra
	if call != -1 {
		inf = &info.Calls[call]
	}
	return signal.FromRaw(inf.Signal, signalPrio(p, inf, call)), inf.Cover
}

func (proc *Proc) smashInput(item *WorkSmash) {
	if proc.fuzzer.faultInjectionEnabled && item.call != -1 {
		proc.failCall(item.p, item.call)
	}
	if proc.fuzzer.comparisonTracingEnabled && item.call != -1 {
		proc.executeHintSeed(item.p, item.call)
	}
	corpus := proc.fuzzer.corpusSnapshot()
	for i := 0; i < 100; i++ {
		p := item.p.Clone()
		p.Mutate(proc.rnd, programLength, proc.fuzzer.choiceTable, corpus)
		log.Logf(1, "#%v: smash mutated", proc.pid)
		proc.execute(proc.execOpts, p, ProgNormal, StatSmash)
	}
}

func (proc *Proc) failCall(p *prog.Prog, call int) {
	for nth := 0; nth < 100; nth++ {
		log.Logf(1, "#%v: injecting fault into call %v/%v", proc.pid, call, nth)
		opts := *proc.execOpts
		opts.Flags |= ipc.FlagInjectFault
		opts.FaultCall = call
		opts.FaultNth = nth
		info := proc.executeRaw(&opts, p, StatSmash)
		if info != nil && len(info.Calls) > call && info.Calls[call].Flags&ipc.CallFaultInjected == 0 {
			break
		}
	}
}

func (proc *Proc) executeHintSeed(p *prog.Prog, call int) {
	log.Logf(1, "#%v: collecting comparisons", proc.pid)
	// First execute the original program to dump comparisons from KCOV.
	info := proc.execute(proc.execOptsComps, p, ProgNormal, StatSeed)
	if info == nil {
		return
	}

	// Then mutate the initial program for every match between
	// a syscall argument and a comparison operand.
	// Execute each of such mutants to check if it gives new coverage.
	p.MutateWithHints(call, info.Calls[call].Comps, func(p *prog.Prog) {
		log.Logf(1, "#%v: executing comparison hint", proc.pid)
		proc.execute(proc.execOpts, p, ProgNormal, StatHint)
	})
}

func (proc *Proc) execute(execOpts *ipc.ExecOpts, p *prog.Prog, flags ProgTypes, stat Stat) *ipc.ProgInfo {
	info := proc.executeRaw(execOpts, p, stat)
	calls, extra := proc.fuzzer.checkNewSignal(p, info)
	for _, callIndex := range calls {
		if stat == StatDependency {
			proc.fuzzer.dManager.SendLog(fmt.Sprintf("new input from StatDependency : %v", p))
			p.Comments = append(p.Comments, "StatDependency")
		}
		proc.enqueueCallTriage(p, flags, callIndex, info.Calls[callIndex])
	}
	if extra {
		proc.enqueueCallTriage(p, flags, -1, info.Extra)
	}

	//ccalls := proc.fuzzer.checkNewCoverage(p, info)
	//for _, callIndex := range ccalls {
	//	proc.enqueueCallTriage(p, flags, callIndex, info.Calls[callIndex])
	//}

	return info
}

func (proc *Proc) enqueueCallTriage(p *prog.Prog, flags ProgTypes, callIndex int, info ipc.CallInfo) {
	// info.Signal points to the output shmem region, detach it before queueing.
	info.Signal = append([]uint32{}, info.Signal...)
	// None of the caller use Cover, so just nil it instead of detaching.
	// Note: triage input uses executeRaw to get coverage.
	info.Cover = nil
	proc.fuzzer.workQueue.enqueue(&WorkTriage{
		p:     p.Clone(),
		call:  callIndex,
		info:  info,
		flags: flags,
	})
}

func (proc *Proc) executeRaw(opts *ipc.ExecOpts, p *prog.Prog, stat Stat) *ipc.ProgInfo {
	if opts.Flags&ipc.FlagDedupCover == 0 {
		log.Fatalf("dedup cover is not enabled")
	}

	// Limit concurrency window and do leak checking once in a while.
	ticket := proc.fuzzer.gate.Enter()
	defer proc.fuzzer.gate.Leave(ticket)

	proc.logProgram(opts, p)
	for try := 0; ; try++ {
		atomic.AddUint64(&proc.fuzzer.stats[stat], 1)
		output, info, hanged, err := proc.env.Exec(opts, p)
		if err != nil {
			if try > 10 {
				log.Fatalf("executor %v failed %v times:\n%v", proc.pid, try, err)
			}
			log.Logf(4, "fuzzer detected executor failure='%v', retrying #%d", err, try+1)
			debug.FreeOSMemory()
			time.Sleep(time.Second)
			continue
		}
		log.Logf(2, "result hanged=%v: %s", hanged, output)
		return info
	}
}

func (proc *Proc) logProgram(opts *ipc.ExecOpts, p *prog.Prog) {
	if proc.fuzzer.outputType == OutputNone {
		return
	}

	data := p.Serialize()
	strOpts := ""
	if opts.Flags&ipc.FlagInjectFault != 0 {
		strOpts = fmt.Sprintf(" (fault-call:%v fault-nth:%v)", opts.FaultCall, opts.FaultNth)
	}

	// The following output helps to understand what program crashed kernel.
	// It must not be intermixed.
	switch proc.fuzzer.outputType {
	case OutputStdout:
		now := time.Now()
		proc.fuzzer.logMu.Lock()
		fmt.Printf("%02v:%02v:%02v executing program %v%v:\n%s\n",
			now.Hour(), now.Minute(), now.Second(),
			proc.pid, strOpts, data)
		proc.fuzzer.logMu.Unlock()
	case OutputDmesg:
		fd, err := syscall.Open("/dev/kmsg", syscall.O_WRONLY, 0)
		if err == nil {
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, "syzkaller: executing program %v%v:\n%s\n",
				proc.pid, strOpts, data)
			syscall.Write(fd, buf.Bytes())
			syscall.Close(fd)
		}
	case OutputFile:
		f, err := os.Create(fmt.Sprintf("%v-%v.prog", proc.fuzzer.name, proc.pid))
		if err == nil {
			if strOpts != "" {
				fmt.Fprintf(f, "#%v\n", strOpts)
			}
			f.Write(data)
			f.Close()
		}
	default:
		log.Fatalf("unknown output type: %v", proc.fuzzer.outputType)
	}
}
