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

		//proc.fuzzer.checkNewCoverage(item.p, info)
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
		Sig:        sig.String(),
		Program:    []byte{},
		Call:       make(map[uint32]*pb.Call),
		Dependency: false,
	}

	for _, c := range data {
		input.Program = append(input.Program, c)
	}

	//log.Logf(2, "data :\n%s", data)
	//log.Logf(2, "input.Program :\n%s", input.Program)

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
			proc.fuzzer.dManager.SendLog(fmt.Sprintf("real new input from StatDependency : \n%s", data))
		}
	}

	proc.fuzzer.dManager.SendNewInput(&input)
}

func (proc *Proc) getCall(sc *pb.IoctlCmd) (res *prog.Syscall) {
	// only work for ioctl
	for n, c := range proc.fuzzer.target.SyscallMap {
		ok := strings.HasPrefix(n, sc.Name)
		if !ok {
			continue
		}
		for _, a := range c.Args {
			ok2 := a.FieldName() == "cmd"
			if !ok2 {
				continue
			}
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
	return
}

func (proc *Proc) dependencyMutate(item *WorkDependency) {

	log.Logf(1, "#%v: DependencyMutate", proc.pid)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("#%v: DependencyMutate", proc.pid))

	task := item.task
	log.Logf(1, "DependencyMutate program : \n%s", task.Program)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("DependencyMutate program : \n%s", task.Program))
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("index  : %d\n write index : %d", task.Index, task.WriteIndex))

	ct := proc.fuzzer.choiceTable

	p, err := proc.fuzzer.target.Deserialize(task.Program, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependencyMutate failed to deserialize program from task.Program: %v", err)
	}

	wp, err := proc.fuzzer.target.Deserialize(task.WriteProgram, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependencyMutate failed to deserialize program from task.WriteProgram: %v", err)
	}
	var usefulCall []*prog.Call
	if int(task.WriteIndex) > len(wp.Calls) {
		log.Fatalf("dependencyMutate int(task.WriteIndex) > len(wp.Calls)")
	}
	for i, c := range wp.Calls {
		if i <= int(task.WriteIndex) {
			usefulCall = append(usefulCall, c)
		} else {
			break
		}
	}
	wdata := wp.Serialize()
	log.Logf(1, "usefulCall program : \n%s", wdata)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("usefulCall program : \n%s", wdata))

	idx := task.Index

	// need combine open
	p.Calls = append(p.Calls[:idx], append(usefulCall, p.Calls[idx:]...)...)
	for i := len(p.Calls) - 1; i >= programLength; i-- {
		p.RemoveCall(i)
	}
	data := p.Serialize()
	log.Logf(1, "final program : \n%s", data)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("final program : \n%s", data))

	idx = idx + task.WriteIndex + 1

	var info *ipc.ProgInfo

	info = proc.execute(proc.execOptsCover, wp, ProgNormal, StatDependency)
	checkWriteAddress1 := checkAddress(task.WriteAddress, info.Calls[task.WriteIndex].Cover)
	if checkWriteAddress1 {
		task.CheckWriteAddress = true
		log.Logf(1, "write program could arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("write input could arrive at write address : %d", task.WriteAddress))
	} else {
		log.Logf(1, "write program could not arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("write input could not arrive at write address : %d", task.WriteAddress))
	}

	info = proc.execute(proc.execOptsCover, p, ProgNormal, StatDependency)
	checkWriteAddress2 := checkAddress(task.WriteAddress, info.Calls[task.WriteIndex].Cover)
	if checkWriteAddress2 {
		task.CheckWriteAddressFinal = true
		log.Logf(1, "write program could arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("write input could arrive at write address : %d", task.WriteAddress))
	} else {
		log.Logf(1, "final program could not arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("final input could not arrive at write address : %d", task.WriteAddress))
	}

	for i := 0; i < 20; i++ {
		info = proc.execute(proc.execOptsCover, p, ProgNormal, StatDependency)
		var cov cover.Cover
		cov.Merge(info.Calls[idx].Cover)

		for u, r := range task.UncoveredAddress {
			checkConditionAddress := checkAddressMap(r.ConditionAddress, cov)
			if !checkConditionAddress {
				continue
			}

			checkUncoveredAddress := checkAddressMap(r.Address, cov)
			if !checkUncoveredAddress {
				continue
			}

			proc.fuzzer.dManager.SendLog(fmt.Sprintf("cover uncovered address : %x", r.Address))

			r := pb.CloneRunTimeData(task.UncoveredAddress[u])
			task.CoveredAddress[u] = r

			data := p.Serialize()
			for _, c := range data {
				r.Program = append(r.Program, c)
			}
			r.TaskStatus = pb.TaskStatus_covered
			r.Idx = idx
			r.CheckAddress = true

			delete(task.UncoveredAddress, u)
		}

		if len(task.UncoveredAddress) > 0 {
			p.MutateIoctl3Arg(proc.rnd, idx, ct)
		} else {
			break
		}
	}

	task.TaskStatus = pb.TaskStatus_tested
	tasks := &pb.Tasks{
		Name: proc.fuzzer.name,
		Task: []*pb.Task{},
	}
	tasks.Task = append(tasks.Task, task)
	proc.fuzzer.dManager.ReturnTasks(tasks)

	return
}

//
//func (proc *Proc) dependencyMutate(item *WorkDependency) {
//
//	log.Logf(1, "#%v: DependencyMutate", proc.pid)
//	proc.fuzzer.dManager.SendLog(fmt.Sprintf("#%v: DependencyMutate", proc.pid))
//
//	dependencyInput := item.dependencyInput
//	log.Logf(1, "DependencyMutate program : \n%s", dependencyInput.Program)
//	proc.fuzzer.dManager.SendLog(fmt.Sprintf("DependencyMutate program : \n%s", dependencyInput.Program))
//
//	for _, u := range dependencyInput.UncoveredAddress {
//		proc.fuzzer.dManager.SendLog(fmt.Sprintf("UncoveredAddress : %v", u.UncoveredAddress))
//		for _, wa := range u.WriteAddress {
//			log.Logf(2, "dependencyMutate data :\n%s", wa.RunTimeDate.Program)
//			if ok, _ := proc.dependencyWriteAddress(wa); ok {
//				updateRunTimeData(u.RunTimeDate, wa.RunTimeDate)
//				u.RunTimeDate.CheckAddress = true
//				u.RunTimeDate.TaskStatus = pb.RunTimeData_cover
//				break
//			} else {
//			}
//		}
//		updateRunTimeDataTaskStatusU(u)
//	}
//	_, _ = proc.fuzzer.dManager.ReturnDependencyInput(dependencyInput)
//	return
//}
//
//func (proc *Proc) dependencyWriteAddress(wa *pb.WriteAddress) (res bool, info *ipc.ProgInfo) {
//	for _, wc := range wa.WriteSyscall {
//		log.Logf(2, "dependencyWriteAddress data :\n%s", wc.RunTimeDate.Program)
//		info = proc.dependencyRecursiveWriteSyscall(wc)
//		if wc.RunTimeDate.TaskStatus == pb.RunTimeData_cover {
//			updatePRunTimeData(wa.RunTimeDate, wc.RunTimeDate)
//
//			p, err := proc.fuzzer.target.Deserialize(wc.RunTimeDate.Program, prog.NonStrict)
//			if err != nil {
//				log.Fatalf("failed to deserialize program from dependencyRecursiveWriteAddress: %v", err)
//			}
//			tempInfo := info
//			ct := proc.fuzzer.choiceTable
//			for i := 0; i < 10; i++ {
//				address := checkAddress(wa.RunTimeDate.Address, tempInfo.Calls[wa.RunTimeDate.Idx].Cover)
//				conditionAddress := checkAddress(wa.RunTimeDate.ConditionAddress, tempInfo.Calls[wa.RunTimeDate.Idx].Cover)
//				if conditionAddress == true {
//					if address == true {
//						// arrive at address
//						updateRunTimeDataCover(wa.RunTimeDate)
//						data := p.Serialize()
//						for _, c := range data {
//							wa.RunTimeDate.Program = append(wa.RunTimeDate.Program, c)
//						}
//						info = tempInfo
//						return true, info
//					}
//				} else {
//
//				}
//
//				p.MutateIoctl3Arg(proc.rnd, wc.RunTimeDate.Idx, ct)
//				tempInfo = proc.execute(proc.execOptsCover, p, ProgNormal, StatDependency)
//			}
//			wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested
//		} else if wc.RunTimeDate.TaskStatus == pb.RunTimeData_recursive {
//			// can not arrive at write address
//			// recursive for getting next critical condition
//			checkCriticalCondition(wc, info)
//		} else {
//
//		}
//	}
//
//	updateRunTimeDataTaskStatusWa(wa)
//	return false, nil
//}
//
//func (proc *Proc) dependencyRecursiveWriteAddress(wa *pb.WriteAddress) (info *ipc.ProgInfo) {
//	for _, wc := range wa.WriteSyscall {
//		info = proc.dependencyRecursiveWriteSyscall(wc)
//		if wc.RunTimeDate.TaskStatus == pb.RunTimeData_cover {
//			updatePRunTimeData(wa.RunTimeDate, wc.RunTimeDate)
//
//			p, err := proc.fuzzer.target.Deserialize(wc.RunTimeDate.Program, prog.NonStrict)
//			if err != nil {
//				log.Fatalf("failed to deserialize program from dependencyRecursiveWriteAddress: %v", err)
//			}
//			tempInfo := info
//			ct := proc.fuzzer.choiceTable
//			for i := 0; i < 10; i++ {
//				address := checkAddress(wa.RunTimeDate.Address, tempInfo.Calls[wa.RunTimeDate.Idx].Cover)
//				var cover cover.Cover
//				cover.Merge(info.Calls[wa.RunTimeDate.Idx].Cover)
//				rightBranchAddress := checkAddresses(wa.RunTimeDate.RightBranchAddress, cover)
//
//				if address == true {
//					// arrive at address
//					updateRunTimeDataCover(wa.RunTimeDate)
//					data := p.Serialize()
//					for _, c := range data {
//						wa.RunTimeDate.Program = append(wa.RunTimeDate.Program, c)
//					}
//					info = tempInfo
//					return info
//				} else if rightBranchAddress == true {
//					// arrive at right branch of critical condition
//					// continue next critical condition
//					if wa.RunTimeDate.CheckRightBranchAddress != true {
//						wa.RunTimeDate.TaskStatus = pb.RunTimeData_recursive
//						wa.RunTimeDate.CheckRightBranchAddress = true
//						data := p.Serialize()
//						for _, c := range data {
//							wa.RunTimeDate.Program = append(wa.RunTimeDate.Program, c)
//						}
//						info = tempInfo
//					}
//				}
//
//				p.MutateIoctl3Arg(proc.rnd, wc.RunTimeDate.Idx, ct)
//				tempInfo = proc.execute(proc.execOptsCover, p, ProgNormal, StatDependency)
//			}
//			wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested
//		} else if wc.RunTimeDate.TaskStatus == pb.RunTimeData_recursive {
//			// can not arrive at write address
//			// recursive for getting next critical condition
//			checkCriticalCondition(wc, info)
//		} else {
//
//		}
//	}
//	updateRunTimeDataTaskStatusWa(wa)
//	return info
//}
//
//// return true once arrive at write address for write syscall
//func (proc *Proc) dependencyRecursiveWriteSyscall(wc *pb.Syscall) (info *ipc.ProgInfo) {
//	log.Logf(2, "dependencyRecursiveWriteSyscall data :\n%s", wc.RunTimeDate.Program)
//	if wc.RunTimeDate.TaskStatus == pb.RunTimeData_untested {
//		return proc.dependencyWriteSyscallUntested(wc)
//	} else if wc.RunTimeDate.TaskStatus == pb.RunTimeData_recursive {
//		return proc.dependencyWriteSyscallRecursive(wc)
//	} else if wc.RunTimeDate.TaskStatus == pb.RunTimeData_tested {
//
//	} else if wc.RunTimeDate.TaskStatus == pb.RunTimeData_out {
//
//	} else if wc.RunTimeDate.TaskStatus == pb.RunTimeData_cover {
//
//	}
//	return nil
//}
//
//func (proc *Proc) dependencyWriteSyscallUntested(wc *pb.Syscall) (info *ipc.ProgInfo) {
//	ct := proc.fuzzer.choiceTable
//
//	log.Logf(2, "dependencyWriteSyscallUntested data :\n%s", wc.RunTimeDate.Program)
//	p, err := proc.fuzzer.target.Deserialize(wc.RunTimeDate.Program, prog.NonStrict)
//	if err != nil {
//		log.Fatalf("failed to deserialize program from dependencyRecursiveWriteAddress: %v", err)
//	}
//	call := proc.getCall(wc)
//	if call == nil {
//		wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested
//		return nil
//	}
//	log.Logf(2, "dependencyWriteSyscallUntested wc.RunTimeDate.Idx : %v", wc.RunTimeDate.Idx)
//	c0c := p.GetCall(proc.rnd, call, wc.RunTimeDate.Idx, ct)
//	p.InsertCall(c0c, wc.RunTimeDate.Idx, programLength)
//
//	data := p.Serialize()
//	for _, c := range data {
//		wc.RunTimeDate.Program = append(wc.RunTimeDate.Program, c)
//	}
//	size := uint32(len(c0c))
//	wc.RunTimeDate.Idx = wc.RunTimeDate.Idx + size - 1
//
//	return proc.dependencyWriteSyscallMutateArgument(wc)
//}
//
//func (proc *Proc) dependencyWriteSyscallMutateArgument(wc *pb.Syscall) (info *ipc.ProgInfo) {
//	p, err := proc.fuzzer.target.Deserialize(wc.RunTimeDate.Program, prog.NonStrict)
//	if err != nil {
//		log.Fatalf("failed to deserialize program from dependencyRecursiveWriteAddress: %v", err)
//	}
//	ct := proc.fuzzer.choiceTable
//	for i := 0; i < 10; i++ {
//		info = proc.execute(proc.execOptsCover, p, ProgNormal, StatDependency)
//		address := checkAddress(wc.RunTimeDate.Address, info.Calls[wc.RunTimeDate.Idx].Cover)
//		if address == true {
//			updateRunTimeDataCover(wc.RunTimeDate)
//			data := p.Serialize()
//			for _, c := range data {
//				wc.RunTimeDate.Program = append(wc.RunTimeDate.Program, c)
//			}
//			return info
//		}
//		p.MutateIoctl3Arg(proc.rnd, wc.RunTimeDate.Idx, ct)
//	}
//	wc.RunTimeDate.TaskStatus = pb.RunTimeData_recursive
//	return info
//}
//
//func (proc *Proc) dependencyWriteSyscallRecursive(wc *pb.Syscall) (info *ipc.ProgInfo) {
//	for _, wa := range wc.WriteAddress {
//		WAinfo := proc.dependencyRecursiveWriteAddress(wa)
//		if wa.RunTimeDate.TaskStatus == pb.RunTimeData_cover {
//			updateRunTimeData(wc.RunTimeDate, wa.RunTimeDate)
//			updateRunTimeDataCover(wc.RunTimeDate)
//			return WAinfo
//		} else if wa.RunTimeDate.TaskStatus == pb.RunTimeData_recursive {
//			if wa.RunTimeDate.CheckRightBranchAddress == true {
//				updateRunTimeData(wc.RunTimeDate, wa.RunTimeDate)
//				WSinfo := proc.dependencyWriteSyscallMutateArgument(wc)
//				if wc.RunTimeDate.TaskStatus == pb.RunTimeData_cover {
//					return WSinfo
//				} else {
//					wc.RunTimeDate.TaskStatus = pb.RunTimeData_recursive
//					info = WAinfo
//				}
//			}
//		}
//	}
//	updateRunTimeDataTaskStatusWc(wc)
//	return info
//}
//
//func forprogam() {
//
//	// for repeat
//	//if wa.Repeat == 0 {
//	//	mini := 1
//	//	wa.Repeat = uint32(proc.rnd.Int31n(int32(programLength-len(p.Calls))-int32(mini)) + int32(mini))
//	//}
//	// log.Logf(1, "repeat : %v", wa.Repeat)
//
//	//	for _, wi := range wa.WriteInput {
//	//
//	//		log.Logf(1, "write program : \n%s", wi.Program)
//	//		proc.fuzzer.dManager.SendLog(fmt.Sprintf("write program : \n%s", wi.Program))
//	//
//	//		wp, err := proc.fuzzer.target.Deserialize(wi.Program, prog.NonStrict)
//	//		if err != nil {
//	//			log.Fatalf("failed to deserialize program from write program: %v", err)
//	//		}
//	//		wpInfo := proc.execute(proc.execOptsCover, wp, ProgNormal, StatDependency)
//	//		u.RunTimeDate.CheckAddress = checkAddress(wi.WriteAddress, wpInfo.Calls[wi.Idx].Cover)
//	//
//	//		p0 := p.Clone()
//	//		p0.Splice(wp, u.Idx, programLength)
//	//
//	//		data := p0.Serialize()
//	//		log.Logf(1, "test case with write program : \n%s", data)
//	//		proc.fuzzer.dManager.SendLog(fmt.Sprintf("test case with write program : \n%s", data))
//	//
//	//		info := proc.execute(proc.execOptsCover, p0, ProgNormal, StatDependency)
//	//		u.RunTimeDate.CheckAddress = checkAddress(wi.WriteAddress, info.Calls[wi.Idx].Cover)
//	//
//	//		ok1, ok2, ok3 := proc.checkCoverage(p, inputCover)
//	//		if ok1 {
//	//			proc.fuzzer.dManager.SendLog(fmt.Sprintf("checkWriteAddress : %x", p.WriteAddress))
//	//		} else {
//	//			proc.fuzzer.dManager.SendLog(fmt.Sprintf("not checkWriteAddress : %x", p.WriteAddress))
//	//		}
//	//		if ok2 {
//	//			proc.fuzzer.dManager.SendLog(fmt.Sprintf("checkConditionAddress : %x", p.Uncover[p.UncoverIdx].ConditionAddress))
//	//		} else {
//	//			proc.fuzzer.dManager.SendLog(fmt.Sprintf("not checkConditionAddress : %x", p.Uncover[p.UncoverIdx].ConditionAddress))
//	//		}
//	//		if ok3 {
//	//			u.RunTimeDate.CheckAddress = true
//	//			goto cover
//	//		} else {
//	//
//	//		}
//	//	}
//}
//
//func updateRunTimeData(parent *pb.RunTimeData, child *pb.RunTimeData) {
//	for _, c := range child.Program {
//		parent.Program = append(parent.Program, c)
//	}
//	parent.Idx = child.Idx
//}
//
//func updatePRunTimeData(parent *pb.RunTimeData, child *pb.RunTimeData) {
//	for _, c := range child.Program {
//		parent.Program = append(parent.Program, c)
//	}
//	parent.Idx = child.Idx + 1
//}
//
//func updateRunTimeDataCover(parent *pb.RunTimeData) {
//	parent.TaskStatus = pb.RunTimeData_cover
//	parent.CheckAddress = true
//}
//
//func updateRunTimeDataTaskStatusWc(wc *pb.Syscall) {
//	untested := 0
//	recursive := 0
//	tested := 0
//	out := 0
//	cover := 0
//	for _, wa := range wc.WriteAddress {
//		switch wa.RunTimeDate.TaskStatus {
//		case pb.RunTimeData_untested:
//			untested++
//		case pb.RunTimeData_recursive:
//			recursive++
//		case pb.RunTimeData_tested:
//			tested++
//		case pb.RunTimeData_out:
//			out++
//		case pb.RunTimeData_cover:
//			cover++
//		default:
//
//		}
//	}
//
//	if wc.RunTimeDate.TaskStatus != pb.RunTimeData_recursive {
//		if recursive > 0 {
//			wc.RunTimeDate.TaskStatus = pb.RunTimeData_recursive
//		} else if out > 0 {
//			wc.RunTimeDate.TaskStatus = pb.RunTimeData_out
//		} else if tested > 0 {
//			wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested
//		} else if cover > 0 {
//			wc.RunTimeDate.TaskStatus = pb.RunTimeData_tested
//		}
//	}
//}
//
//func updateRunTimeDataTaskStatusWa(wa *pb.WriteAddress) {
//	untested := 0
//	recursive := 0
//	tested := 0
//	out := 0
//	cover := 0
//	for _, wc := range wa.WriteSyscall {
//		switch wc.RunTimeDate.TaskStatus {
//		case pb.RunTimeData_untested:
//			untested++
//		case pb.RunTimeData_recursive:
//			recursive++
//		case pb.RunTimeData_tested:
//			tested++
//		case pb.RunTimeData_out:
//			out++
//		case pb.RunTimeData_cover:
//			cover++
//		default:
//
//		}
//	}
//	if cover > 0 {
//		wa.RunTimeDate.TaskStatus = pb.RunTimeData_cover
//	} else if untested > 0 {
//		wa.RunTimeDate.TaskStatus = pb.RunTimeData_untested
//	} else if recursive > 0 {
//		wa.RunTimeDate.TaskStatus = pb.RunTimeData_recursive
//	} else if out > 0 {
//		wa.RunTimeDate.TaskStatus = pb.RunTimeData_out
//	} else if tested > 0 {
//		wa.RunTimeDate.TaskStatus = pb.RunTimeData_tested
//	}
//}
//
//func updateRunTimeDataTaskStatusU(u *pb.UncoveredAddress) {
//	untested := 0
//	recursive := 0
//	tested := 0
//	out := 0
//	cover := 0
//	for _, wa := range u.WriteAddress {
//		switch wa.RunTimeDate.TaskStatus {
//		case pb.RunTimeData_untested:
//			untested++
//		case pb.RunTimeData_recursive:
//			recursive++
//		case pb.RunTimeData_tested:
//			tested++
//		case pb.RunTimeData_out:
//			out++
//		case pb.RunTimeData_cover:
//			cover++
//		default:
//
//		}
//	}
//	if cover > 0 {
//		u.RunTimeDate.TaskStatus = pb.RunTimeData_cover
//	} else if untested > 0 {
//		u.RunTimeDate.TaskStatus = pb.RunTimeData_untested
//	} else if recursive > 0 {
//		u.RunTimeDate.TaskStatus = pb.RunTimeData_recursive
//	} else if out > 0 {
//		u.RunTimeDate.TaskStatus = pb.RunTimeData_out
//	} else if tested > 0 {
//		u.RunTimeDate.TaskStatus = pb.RunTimeData_tested
//	}
//}
//
func checkAddress(Address uint32, cover []uint32) (res bool) {
	res = false
	for _, c := range cover {
		if c == Address {
			res = true
			return
		}
	}
	return
}

func checkAddressMap(Address uint32, cover cover.Cover) (res bool) {
	res = false
	if _, ok := cover[Address]; ok {
		res = true
		return
	}
	return
}

func checkCondition(condition *pb.Condition, cover cover.Cover) (res bool) {
	res = false
	if _, ok := cover[condition.SyzkallerConditionAddress]; ok {
		for _, a := range condition.SyzkallerRightBranchAddress {
			if _, ok := cover[a]; ok {
				res = true
				return
			}
		}
	}
	return
}

func checkAddresses(Address []uint32, cover cover.Cover) (res bool) {
	//func checkAddresses(Address map[uint32]uint32, cover cover.Cover) (res bool) {
	res = false
	for _, a := range Address {
		if _, ok := cover[a]; ok {
			res = true
			return
		}
	}
	return
}

func checkAddressesMap(Address map[uint32]uint32, cover cover.Cover) (res bool) {
	//func checkAddresses(Address map[uint32]uint32, cover cover.Cover) (res bool) {
	res = false
	for a := range Address {
		if _, ok := cover[a]; ok {
			res = true
			return
		}
	}
	return
}

//func checkCriticalCondition(wc *pb.Syscall, info *ipc.ProgInfo) (res bool) {
//	var cover cover.Cover
//	cover.Merge(info.Calls[wc.RunTimeDate.Idx].Cover)
//	for _, condition := range wc.CriticalCondition {
//		if checkCondition(condition, cover) {
//
//		} else {
//			wc.RunTimeDate.TaskStatus = pb.RunTimeData_recursive
//			wc.RunTimeDate.ConditionAddress = condition.SyzkallerConditionAddress
//			for _, ra := range condition.SyzkallerRightBranchAddress {
//				//wc.RunTimeDate.RightBranchAddress[ra] = 0
//				wc.RunTimeDate.RightBranchAddress = append(wc.RunTimeDate.RightBranchAddress, ra)
//			}
//			wc.WriteAddress = nil
//			return true
//		}
//	}
//	log.Logf(1, "every critical condition is right but we can not arrive at write address")
//	return false
//}

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
			proc.fuzzer.dManager.SendLog(fmt.Sprintf("new input from StatDependency : \n%s", p.Serialize()))
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
