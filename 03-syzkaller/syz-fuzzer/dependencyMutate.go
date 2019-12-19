package main

import (
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/cover"
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/hash"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/ipc"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/prog"
	"github.com/golang/protobuf/proto"
	"strings"
)

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

func (proc *Proc) getSyscall(sc *pb.FileOperationsFunction) (res *prog.Syscall) {
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
			checkConditionAddress := checkAddressInCover(r.ConditionAddress, cov)
			if !checkConditionAddress {
				continue
			}
			checkUncoveredAddress := checkAddressInCover(r.Address, cov)
			if !checkUncoveredAddress {
				continue
			}
			r := proto.Clone(task.UncoveredAddress[u]).(*pb.RunTimeData)
			task.CoveredAddress[u] = r
			data := p0c.Serialize()
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

func checkAddressInArray(Address uint32, cover []uint32) (res bool) {
	res = false
	for _, c := range cover {
		if c == Address {
			res = true
			return
		}
	}
	return
}

func checkAddressInCover(Address uint32, cover cover.Cover) (res bool) {
	res = false
	if _, ok := cover[Address]; ok {
		res = true
		return
	}
	return
}

func checkConditionInCover(condition *pb.Condition, cover cover.Cover) (res bool) {
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

func checkAddressArrayInCover(Address []uint32, cover cover.Cover) (res bool) {
	//func checkAddressArrayInCover(Address map[uint32]uint32, cover cover.Cover) (res bool) {
	res = false
	for _, a := range Address {
		if _, ok := cover[a]; ok {
			res = true
			return
		}
	}
	return
}

func checkAddressMapInCover(Address map[uint32]uint32, cover cover.Cover) (res bool) {
	//func checkAddressArrayInCover(Address map[uint32]uint32, cover cover.Cover) (res bool) {
	res = false
	for a := range Address {
		if _, ok := cover[a]; ok {
			res = true
			return
		}
	}
	return
}

func (fuzzer *Fuzzer) checkAddressIsCovered(id int, address uint32) (res bool) {
	fuzzer.coverMu.RLock()
	if c, ok := fuzzer.cover[id]; ok {
		call := c.Address
		if _, ok := call[address]; !ok {
			return false
		}
		return true
	}

	fuzzer.cover[id] = &pb.Call{
		Idx:     0,
		Address: make(map[uint32]uint32),
	}
	call := fuzzer.cover[id].Address
	call[address] = 0
	return false
}

// SendNeedInput :
func (proc *Proc) SendNeedInput(p *prog.Prog, info *ipc.ProgInfo) {
	data := p.Serialize()
	sig := hash.Hash(data)
	input := pb.Input{
		Sig:     sig.String(),
		Program: []byte{},
		Call:    make(map[uint32]*pb.Call),
		Stat:    pb.FuzzingStat_StatTriage,
	}

	for _, c := range data {
		input.Program = append(input.Program, c)
	}

	//log.Logf(2, "data :\n%s", data)
	//log.Logf(2, "input.Program :\n%s", input.Program)

	for i, c := range info.Calls {
		cc := &pb.Call{
			Idx:     uint32(i),
			Address: make(map[uint32]uint32),
		}
		input.Call[uint32(i)] = cc
		for _, a := range c.Cover {
			cc.Address[a] = 0
		}
	}

	for _, c := range p.Comments {
		i, ok := pb.FuzzingStat_value[c]
		if ok {
			input.Stat = pb.FuzzingStat(i)
		}
	}
	proc.fuzzer.dManager.SendNeedInput(&input)
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
//				address := checkAddressInArray(wa.RunTimeDate.Address, tempInfo.Calls[wa.RunTimeDate.Idx].Cover)
//				conditionAddress := checkAddressInArray(wa.RunTimeDate.ConditionAddress, tempInfo.Calls[wa.RunTimeDate.Idx].Cover)
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
//				address := checkAddressInArray(wa.RunTimeDate.Address, tempInfo.Calls[wa.RunTimeDate.Idx].Cover)
//				var cover cover.Cover
//				cover.Merge(info.Calls[wa.RunTimeDate.Idx].Cover)
//				rightBranchAddress := checkAddressArrayInCover(wa.RunTimeDate.RightBranchAddress, cover)
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
//	call := proc.getSyscall(wc)
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
//		address := checkAddressInArray(wc.RunTimeDate.Address, info.Calls[wc.RunTimeDate.Idx].Cover)
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
//	//		u.RunTimeDate.CheckAddress = checkAddressInArray(wi.WriteAddress, wpInfo.Calls[wi.Idx].Cover)
//	//
//	//		p0 := p.Clone()
//	//		p0.Splice(wp, u.Idx, programLength)
//	//
//	//		data := p0.Serialize()
//	//		log.Logf(pb.DebugLevel, "test case with write program : \n%s", data)
//	//		proc.fuzzer.dManager.SendLog(fmt.Sprintf("test case with write program : \n%s", data))
//	//
//	//		info := proc.execute(proc.execOptsCover, p0, ProgNormal, StatDependency)
//	//		u.RunTimeDate.CheckAddress = checkAddressInArray(wi.WriteAddress, info.Calls[wi.Idx].Cover)
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
