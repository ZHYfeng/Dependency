package main

import (
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/hash"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/log"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/prog"
	"github.com/golang/protobuf/proto"
)

func (proc *Proc) dependency(item *WorkDependency) {

	task := item.task

	if proc.dependencyMutateCheckATask(task) {

		indexInsert := []int{0, int(task.Index)}
		for _, i := range indexInsert {

			insertTaskRunTImeData := proc.dependencyMutateInsert(task, i)
			proc.dependencyMutateCheck(task, insertTaskRunTImeData)
			proc.dependencyMutateArguement(task, insertTaskRunTImeData)

			removeTaskRunTimeData := proc.dependencyMutateRemove(task, insertTaskRunTImeData)
			proc.dependencyMutateCheck(task, removeTaskRunTimeData)
			proc.dependencyMutateArguement(task, insertTaskRunTImeData)

		}

	} else {

	}

	if len(task.UncoveredAddress) == 0 {
		task.TaskStatus = pb.TaskStatus_covered
	} else if !task.CheckWriteAddress {
		task.TaskStatus = pb.TaskStatus_unstable
	} else {
		task.TaskStatus = pb.TaskStatus_tested
		for _, ua := range task.UncoveredAddress {
			if ua.TaskStatus == pb.TaskStatus_unstable {
				task.TaskStatus = pb.TaskStatus_unstable
			}
		}
	}
	tasks := &pb.Tasks{
		Name:      proc.fuzzer.name,
		Kind:      pb.TaskKind_Normal,
		TaskMap:   map[string]*pb.Task{},
		TaskArray: []*pb.Task{},
	}
	task.Count = 1
	tasks.AddTask(task)
	proc.fuzzer.dManager.ReturnTasks(tasks)

	return
}

func (proc *Proc) dependencyMutateArguement(task *pb.Task, taskRunTimeData *pb.TaskRunTimeData) {

	if task.CheckWriteAddress && task.Kind == 2 {
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

func (proc *Proc) dependencyMutateCheckATask(task *pb.Task) bool {

	ProgWrite, err := proc.fuzzer.target.Deserialize(task.WriteProgram, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}

	idx := int(task.WriteIndex)
	info := proc.execute(proc.execOptsCover, ProgWrite, ProgNormal, StatDependency)
	checkWriteAddress1 := checkAddressInArray(task.WriteAddress, info.Calls[idx].Cover)
	if checkWriteAddress1 {
		task.CheckWriteAddress = true
	} else {
		task.CheckWriteAddress = false
	}

	ProgCondition, err := proc.fuzzer.target.Deserialize(task.Program, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}
	var temp []uint32
	idx = int(task.Index)
	info = proc.execute(proc.execOptsCover, ProgCondition, ProgNormal, StatDependency)
	for ua, r := range task.UncoveredAddress {
		if checkAddressInArray(r.ConditionAddress, info.Calls[idx].Cover) {
			r.CheckCondition = true
			if checkAddressInArray(ua, info.Calls[idx].Cover) {
				r.CheckAddress = true
				r.TaskStatus = pb.TaskStatus_covered
				task.CoveredAddress[ua] = r
				temp = append(temp, ua)
			} else {
				r.CheckAddress = false
				if r.TaskStatus < pb.TaskStatus_tested {
					r.TaskStatus = pb.TaskStatus_tested
				}
			}
		} else {
			r.CheckCondition = false
			if r.TaskStatus < pb.TaskStatus_unstable {
				r.TaskStatus = pb.TaskStatus_unstable
			}
		}
	}

	for _, ua := range temp {
		delete(task.UncoveredAddress, ua)
	}

	return task.CheckWriteAddress
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
		Hash:              hash.String(data),
		Program:           nil,
		WriteIdx:          uint32(idx) + task.WriteIndex,
		ConditionIdx:      task.Index + task.WriteIndex + 1,
		CheckWriteAddress: false,
		UncoveredAddress:  map[uint32]*pb.RunTimeData{},
		CoveredAddress:    map[uint32]*pb.RunTimeData{},
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
		r.CheckAddress = false
		r.CheckCondition = false
	}

	task.TaskRunTimeData = append(task.TaskRunTimeData, insertTaskRunTImeData)

	return insertTaskRunTImeData
}

func (proc *Proc) dependencyMutateRemove(task *pb.Task, taskRunTimeData *pb.TaskRunTimeData) *pb.TaskRunTimeData {

	removeData, removeIdx := removeSameResource(taskRunTimeData.Program)

	removeTaskRunTimeData := &pb.TaskRunTimeData{
		Hash:              hash.String(removeData),
		Program:           []byte{},
		WriteIdx:          removeIdx[taskRunTimeData.WriteIdx],
		ConditionIdx:      removeIdx[taskRunTimeData.ConditionIdx],
		CheckWriteAddress: false,
		UncoveredAddress:  map[uint32]*pb.RunTimeData{},
		CoveredAddress:    map[uint32]*pb.RunTimeData{},
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
		r.CheckAddress = false
		r.CheckCondition = false
	}

	task.TaskRunTimeData = append(task.TaskRunTimeData, removeTaskRunTimeData)

	return removeTaskRunTimeData
}

// kind 1: final 2: remove
func (proc *Proc) dependencyMutateCheck(task *pb.Task, taskRunTimeData *pb.TaskRunTimeData) {

	Prog, err := proc.fuzzer.target.Deserialize(taskRunTimeData.Program, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}

	info := proc.execute(proc.execOptsCover, Prog, ProgNormal, StatDependency)
	checkWriteAddress2 := checkAddressInArray(task.WriteAddress, info.Calls[taskRunTimeData.WriteIdx].Cover)
	if checkWriteAddress2 {
		taskRunTimeData.CheckWriteAddress = true
	} else {
		taskRunTimeData.CheckWriteAddress = false
	}

	var temp []uint32
	for ua, r := range taskRunTimeData.UncoveredAddress {
		if checkAddressInArray(r.ConditionAddress, info.Calls[taskRunTimeData.ConditionIdx].Cover) {
			r.CheckCondition = true
			if checkAddressInArray(ua, info.Calls[taskRunTimeData.ConditionIdx].Cover) {
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
		} else {
			r.CheckCondition = false
			if r.TaskStatus < pb.TaskStatus_unstable {
				r.TaskStatus = pb.TaskStatus_unstable
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

	return
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
