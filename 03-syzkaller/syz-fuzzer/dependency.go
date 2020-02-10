package main

import (
	"fmt"
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/hash"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/log"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/prog"
	"github.com/golang/protobuf/proto"
)

func (proc *Proc) dependency(task *pb.Task, kind pb.TaskKind) string {
	res := "dependency : " + "\n"
	res += "task hash : " + task.ComputeHash() + "\n"
	r, ok := proc.dependencyMutateCheckATask(task)
	res += r
	if ok {

		indexInsert := []int{0, int(task.Index)}
		for _, i := range indexInsert {

			res += "insert : " + "\n"
			insertTaskRunTImeData := proc.dependencyMutateInsert(task, i)
			res += proc.dependencyMutateCheck(task, insertTaskRunTImeData)
			proc.dependencyMutateArguement(task, insertTaskRunTImeData)

			res += "remove : " + "\n"
			removeTaskRunTimeData := proc.dependencyMutateRemove(task, insertTaskRunTImeData)
			res += proc.dependencyMutateCheck(task, removeTaskRunTimeData)
			proc.dependencyMutateArguement(task, insertTaskRunTImeData)

		}

	} else {

	}

	if len(task.UncoveredAddress) == 0 {
		task.TaskStatus = pb.TaskStatus_covered
	} else if !task.Check {
		task.TaskStatus = pb.TaskStatus_unstable
	} else {
		for _, TRD := range task.TaskRunTimeData {
			for _, ua := range TRD.UncoveredAddress {
				if task.TaskStatus < ua.TaskStatus {
					task.TaskStatus = ua.TaskStatus
				}
			}
		}
	}

	tasks := &pb.Tasks{
		Name:      proc.fuzzer.name,
		Kind:      kind,
		TaskMap:   map[string]*pb.Task{},
		TaskArray: []*pb.Task{},
	}
	task.Count = 1
	tasks.AddTask(task)
	proc.fuzzer.dManager.ReturnTasks(tasks)

	return res
}

func (proc *Proc) dependencyMutateArguement(task *pb.Task, taskRunTimeData *pb.TaskRunTimeData) {

	if task.Check && task.Kind == 2 {
		if taskRunTimeData.ConditionIdx > taskRunTimeData.WriteIdx {
			if !proc.fuzzer.comparisonTracingEnabled {

				Prog, err := proc.fuzzer.target.Deserialize(taskRunTimeData.Program, prog.NonStrict)
				if err != nil {
					log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
				}

				proc.executeDependencyHintSeed(Prog, int(taskRunTimeData.WriteIdx))
				proc.executeDependencyHintSeed(Prog, int(taskRunTimeData.ConditionIdx))
			}
		}
	}

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

func (proc *Proc) dependencyMutateCheckATask(task *pb.Task) (string, bool) {
	res := "dependencyMutateCheckATask : " + "\n"

	ProgWrite, err := proc.fuzzer.target.Deserialize(task.WriteProgram, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}

	idx1 := int(task.WriteIndex)
	info1 := proc.execute(proc.execOptsCover, ProgWrite, ProgNormal, StatDependency)

	ProgCondition, err := proc.fuzzer.target.Deserialize(task.Program, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}
	var temp []uint32
	idx2 := int(task.Index)
	info2 := proc.execute(proc.execOptsCover, ProgCondition, ProgNormal, StatDependency)

	for _, rTD := range task.UncoveredAddress {
		check1 := checkAddressInArray(rTD.WriteAddress, info1.Calls[idx1].Cover)
		res += fmt.Sprintf("check write address : %t : 0xffffffff%x\n", check1, rTD.WriteAddress)
		if check1 {
			rTD.CheckWrite = true
			if rTD.TaskStatus <= pb.TaskStatus_stable_write {
				rTD.TaskStatus = pb.TaskStatus_stable_write
			}
		} else {
			rTD.CheckWrite = false
			if rTD.TaskStatus <= pb.TaskStatus_unstable_write {
				rTD.TaskStatus = pb.TaskStatus_unstable_write

				if pb.CollectUnstable {
					unstableInput := &pb.UnstableInput{
						Sig:          task.Sig,
						Program:      task.Program,
						UnstablePath: []*pb.Paths{},
						Address:      map[uint32]uint32{},
					}
					unstableInput.Address[rTD.ConditionAddress] = 1 << task.Index
					paths := &pb.Paths{
						Path: map[uint32]*pb.Path{},
					}
					for i, c := range info1.Calls {
						paths.Path[uint32(i)] = &pb.Path{
							Address: c.Cover,
						}
					}
					unstableInput.UnstablePath = append(unstableInput.UnstablePath, paths)

					proc.fuzzer.dManager.SendUnstableInput(unstableInput)

				}
			}
		}

		check2 := checkAddressInArray(rTD.ConditionAddress, info2.Calls[idx2].Cover)
		res += fmt.Sprintf("check condition address : %t : 0xffffffff%x\n", check2, rTD.ConditionAddress)
		if check2 {
			rTD.CheckCondition = true
			check3 := checkAddressInArray(rTD.Address, info2.Calls[idx2].Cover)
			res += fmt.Sprintf("check uncovered address : %t : 0xffffffff%x\n", check3, rTD.Address)
			if check3 {
				rTD.CheckAddress = true
				rTD.TaskStatus = pb.TaskStatus_covered
				task.CoveredAddress[rTD.Address] = rTD
				temp = append(temp, rTD.Address)
			} else {
				rTD.CheckAddress = false
				if rTD.TaskStatus < pb.TaskStatus_stable_condition {
					rTD.TaskStatus = pb.TaskStatus_stable_condition
				}
			}
			res += fmt.Sprintf("check branch address\n")
			rTD.CheckRightBranchAddress = []bool{}
			for _, a := range rTD.RightBranchAddress {
				check4 := checkAddressInArray(a, info2.Calls[idx2].Cover)
				res += fmt.Sprintf("check branch address : %t : 0xffffffff%x\n", check4, a)
				rTD.CheckRightBranchAddress = append(rTD.CheckRightBranchAddress, check4)
			}

		} else {
			rTD.CheckCondition = false
			if rTD.TaskStatus <= pb.TaskStatus_unstable_condition {
				rTD.TaskStatus = pb.TaskStatus_unstable_condition

				if pb.CollectUnstable {
					unstableInput := &pb.UnstableInput{
						Sig:          task.Sig,
						Program:      task.Program,
						UnstablePath: []*pb.Paths{},
						Address:      map[uint32]uint32{},
					}
					unstableInput.Address[rTD.ConditionAddress] = 1 << task.Index
					paths := &pb.Paths{
						Path: map[uint32]*pb.Path{},
					}
					for i, c := range info2.Calls {
						paths.Path[uint32(i)] = &pb.Path{
							Address: c.Cover,
						}
					}
					unstableInput.UnstablePath = append(unstableInput.UnstablePath, paths)
					proc.fuzzer.dManager.SendUnstableInput(unstableInput)

				}
			}
		}

		task.Check = task.Check || (check1 && check2)
	}

	for _, ua := range temp {
		delete(task.UncoveredAddress, ua)
	}

	return res, task.Check
}

func (proc *Proc) dependencyMutateInsert(task *pb.Task, idx int) *pb.TaskRunTimeData {

	ProgWrite, err := proc.fuzzer.target.Deserialize(task.WriteProgram, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}

	ProgCondition, err := proc.fuzzer.target.Deserialize(task.Program, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}

	var usefulSyscall []*prog.Call
	if int(task.WriteIndex) > len(ProgWrite.Calls) {
		log.Fatalf("dependency int(task.WriteIndex) > len(ProgWrite.Calls)")
	}
	for i, c := range ProgWrite.Calls {
		if i <= int(task.WriteIndex) {
			usefulSyscall = append(usefulSyscall, c)
		} else {
			break
		}
	}
	p := ProgCondition.Clone()
	p.Calls = append(p.Calls[:idx], append(usefulSyscall, p.Calls[idx:]...)...)
	data := p.Serialize()

	insertTaskRunTImeData := &pb.TaskRunTimeData{
		Hash:             hash.String(data),
		Program:          nil,
		WriteIdx:         uint32(idx) + task.WriteIndex,
		ConditionIdx:     task.Index + task.WriteIndex + 1,
		UncoveredAddress: map[uint32]*pb.RunTimeData{},
		CoveredAddress:   map[uint32]*pb.RunTimeData{},
	}

	for _, c := range data {
		insertTaskRunTImeData.Program = append(insertTaskRunTImeData.Program, c)
	}

	for ua, r := range task.UncoveredAddress {
		insertTaskRunTImeData.UncoveredAddress[ua] = proto.Clone(r).(*pb.RunTimeData)
	}

	for _, r := range insertTaskRunTImeData.UncoveredAddress {
		r.TaskStatus = pb.TaskStatus_untested
		for _, c := range data {
			r.Program = append(r.Program, c)
		}
		r.Idx = insertTaskRunTImeData.ConditionIdx
		r.CheckWrite = false
		r.CheckAddress = false
		r.CheckCondition = false
	}

	task.TaskRunTimeData = append(task.TaskRunTimeData, insertTaskRunTImeData)

	return insertTaskRunTImeData
}

func (proc *Proc) dependencyMutateRemove(task *pb.Task, taskRunTimeData *pb.TaskRunTimeData) *pb.TaskRunTimeData {

	removeData, removeIdx := removeSameResource(taskRunTimeData.Program)

	removeTaskRunTimeData := &pb.TaskRunTimeData{
		Hash:             hash.String(removeData),
		Program:          []byte{},
		WriteIdx:         removeIdx[taskRunTimeData.WriteIdx],
		ConditionIdx:     removeIdx[taskRunTimeData.ConditionIdx],
		UncoveredAddress: map[uint32]*pb.RunTimeData{},
		CoveredAddress:   map[uint32]*pb.RunTimeData{},
	}

	for _, c := range removeData {
		removeTaskRunTimeData.Program = append(removeTaskRunTimeData.Program, c)
	}

	for ua, r := range task.UncoveredAddress {
		removeTaskRunTimeData.UncoveredAddress[ua] = proto.Clone(r).(*pb.RunTimeData)
	}

	for _, r := range removeTaskRunTimeData.UncoveredAddress {
		r.TaskStatus = pb.TaskStatus_untested
		for _, c := range removeData {
			r.Program = append(r.Program, c)
		}
		r.Idx = removeTaskRunTimeData.ConditionIdx
		r.CheckWrite = false
		r.CheckAddress = false
		r.CheckCondition = false
	}

	task.TaskRunTimeData = append(task.TaskRunTimeData, removeTaskRunTimeData)

	return removeTaskRunTimeData
}

// kind 1: final 2: remove
func (proc *Proc) dependencyMutateCheck(task *pb.Task, taskRunTimeData *pb.TaskRunTimeData) string {
	res := "dependencyMutateCheck : " + "\n"
	Prog, err := proc.fuzzer.target.Deserialize(taskRunTimeData.Program, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}

	info := proc.execute(proc.execOptsCover, Prog, ProgNormal, StatDependency)

	var temp []uint32
	for ua, r := range taskRunTimeData.UncoveredAddress {

		check1 := checkAddressInArray(r.WriteAddress, info.Calls[taskRunTimeData.WriteIdx].Cover)
		res += fmt.Sprintf("check write address : %t : 0xffffffff%x\n", check1, r.WriteAddress)
		if check1 {
			r.CheckWrite = true
			check2 := checkAddressInArray(r.ConditionAddress, info.Calls[taskRunTimeData.ConditionIdx].Cover)
			res += fmt.Sprintf("check condition address : %t : 0xffffffff%x\n", check2, r.ConditionAddress)
			if check2 {
				r.CheckCondition = true
				check3 := checkAddressInArray(ua, info.Calls[taskRunTimeData.ConditionIdx].Cover)
				res += fmt.Sprintf("check uncovered address : %t : 0xffffffff%x\n", check3, ua)
				if check3 {
					r.CheckAddress = true
					r.TaskStatus = pb.TaskStatus_covered
					taskRunTimeData.CoveredAddress[ua] = r
					temp = append(temp, ua)
				} else {
					r.CheckAddress = false
					if r.TaskStatus < pb.TaskStatus_tested {
						r.TaskStatus = pb.TaskStatus_tested
					}
				}

				r.CheckRightBranchAddress = []bool{}
				for _, a := range r.RightBranchAddress {
					check4 := checkAddressInArray(a, info.Calls[taskRunTimeData.ConditionIdx].Cover)
					res += fmt.Sprintf("check branch address : %t : 0xffffffff%x\n", check4, a)
					r.CheckRightBranchAddress = append(r.CheckRightBranchAddress, check4)
				}

			} else {
				r.CheckCondition = false
				if r.TaskStatus < pb.TaskStatus_unstable_insert_condition {
					r.TaskStatus = pb.TaskStatus_unstable_insert_condition

				}
			}
		} else {
			r.CheckWrite = false
			if r.TaskStatus < pb.TaskStatus_unstable_insert_write {
				r.TaskStatus = pb.TaskStatus_unstable_insert_write
			}
		}
	}

	for _, ua := range temp {
		delete(taskRunTimeData.UncoveredAddress, ua)
	}
	for ua, r := range taskRunTimeData.CoveredAddress {
		if _, ok := task.UncoveredAddress[ua]; ok {
			delete(task.UncoveredAddress, ua)
		}
		if _, ok := task.CoveredAddress[ua]; ok {

		} else {
			task.CoveredAddress[ua] = r
		}
	}

	return res
}

func (proc *Proc) dependencyBoot(item *WorkBoot) {
	task := item.task
	p, err := proc.fuzzer.target.Deserialize(task.WriteProgram, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.WriteProgram: %v", err)
	}
	idx := int(task.Index)
	var index []int
	for i := 0; i < 32; i++ {
		if (1<<uint(i))&task.WriteIndex > 0 {
			index = append(index, i)
		}
	}
	l := len(p.Calls)
	for i := idx + 1; i < l; i++ {
		p.RemoveCall(i)
	}
	for ii, i := range index {
		if i >= l-1 {
			break
		}
		p.RemoveCall(i - ii)
	}

	infoWrite := proc.execute(proc.execOptsCover, p, ProgNormal, StatDependency)
	for UncoveredAddress := range task.UncoveredAddress {
		checkUncoveredAddress := checkAddressInArray(UncoveredAddress, infoWrite.Calls[idx].Cover)
		if checkUncoveredAddress {
			task.CoveredAddress[UncoveredAddress] = task.UncoveredAddress[UncoveredAddress]
		}
	}
	data := p.Serialize()
	for _, c := range data {
		task.WriteProgram = append(task.WriteProgram, c)
	}

	for address := range task.CoveredAddress {
		delete(task.UncoveredAddress, address)
	}

	if len(task.CoveredAddress) != 0 {
		task.TaskStatus = pb.TaskStatus_covered
		input := pb.Input{
			Sig:     task.Sig,
			Program: []byte{},
			Call:    make(map[uint32]*pb.Call),
			Stat:    pb.FuzzingStat_StatTriage,
		}

		for _, c := range data {
			input.Program = append(input.Program, c)
		}
		if item.call != -1 {
			cc := &pb.Call{
				Idx:     uint32(item.call),
				Address: make(map[uint32]uint32),
			}
			input.Call[uint32(item.call)] = cc
			for _, a := range infoWrite.Calls[idx].Cover {
				cc.Address[a] = 0
			}
		}
		input.Stat = pb.FuzzingStat_StatDependencyBoot
		proc.fuzzer.dManager.SendBootInput(&input)
	} else {
		task.TaskStatus = pb.TaskStatus_tested
	}
	tasks := &pb.Tasks{
		Name:      proc.fuzzer.name,
		Kind:      pb.TaskKind_Boot,
		TaskMap:   map[string]*pb.Task{},
		TaskArray: []*pb.Task{},
	}
	tasks.AddTask(task)
	proc.fuzzer.dManager.ReturnTasks(tasks)

	return
}

func (proc *Proc) checkInput(input *pb.Input) {
	res := ""
	proc.fuzzer.dManager.MuDependency.Lock()
	ua, tasks, r := proc.fuzzer.dManager.DataDependency.GetTaskByInput(input)
	proc.fuzzer.dManager.MuDependency.Unlock()
	res += r
	for _, t := range tasks {
		// mutate the argument
		t.Kind = 2
		res += proc.dependency(t, pb.TaskKind_Ckeck)
	}

	proc.fuzzer.dManager.SendLog(res)
	proc.fuzzer.dManager.SSendLog()

	proc.fuzzer.dManager.MuDependency.Lock()
	if len(tasks) > 0 {
		if _, ok := proc.fuzzer.dManager.DataDependency.UncoveredAddress[ua.UncoveredAddress]; ok {
			delete(proc.fuzzer.dManager.DataDependency.UncoveredAddress, ua.UncoveredAddress)
		}
	}
	proc.fuzzer.dManager.MuDependency.Unlock()

}
