package main

import (
	"fmt"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/cover"
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/hash"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/ipc"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/log"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/prog"
	"github.com/golang/protobuf/proto"
	"strings"
)

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

func (proc *Proc) executeDependencyHintSeed(p *prog.Prog, call int) {
	log.Logf(pb.DebugLevel, "#%v: collecting comparisons", proc.pid)
	// First execute the original program to dump comparisons from KCOV.
	info := proc.execute(proc.execOptsComps, p, ProgNormal, StatDependency)
	if info == nil {
		return
	}

	// Then mutate the initial program for every match between
	// a syscall argument and a comparison operand.
	// Execute each of such mutants to check if it gives new coverage.
	p.MutateWithHints(call, info.Calls[call].Comps, func(p *prog.Prog) {
		log.Logf(pb.DebugLevel, "#%v: executing comparison hint", proc.pid)
		proc.execute(proc.execOpts, p, ProgNormal, StatDependency)
	})
}

func dealAddress(open string) string {
	res := ""
	add := true
	for i, c := range open {
		if c == '&' && open[i+1] == '(' {
			add = false
		}
		if add {
			res = res + string(c)
		} else {
			if c == ')' {
				add = true
			}
		}
	}
	return res
}

func removeSameResource(p []byte) ([]byte, []uint32) {
	var calls []string
	var i = 0
	var j = 0
	calls = append(calls, "")
	for _, c := range p {
		if c != '\n' {
			calls[i] = calls[i] + string(c)
		} else {
			i++
			calls = append(calls, "")
		}
	}
	var newIdx = make([]uint32, len(calls))

	for i, call := range calls {
		if strings.Index(call, " = ") != -1 {
			j = strings.Index(call, " = ")
			res := call[:j]
			res = res + ","
			open := call[j:]
			open = dealAddress(open)
			for k := i + 1; k < len(calls); k++ {
				call := dealAddress(calls[k])
				if strings.Contains(call, open) {
					if strings.Index(calls[k], " = ") != -1 {
						rres := calls[k][:strings.Index(calls[k], " = ")]
						rres = rres + ","
						for h := k + 1; h < len(calls); h++ {
							calls[h] = strings.ReplaceAll(calls[h], rres, res)
						}
						calls[k] = ""
						newIdx[k] = uint32(i)
					}
				}
			}
		}
	}

	var pc []byte
	var idx uint32
	idx = 0
	for i, call := range calls {
		if call != "" {
			callb := []byte(call)
			for _, c := range callb {
				pc = append(pc, c)
			}
			pc = append(pc, '\n')
			newIdx[i] = idx
			idx++
		} else {
			newIdx[i] = newIdx[newIdx[i]]
		}
	}

	return pc, newIdx
}

func (proc *Proc) dependencyMutateWriteInputCheckWriteAddress(task *pb.Task) *prog.Prog {

	log.Logf(pb.DebugLevel, "write program : \n%s", task.WriteProgram)
	wp, err := proc.fuzzer.target.Deserialize(task.WriteProgram, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.WriteProgram: %v", err)
	}
	idx := int(task.WriteIndex)

	infoWrite := proc.execute(proc.execOptsCover, wp, ProgNormal, StatDependency)
	checkWriteAddress1 := checkAddress(task.WriteAddress, infoWrite.Calls[idx].Cover)
	if checkWriteAddress1 {
		task.CheckWriteAddress = true
		log.Logf(pb.DebugLevel, "write program could arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("write input could arrive at write address : %x", task.WriteAddress))
	} else {
		log.Logf(pb.DebugLevel, "write program could not arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("write input could not arrive at write address : %x", task.WriteAddress))
	}
	return wp
}

func (proc *Proc) dependencyMutateFinalInputCheckWriteAddress(task *pb.Task, wp *prog.Prog, idx int) *prog.Prog {
	p, err := proc.fuzzer.target.Deserialize(task.Program, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}
	var usefulCall []*prog.Call
	if int(task.WriteIndex) > len(wp.Calls) {
		log.Fatalf("dependency int(task.WriteIndex) > len(wp.Calls)")
	}
	for i, c := range wp.Calls {
		if i <= int(task.WriteIndex) {
			usefulCall = append(usefulCall, c)
		} else {
			break
		}
	}
	p.Calls = append(p.Calls[:idx], append(usefulCall, p.Calls[idx:]...)...)
	task.FinalIdx = task.Index + task.WriteIndex + 1
	task.FinalWriteIdx = uint32(idx) + task.WriteIndex
	idx = int(task.FinalIdx)

	data := p.Serialize()
	log.Logf(pb.DebugLevel, "final program : \n%s", data)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("final program : \n%s", data))

	infoFinal := proc.execute(proc.execOptsCover, p, ProgNormal, StatDependency)
	checkWriteAddress2 := checkAddress(task.WriteAddress, infoFinal.Calls[task.FinalWriteIdx].Cover)
	if checkWriteAddress2 {
		task.CheckWriteAddressFinal = true
		log.Logf(pb.DebugLevel, "final program could arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("final input could arrive at write address : %x", task.WriteAddress))
	} else {
		task.CheckWriteAddressFinal = false
		log.Logf(pb.DebugLevel, "final program could not arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("final input could not arrive at write address : %x", task.WriteAddress))
	}
	proc.fuzzer.dManager.SSendLog()
	return p
}

func (proc *Proc) dependencyMutateRemoveInputCheckWriteAddress(task *pb.Task, data []byte) *prog.Prog {
	idx := task.FinalIdx
	removeData, removeIdx := removeSameResource(data)
	task.RemoveIdx = removeIdx[idx]
	removeP, err := proc.fuzzer.target.Deserialize(removeData, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}
	log.Logf(pb.DebugLevel, "remove program : \n%s", removeData)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("remove program : \n%s", removeData))
	writeIdx := removeIdx[task.FinalWriteIdx]
	task.RemoveWriteIdx = writeIdx

	info := proc.execute(proc.execOptsCover, removeP, ProgNormal, StatDependency)
	checkWriteAddress3 := checkAddress(task.WriteAddress, info.Calls[writeIdx].Cover)
	if checkWriteAddress3 {
		task.CheckWriteAddressRemove = true
		log.Logf(pb.DebugLevel, "remove program could arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("remove input could arrive at write address : %x", task.WriteAddress))
	} else {
		task.CheckWriteAddressRemove = false
		log.Logf(pb.DebugLevel, "remove program could not arrive at write address : %d", task.WriteAddress)
		proc.fuzzer.dManager.SendLog(fmt.Sprintf("remove input could not arrive at write address : %x", task.WriteAddress))
	}
	return removeP
}

func (proc *Proc) dependencyMutateParameter(task *pb.Task, p *prog.Prog) {
	ct := proc.fuzzer.choiceTable
	idx := int(task.Index)
	for i := 0; i < 40; i++ {
		p0c := p.Clone()

		if len(task.UncoveredAddress) > 0 {
			p0c.MutateIoctl3Arg(proc.rnd, idx, ct)
		} else {
			break
		}

		infoWrite := proc.execute(proc.execOptsCover, p0c, ProgNormal, StatDependency)
		cov := cover.Cover{}
		cov.Merge(infoWrite.Calls[idx].Cover)

		for u, r := range task.UncoveredAddress {
			checkConditionAddress := checkAddressMap(r.ConditionAddress, cov)
			if !checkConditionAddress {
				continue
			}
			checkUncoveredAddress := checkAddressMap(r.Address, cov)
			if !checkUncoveredAddress {
				continue
			}
			proc.fuzzer.dManager.SendLog(fmt.Sprintf("mutate : %d cover uncovered address : %x", i, r.Address))
			r := proto.Clone(task.UncoveredAddress[u]).(*pb.RunTimeData)
			task.CoveredAddress[u] = r
			data := p0c.Serialize()
			proc.fuzzer.dManager.SendLog(fmt.Sprintf("program : \n%s", data))
			for _, c := range data {
				r.Program = append(r.Program, c)
			}
			r.TaskStatus = pb.TaskStatus_covered
			r.Idx = uint32(i)
			r.CheckAddress = true
			delete(task.UncoveredAddress, u)
		}
	}
}

func (proc *Proc) dependency(item *WorkDependency) {

	log.Logf(pb.DebugLevel, "#%v: DependencyMutate", proc.pid)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("#%v: DependencyMutate", proc.pid))

	task := item.task
	log.Logf(pb.DebugLevel, "DependencyMutate program : \n%s", task.Program)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("DependencyMutate program : \n%s", task.Program))
	log.Logf(pb.DebugLevel, "index  : %d write index : %d", task.Index, task.WriteIndex)
	proc.fuzzer.dManager.SendLog(fmt.Sprintf("index  : %d write index : %d", task.Index, task.WriteIndex))
	proc.fuzzer.dManager.SSendLog()

	wp := proc.dependencyMutateWriteInputCheckWriteAddress(task)
	if task.CheckWriteAddress {
		indexes := []int{0, int(task.Index)}
		for _, i := range indexes {
			p := proc.dependencyMutateFinalInputCheckWriteAddress(task, wp, i)
			data := p.Serialize()
			if task.CheckWriteAddressFinal {
				writeIdx := int(task.FinalWriteIdx)
				idx := int(task.FinalIdx)
				log.Logf(pb.DebugLevel, "final index  : %d final write index : %d", task.FinalIdx, task.FinalWriteIdx)
				proc.fuzzer.dManager.SendLog(fmt.Sprintf("final index  : %d final write index : %d",
					task.FinalWriteIdx, task.FinalIdx))
				proc.fuzzer.dManager.SSendLog()
				if !proc.fuzzer.comparisonTracingEnabled && item.call != -1 {
					proc.executeDependencyHintSeed(p, writeIdx)
					proc.executeDependencyHintSeed(p, idx)
				}
			}

			removeP := proc.dependencyMutateRemoveInputCheckWriteAddress(task, data)
			if task.CheckWriteAddressRemove {
				removeWriteIdx := int(task.RemoveWriteIdx)
				removeIdx := int(task.RemoveIdx)
				log.Logf(pb.DebugLevel, "remove index  : %d remove write index : %d", task.RemoveIdx, task.RemoveWriteIdx)
				proc.fuzzer.dManager.SendLog(fmt.Sprintf("remove index  : %d remove write index : %d",
					task.RemoveWriteIdx, task.RemoveIdx))
				proc.fuzzer.dManager.SSendLog()
				if removeIdx > removeWriteIdx {
					if !proc.fuzzer.comparisonTracingEnabled && item.call != -1 {
						proc.executeDependencyHintSeed(removeP, removeWriteIdx)
						proc.executeDependencyHintSeed(removeP, removeIdx)
					}
				}
			}
		}
	} else {

	}

	if len(task.UncoveredAddress) == 0 {
		task.TaskStatus = pb.TaskStatus_covered
	} else if !task.CheckWriteAddress {
		task.TaskStatus = pb.TaskStatus_unstable
	} else {
		task.TaskStatus = pb.TaskStatus_tested
	}
	tasks := &pb.Tasks{
		Name: proc.fuzzer.name,
		Task: []*pb.Task{},
	}
	tasks.Task = append(tasks.Task, task)
	proc.fuzzer.dManager.ReturnTasks(tasks)

	return
}

//
//func (proc *Proc) dependency(item *WorkDependency) {
//
//	log.Logf(pb.DebugLevel, "#%v: DependencyMutate", proc.pid)
//	proc.fuzzer.dManager.SendLog(fmt.Sprintf("#%v: DependencyMutate", proc.pid))
//
//	dependencyInput := item.dependencyInput
//	log.Logf(pb.DebugLevel, "DependencyMutate program : \n%s", dependencyInput.Program)
//	proc.fuzzer.dManager.SendLog(fmt.Sprintf("DependencyMutate program : \n%s", dependencyInput.Program))
//
//	for _, u := range dependencyInput.UncoveredAddress {
//		proc.fuzzer.dManager.SendLog(fmt.Sprintf("UncoveredAddress : %v", u.UncoveredAddress))
//		for _, wa := range u.WriteAddress {
//			log.Logf(2, "dependency data :\n%s", wa.RunTimeDate.Program)
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
//	// log.Logf(pb.DebugLevel, "repeat : %v", wa.Repeat)
//
//	//	for _, wi := range wa.WriteInput {
//	//
//	//		log.Logf(pb.DebugLevel, "write program : \n%s", wi.Program)
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
//	//		log.Logf(pb.DebugLevel, "test case with write program : \n%s", data)
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
//	log.Logf(pb.DebugLevel, "every critical condition is right but we can not arrive at write address")
//	return false
//}

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

func (fuzzer *Fuzzer) checkIsCovered(id int, address uint32) (res bool) {
	fuzzer.coverMu.RLock()
	if c, ok := fuzzer.cover[id]; ok {
		call := c.Address
		if _, ok := call[address]; !ok {
			return false
		} else {
			return true
		}
	} else {
		fuzzer.cover[id] = &pb.Call{
			Idx:     0,
			Address: make(map[uint32]uint32),
		}
		call := fuzzer.cover[id].Address
		call[address] = 0
		return false
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
			//Address: []uint32{},
		}
		newCall := input.Call[uint32(i)]

		id := p.Calls[i].Meta.ID
		if _, ok := fuzzer.cover[id]; !ok {
			fuzzer.cover[id] = &pb.Call{
				Idx:     0,
				Address: make(map[uint32]uint32),
				//Address: []uint32{},
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
		//fuzzer.dManager.SendNewInput(input)
	}

	//for _, cc := range info.Calls {
	//	log.Logf(pb.DebugLevel, "Dependency gRPC checkNewCoverage address : %v", cc.Cover)
	//}

	fuzzer.coverMu.Unlock()
	return
}
