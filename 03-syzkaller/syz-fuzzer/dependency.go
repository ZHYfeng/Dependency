package main

import (
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/log"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/prog"
)

func (proc *Proc) dependency(item *WorkDependency) {

	log.Logf(pb.DebugLevel, "#%v: DependencyMutate", proc.pid)
	task := item.task
	log.Logf(pb.DebugLevel, "DependencyMutate program : \n%s", task.Program)
	log.Logf(pb.DebugLevel, "index  : %d write index : %d", task.Index, task.WriteIndex)

	writeProg, err := proc.fuzzer.target.Deserialize(task.WriteProgram, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.WriteProgram: %v", err)
	}
	proc.dependencyMutateCheckWriteAddress(task, writeProg)

	if task.CheckWriteAddress {

		Prog, err := proc.fuzzer.target.Deserialize(task.Program, prog.NonStrict)
		if err != nil {
			log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
		}

		indexInsert := []int{0, int(task.Index)}
		for _, i := range indexInsert {

			tempProg := proc.dependencyMutateInsert(task, Prog, writeProg, i)
			proc.dependencyMutateCheck(task, tempProg)

			writeIndex := int(task.FinalWriteIdx)
			index := int(task.FinalIdx)
			proc.dependencyMutateArguement(task, tempProg, writeIndex, index)

			data := tempProg.Serialize()
			removeProg := proc.dependencyMutateRemove(task, data)
			proc.dependencyMutateCheck(task, removeProg)

			removeWriteIdx := int(task.RemoveWriteIdx)
			removeIdx := int(task.RemoveIdx)
			proc.dependencyMutateArguement(task, removeProg, removeWriteIdx, removeIdx)

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
		Name:  proc.fuzzer.name,
		Kind:  pb.TaskKind_Normal,
		Task:  map[string]*pb.Task{},
		Tasks: []*pb.Task{},
	}
	tasks.Tasks = append(tasks.Tasks, task)
	proc.fuzzer.dManager.ReturnTasks(tasks)

	return
}

func (proc *Proc) dependencyMutateArguement(task *pb.Task, Prog *prog.Prog, writeIndex int, index int) {

	if task.CheckWriteAddressRemove && task.Kind == 2 {
		if index > writeIndex {
			if !proc.fuzzer.comparisonTracingEnabled {
				proc.executeDependencyHintSeed(Prog, writeIndex)
				proc.executeDependencyHintSeed(Prog, index)
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

func (proc *Proc) dependencyMutateCheckWriteAddress(task *pb.Task, writeProg *prog.Prog) {
	log.Logf(pb.DebugLevel, "write program : \n%s", task.WriteProgram)
	idx := int(task.WriteIndex)
	info := proc.execute(proc.execOptsCover, writeProg, ProgNormal, StatDependency)
	checkWriteAddress1 := checkAddressInArray(task.WriteAddress, info.Calls[idx].Cover)
	if checkWriteAddress1 {
		task.CheckWriteAddress = true
		log.Logf(pb.DebugLevel, "write program could arrive at write address : %d", task.WriteAddress)
	} else {
		log.Logf(pb.DebugLevel, "write program could not arrive at write address : %d", task.WriteAddress)
	}
	return
}

func (proc *Proc) dependencyMutateInsert(task *pb.Task, Prog *prog.Prog, writeProg *prog.Prog, idx int) *prog.Prog {

	var usefulSyscall []*prog.Call
	if int(task.WriteIndex) > len(writeProg.Calls) {
		log.Fatalf("dependency int(task.WriteIndex) > len(wp.Calls)")
	}
	for i, c := range writeProg.Calls {
		if i <= int(task.WriteIndex) {
			usefulSyscall = append(usefulSyscall, c)
		} else {
			break
		}
	}

	p := Prog.Clone()
	p.Calls = append(p.Calls[:idx], append(usefulSyscall, p.Calls[idx:]...)...)
	task.FinalIdx = task.Index + task.WriteIndex + 1
	task.FinalWriteIdx = uint32(idx) + task.WriteIndex

	data := p.Serialize()
	log.Logf(pb.DebugLevel, "final program : \n%s", data)
	log.Logf(pb.DebugLevel, "final index  : %d final write index : %d", task.FinalIdx, task.FinalWriteIdx)

	return p
}

func (proc *Proc) dependencyMutateCheck(task *pb.Task, Prog *prog.Prog) {

	infoFinal := proc.execute(proc.execOptsCover, Prog, ProgNormal, StatDependency)
	checkWriteAddress2 := checkAddressInArray(task.WriteAddress, infoFinal.Calls[task.FinalWriteIdx].Cover)
	if checkWriteAddress2 {
		task.CheckWriteAddressFinal = true
		log.Logf(pb.DebugLevel, "final program could arrive at write address : %d", task.WriteAddress)
	} else {
		task.CheckWriteAddressFinal = false
		log.Logf(pb.DebugLevel, "final program could not arrive at write address : %d", task.WriteAddress)
	}
}

func (proc *Proc) dependencyMutateRemove(task *pb.Task, data []byte) *prog.Prog {

	idx := task.FinalIdx
	removeData, removeIdx := removeSameResource(data)
	task.RemoveIdx = removeIdx[idx]
	removeProg, err := proc.fuzzer.target.Deserialize(removeData, prog.NonStrict)
	if err != nil {
		log.Fatalf("dependency failed to deserialize program from task.Program: %v", err)
	}
	log.Logf(pb.DebugLevel, "remove program : \n%s", removeData)
	writeIdx := removeIdx[task.FinalWriteIdx]
	task.RemoveWriteIdx = writeIdx
	log.Logf(pb.DebugLevel, "remove index  : %d remove write index : %d", task.RemoveIdx, task.RemoveWriteIdx)
	return removeProg
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
		Name:  proc.fuzzer.name,
		Kind:  pb.TaskKind_Boot,
		Task:  map[string]*pb.Task{},
		Tasks: []*pb.Task{},
	}
	tasks.Tasks = append(tasks.Tasks, task)
	proc.fuzzer.dManager.ReturnTasks(tasks)

	return
}
