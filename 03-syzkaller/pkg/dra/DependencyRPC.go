package dra

import (
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/log"
	"github.com/golang/protobuf/proto"
)

func (m *Statistic) mergeStatistic(d *Statistic) {

	if m.Name != d.Name {
		log.Fatalf("MergeStatistic with error name")
		return
	}

	m.ExecuteNum = m.ExecuteNum + d.ExecuteNum
	m.Time = m.Time + d.Time
	m.NewTestCaseNum = m.NewTestCaseNum + d.NewTestCaseNum
	m.NewAddressNum = m.NewAddressNum + d.NewAddressNum

	return
}

func (m *Input) mergeInput(d *Input) {
	for i, u := range d.Call {
		var call *Call
		if c, ok := m.Call[i]; ok {
			call = c
		} else {
			call = &Call{
				Address: make(map[uint32]uint32),
				Idx:     u.Idx,
			}
			m.Call[i] = call
		}

		if(call.Address == nil) {
			// templog := "debug mergeInput :\n"
			// templog += string(m.Program) + "\n" + string(d.Program) + "\n"
			// templog += "index : " + strconv.FormatInt(int64(i), 10) + "\n"
			// f, _ := os.OpenFile("./debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			// _, _ = f.WriteString(string(templog))
			// _ = f.Close()
			call.Address = make(map[uint32]uint32)
		}
		
		for a := range u.Address {
			call.Address[a] = 0
		}
	}

	for i, c := range d.UncoveredAddress {
		if index, ok := m.UncoveredAddress[i]; ok {
			m.UncoveredAddress[i] = index | c
		} else {
			if m.UncoveredAddress == nil {
				m.UncoveredAddress = map[uint32]uint32{}
			}
			m.UncoveredAddress[i] = c
		}
	}

	for i, c := range d.WriteAddress {
		if index, ok := m.WriteAddress[i]; ok {
			m.WriteAddress[i] = index | c
		} else {
			if m.WriteAddress == nil {
				m.WriteAddress = map[uint32]uint32{}
			}
			m.WriteAddress[i] = c
		}
	}

	return
}

func (m *UncoveredAddress) mergeUncoveredAddress(d *UncoveredAddress) {
	for i, c := range d.Input {
		if index, ok := m.Input[i]; ok {
			m.Input[i] = index | c
		} else {
			if m.Input == nil {
				m.Input = map[string]uint32{}
			}
			m.Input[i] = c
		}
	}

	for i, c := range d.WriteAddress {
		if _, ok := m.WriteAddress[i]; ok {

		} else {
			if m.WriteAddress == nil {
				m.WriteAddress = map[uint32]*WriteAddressAttributes{}
			}
			m.WriteAddress[i] = proto.Clone(c).(*WriteAddressAttributes)
		}
	}

	return
}

func (m *FileOperations) mergeFileOperations(d *FileOperations) {
	if m.Name != d.Name {
		log.Fatalf("mergeFileOperations with error name")
		return
	}

	for i, c := range d.FileOperationsFunction {
		if _, ok := m.FileOperationsFunction[i]; ok {
			proto.Merge(m.FileOperationsFunction[i], c)
		} else {
			if m.FileOperationsFunction == nil {
				m.FileOperationsFunction = map[int64]*FileOperationsFunction{}
			}
			m.FileOperationsFunction[i] = c
		}
	}
}

func (m *WriteAddress) mergeWriteAddress(d *WriteAddress) {

	for i, c := range d.UncoveredAddress {
		if _, ok := m.UncoveredAddress[i]; ok {

		} else {
			if m.UncoveredAddress == nil {
				m.UncoveredAddress = map[uint32]*WriteAddressAttributes{}
			}
			m.UncoveredAddress[i] = proto.Clone(c).(*WriteAddressAttributes)
		}
	}

	for i, c := range d.FileOperationsFunction {
		if _, ok := m.FileOperationsFunction[i]; ok {
			m.FileOperationsFunction[i] = m.FileOperationsFunction[i] | c
		} else {
			if m.FileOperationsFunction == nil {
				m.FileOperationsFunction = map[string]uint32{}
			}
			m.FileOperationsFunction[i] = c
		}
	}

	for i, c := range d.Input {
		if index, ok := m.Input[i]; ok {
			m.Input[i] = index | c
		} else {
			if m.Input == nil {
				m.Input = map[string]uint32{}
			}
			m.Input[i] = c
		}
	}

	return
}

func (m *RunTimeData) mergeRunTimeData(d *RunTimeData) {
	if d == nil {
		return
	}

	return
}

func (ss Server) increasePriority(m *Task) {
	if t, ok := ss.corpusDependency.Tasks.Task[m.getHash()]; ok {
		t.increasePriority()
	}
}

func (ss Server) reducePriority(m *Task) {
	if t, ok := ss.corpusDependency.Tasks.Task[m.getHash()]; ok {
		t.reducePriority()
	}
}

func (m *Task) increasePriority() {
	m.Priority++
}

func (m *Task) reducePriority() {
	m.Priority--
}

func (m *Task) mergeTask(s *Task) {
	if m.CoveredAddress == nil {
		m.CoveredAddress = map[uint32]*RunTimeData{}
	}
	for u, p := range s.CoveredAddress {
		m.CoveredAddress[u] = proto.Clone(p).(*RunTimeData)
	}

	for u := range m.UncoveredAddress {
		_, ok := m.CoveredAddress[u]
		if ok {
			delete(m.UncoveredAddress, u)
		}
	}
	if m.TaskStatus == TaskStatus_testing {
		m.TaskStatus = s.TaskStatus
	}

	m.CheckWriteAddress = s.CheckWriteAddress || m.CheckWriteAddress
	m.CheckWriteAddressFinal = s.CheckWriteAddressFinal || m.CheckWriteAddressFinal
	m.CheckWriteAddressRemove = s.CheckWriteAddressRemove || m.CheckWriteAddressRemove

	return
}

func (ss Server) pickTask(name string) *Tasks {
	var tasks *Tasks
	f, ok := ss.fuzzers[name]
	if ok {
		f.taskMu.Lock()
		if len(f.highTasks.Task) > 0 {
			last := len(f.highTasks.Task)
			if last > taskNum {
				last = taskNum
			}
			tasks = f.highTasks.pop(last)
			tasks.Kind = TaskKind_High
		} else {
			last := len(f.newTask.Task)
			if last > taskNum {
				last = taskNum
			}
			tasks = f.newTask.pop(last)
		}
		f.taskMu.Unlock()
	}
	return tasks
}

// pickBootTask : pick one task once
func (ss Server) pickBootTask(name string) *Tasks {
	var tasks *Tasks
	f, ok := ss.fuzzers[name]
	if ok {
		f.taskMu.Lock()
		last := len(f.bootTasks.Task)
		tasks = f.bootTasks.pop(last)
		tasks.Kind = TaskKind_Boot
		f.taskMu.Unlock()
	}
	return tasks
}

func (ss *Server) addNewInput(s *Input) {
	if i, ok := ss.corpusDependency.NewInput[s.Sig]; ok {
		i.mergeInput(s)
	} else {
		ss.corpusDependency.NewInput[s.Sig] = s
	}

	return
}

func (ss *Server) addInput(s *Input) {
	if i, ok := ss.corpusDependency.Input[s.Sig]; ok {
		i.mergeInput(s)
	} else {
		ss.corpusDependency.Input[s.Sig] = s
	}

	ss.addWriteAddressMapInput(s)
	ss.addUncoveredAddressMapInput(s)

	// ss.corpusDependency.Input[s.Sig].Call = make(map[uint32]*Call)
	return
}

func (ss *Server) addWriteAddressMapInput(s *Input) {
	sig := s.Sig
	for index, call := range s.Call {
		indexBits := uint32(1 << index)
		for a := range call.Address {
			if wa, ok := ss.corpusDependency.WriteAddress[a]; ok {
				var usefulIndexBits uint32
				waIndex, ok := wa.Input[sig]
				if ok {
					if (waIndex|indexBits)^waIndex > 0 {
						usefulIndexBits = (waIndex | indexBits) ^ waIndex
						wa.Input[sig] = waIndex | indexBits
					}
				} else {
					if wa.Input == nil {
						wa.Input = map[string]uint32{}
					}
					usefulIndexBits = indexBits
					wa.Input[sig] = indexBits
				}
				ss.addWriteAddressTask(wa, sig, usefulIndexBits)
				input, ok := ss.corpusDependency.Input[sig]
				if ok {
					if iIndex, ok := input.WriteAddress[a]; ok {
						input.WriteAddress[a] = iIndex | indexBits
					} else {
						if input.WriteAddress == nil {
							input.WriteAddress = map[uint32]uint32{}
						}
						input.WriteAddress[a] = indexBits
					}
				}
			} else {
			}
		}
	}
	return
}

func (ss *Server) addUncoveredAddressMapInput(s *Input) {
	sig := s.Sig
	for u1, i1 := range s.UncoveredAddress {
		if u2, ok := ss.corpusDependency.UncoveredAddress[u1]; ok {
			if i2, ok := u2.Input[sig]; ok {
				u2.Input[sig] = i2 | i1
			} else {
				if u2.Input == nil {
					u2.Input = map[string]uint32{}
				}
				u2.Input[sig] = i1
			}
		}
	}
	return
}

func (ss *Server) checkUncoveredAddress(uncoveredAddress uint32) bool {
	_, ok := ss.corpusDependency.UncoveredAddress[uncoveredAddress]
	if !ok {
		return false
	}
	ss.deleteUncoveredAddress(uncoveredAddress)
	return true
}

func (ss *Server) deleteUncoveredAddress(uncoveredAddress uint32) {
	u, ok := ss.corpusDependency.UncoveredAddress[uncoveredAddress]
	if !ok {
		return
	}

	for sig := range u.Input {
		input, ok := ss.corpusDependency.Input[sig]
		if !ok {
			log.Fatalf("deleteUncoveredAddress not find sig")
			continue
		} else {
		}
		_, ok1 := input.UncoveredAddress[uncoveredAddress]
		if !ok1 {
			log.Fatalf("deleteUncoveredAddress input not find uncoveredAddress")
		} else {
			delete(input.UncoveredAddress, uncoveredAddress)

		}
	}

	for wa := range u.WriteAddress {
		waa, ok := ss.corpusDependency.WriteAddress[wa]
		if !ok {
			log.Fatalf("deleteUncoveredAddress not find wa")
			continue
		} else {
		}
		_, ok1 := waa.UncoveredAddress[uncoveredAddress]
		if !ok1 {
			log.Fatalf("deleteUncoveredAddress write address not find uncoveredAddress")
		} else {
			delete(waa.UncoveredAddress, uncoveredAddress)
		}
	}

	ss.corpusDependency.CoveredAddress[uncoveredAddress] = u
	delete(ss.corpusDependency.UncoveredAddress, uncoveredAddress)

	return
}

func (ss *Server) addCoveredAddress(input *Input) {
	var isDependency uint32
	if input.Stat == FuzzingStat_StatDependency || input.Stat == FuzzingStat_StatDependencyBoot {
		isDependency = 1
	} else {
		isDependency = 0
	}
	var newAddressNum uint64
	newAddressNum = 0
	var aa []uint32
	for _, call := range input.Call {
		for a := range call.Address {
			_, ok := ss.stat.Coverage.Coverage[a]
			if ok {

			} else {
				newAddressNum++
				ss.stat.Coverage.Coverage[a] = isDependency
			}
			aa = append(aa, a)
		}
	}
	t := time.Now()
	elapsed := t.Sub(ss.timeStart)
	ss.stat.Coverage.Time = append(ss.stat.Coverage.Time, &Time{
		Time: elapsed.Seconds(),
		Num:  int64(len(ss.stat.Coverage.Coverage)),
	})
	log.Logf(DebugLevel, "(ss Server) addCoveredAddress : checkUncoveredAddress")
	for _, a := range aa {
		ss.checkUncoveredAddress(a)
	}
	s, ok := ss.stat.Stat[int32(input.Stat)]
	if ok {
		s.NewTestCaseNum++
		s.NewAddressNum = s.NewAddressNum + newAddressNum
	} else {
		ss.stat.Stat[int32(input.Stat)] = &Statistic{
			Name:           input.Stat,
			ExecuteNum:     0,
			Time:           0,
			NewTestCaseNum: 1,
			NewAddressNum:  newAddressNum,
		}
	}

	if newAddressNum > 100 {
		usefulInput := &UsefulInput{
			Input:      input,
			Time:       elapsed.Seconds(),
			Num:        newAddressNum,
			NewAddress: aa,
		}
		ss.stat.UsefulInput = append(ss.stat.UsefulInput, usefulInput)
	}

	return
}

func (ss *Server) addUncoveredAddress(s *UncoveredAddress) {
	_, ok := ss.stat.Coverage.Coverage[s.UncoveredAddress]
	if ok {
		return
	}

	if i, ok := ss.corpusDependency.UncoveredAddress[s.UncoveredAddress]; ok {
		i.mergeUncoveredAddress(s)
	} else {
		ss.corpusDependency.UncoveredAddress[s.UncoveredAddress] = s
	}
	ss.addWriteAddressMapUncoveredAddress(s)

	return
}

func (ss *Server) addWriteAddressMapUncoveredAddress(s *UncoveredAddress) {
	uncoveredAddress := s.UncoveredAddress
	for w1, w3 := range s.WriteAddress {
		if w2, ok := ss.corpusDependency.WriteAddress[w1]; ok {
			w2.UncoveredAddress[uncoveredAddress] = w3
		}
	}
	return
}

func (ss *Server) addWriteAddress(s *WriteAddress) {
	if i, ok := ss.corpusDependency.WriteAddress[s.WriteAddress]; ok {
		i.mergeWriteAddress(s)
	} else {
		ss.corpusDependency.WriteAddress[s.WriteAddress] = s
	}

	for sig, indexBits1 := range s.Input {
		waInput, ok := ss.corpusDependency.Input[sig]
		if ok {
			indexBits2, ok1 := waInput.WriteAddress[s.WriteAddress]
			if ok1 {
				waInput.WriteAddress[s.WriteAddress] = indexBits2 | indexBits1
			} else {
				if waInput.WriteAddress == nil {
					waInput.WriteAddress = map[uint32]uint32{}
				}
				waInput.WriteAddress[s.WriteAddress] = indexBits1
			}
		} else {
			log.Logf(0, "addWriteAddress not find sig")
		}
	}
}

func (ss *Server) addInputTask(d *Input) {
	sig := d.Sig
	for u, inputIndexBits := range d.UncoveredAddress {
		ua, ok := ss.corpusDependency.UncoveredAddress[u]
		if !ok {
			return
		}
		for w := range ua.WriteAddress {
			wa, ok := ss.corpusDependency.WriteAddress[w]
			if !ok {
				return
			}
			for writeSig, indexBits := range wa.Input {
				ss.addTasks(sig, inputIndexBits, writeSig, indexBits, w, u, false)
			}
			for name := range wa.FileOperationsFunction {
				if name == "init" {
					ss.addBootTasks(sig, inputIndexBits, u)
				}
			}
		}
	}

}

func (ss *Server) addWriteAddressTask(wa *WriteAddress, writeSig string, indexBits uint32) {
	for u := range wa.UncoveredAddress {
		ua, ok := ss.corpusDependency.UncoveredAddress[u]
		if !ok {
			return
		}
		for sig, inputIndexBits := range ua.Input {
			ss.addTasks(sig, inputIndexBits, writeSig, indexBits, wa.WriteAddress, u, true)
		}
	}
}

func (ss *Server) addTasks(sig string, indexBits uint32, writeSig string,
	writeIndexBits uint32, writeAddress uint32, uncoveredAddress uint32, high bool) {

	var i uint32
	var index []uint32
	var writeIndex []uint32
	for i = 0; i < 32; i++ {
		if (1<<i)&indexBits > 0 {
			index = append(index, i)
		}
	}
	for i = 0; i < 32; i++ {
		if (1<<i)&writeIndexBits > 0 {
			writeIndex = append(writeIndex, i)
		}
	}

	for _, i := range index {
		for _, wi := range writeIndex {
			if high {
				ss.addTask(ss.getTask(sig, i, writeSig, wi, writeAddress, uncoveredAddress), ss.corpusDependency.HighTask)
			}
			ss.addTask(ss.getTask(sig, i, writeSig, wi, writeAddress, uncoveredAddress), ss.corpusDependency.Tasks)
		}
	}
	return
}

func (ss *Server) addBootTasks(sig string, indexBits uint32, uncoveredAddress uint32) {
	input, ok := ss.corpusDependency.Input[sig]
	if !ok {
		log.Fatalf("addBootTasks do not find sig")
	}
	ua, ok := ss.corpusDependency.UncoveredAddress[uncoveredAddress]
	if !ok {
		log.Fatalf("addBootTasks do not find uncovered address")
	}

	var writeIndexBits uint32
	writeIndexBits = 0
	for i, c := range input.Call {
		for a := range c.Address {
			for wa := range ua.WriteAddress {
				if wa == a {
					writeIndexBits = writeIndexBits | (1 << i)
				}
			}
		}
	}

	var i uint32
	var index []uint32
	for i = 0; i < 32; i++ {
		if (1<<i)&indexBits > 0 {
			index = append(index, i)
		}
	}
	for _, i := range index {
		ss.addTask(ss.getTask(sig, i, sig, writeIndexBits, uncoveredAddress, uncoveredAddress), ss.corpusDependency.BootTask)
	}
	return
}

func (m *Task) getHash() string {
	if m.Hash == "" {
		m.Hash = m.Sig + strconv.FormatInt(int64(m.Index), 10) + m.WriteSig + strconv.FormatInt(int64(m.WriteIndex), 10)
	}
	return m.Hash
}

func (ss *Server) getTask(sig string, index uint32, writeSig string, writeIndex uint32,
	writeAddress uint32, uncoveredAddress uint32) *Task {
	task := &Task{
		Sig:                    sig,
		Index:                  index,
		Program:                []byte{},
		WriteSig:               writeSig,
		WriteIndex:             writeIndex,
		WriteProgram:           []byte{},
		WriteAddress:           writeAddress,
		Priority:               10,
		Hash:                   "",
		UncoveredAddress:       map[uint32]*RunTimeData{},
		CoveredAddress:         map[uint32]*RunTimeData{},
		TaskStatus:             TaskStatus_untested,
		CheckWriteAddress:      false,
		CheckWriteAddressFinal: false,
	}

	task.Hash = task.getHash()

	input, ok := ss.corpusDependency.Input[sig]
	if !ok {
		log.Fatalf("getTask with error sig")
	}
	for _, c := range input.Program {
		task.Program = append(task.Program, c)
	}

	writeInput, ok := ss.corpusDependency.Input[writeSig]
	if !ok {
		log.Fatalf("getTask with error writeSig")
	}
	for _, c := range writeInput.Program {
		task.WriteProgram = append(task.WriteProgram, c)
	}

	ua, ok := ss.corpusDependency.UncoveredAddress[uncoveredAddress]
	if !ok {
		log.Fatalf("getTask with error uncoveredAddress")
	}
	ca := ua.ConditionAddress

	task.UncoveredAddress[uncoveredAddress] = &RunTimeData{
		Program:                 []byte{},
		TaskStatus:              TaskStatus_untested,
		RcursiveCount:           0,
		Priority:                ss.getPriority(task.WriteAddress, uncoveredAddress),
		Idx:                     index,
		CheckCondition:          false,
		ConditionAddress:        ca,
		CheckAddress:            false,
		Address:                 uncoveredAddress,
		CheckRightBranchAddress: false,
		RightBranchAddress:      []uint32{},
	}

	return task
}

func (ss *Server) addTask(task *Task, tasks *Tasks) {
	var uncoveredAddress uint32
	var dr *RunTimeData
	if len(task.UncoveredAddress) == 1 {
		for u, r := range task.UncoveredAddress {
			uncoveredAddress = u
			dr = r
		}
	} else {
		log.Fatalf("addTask more than one uncovered address")
	}

	hash := task.getHash()
	if t, ok := tasks.Task[hash]; ok {
		if _, ok := t.UncoveredAddress[uncoveredAddress]; ok {
			t.UncoveredAddress[uncoveredAddress].updatePriority(dr.Priority)
		} else {
			if t.UncoveredAddress == nil {
				t.UncoveredAddress = map[uint32]*RunTimeData{}
			}
			t.UncoveredAddress[uncoveredAddress] = proto.Clone(dr).(*RunTimeData)
			t.TaskStatus = TaskStatus_untested
		}
		t.updatePriority(task.Priority)
		return
	}
	tasks.addTask(task)
}

func (m *Tasks) addTask(t *Task) {
	if m.Task == nil {
		m.Task = map[string]*Task{}
	}
	if m.Tasks == nil {
		m.Tasks = []*Task{}
	}
	if len(m.Task) != len(m.Tasks) {
		log.Fatalf("%s : len(m.Task) != len(m.Tasks)", m.Name)
	}
	if _, ok := m.Task[t.getHash()]; ok {

	} else {
		m.Task[t.getHash()] = t
		m.Tasks = append(m.Tasks, t)
	}
}

func (m *Tasks) addTasks(t *Tasks) {
	for _, tt := range t.Task {
		m.addTask(tt)
	}
	for _, tt := range t.Tasks {
		m.addTask(tt)
	}
}

func (m *Tasks) emptyTask() {
	m.Task = map[string]*Task{}
	m.Tasks = []*Task{}
}

func (m *Tasks) pop(number int) *Tasks {
	tasks := &Tasks{
		Name:  m.Name,
		Kind:  TaskKind_Normal,
		Task:  map[string]*Task{},
		Tasks: []*Task{},
	}
	temp := []*Task{}
	if number > len(m.Tasks) {
		log.Fatalf("%s : number > len(m.Tasks)", m.Name)
	}
	temp = append(temp, m.Tasks[:number]...)
	m.Tasks = m.Tasks[number:]
	for _, t := range temp {
		delete(m.Task, t.getHash())
		tasks.addTask(t)
	}
	return tasks
}

func (m *Task) updatePriority(p1 int32) {
	m.Priority = m.Priority + p1
	return
}

func (m *Task) getRealPriority() float64 {
	var p uint32
	for _, r := range m.UncoveredAddress {
		p += r.Priority
	}
	res := float64(p) * (math.Pow(2, float64(m.Priority)))
	return res
}

func (m *RunTimeData) updatePriority(p1 uint32) {
	m.Priority += p1
	return
}

func (ss *Server) getPriority(writeAddress uint32, uncoveredAddress uint32) uint32 {
	u, ok := ss.corpusDependency.UncoveredAddress[uncoveredAddress]
	if !ok {
		log.Fatalf("getPriority not find uncoveredAddress")
	}
	bbcount := u.Bbcount
	waa, ok := u.WriteAddress[writeAddress]
	if !ok {
		log.Fatalf("getPriority not find writeAddress")
	}
	pp := waa.Prio
	priority := pp * bbcount
	return priority
}

func (ss *Server) writeCorpusToDisk() {
	out, err := proto.Marshal(ss.corpusDependency)
	if err != nil {
		log.Fatalf("Failed to encode address: %s", err)
	}
	path := "data.bin"
	_ = os.Remove(path)
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		log.Fatalf("Failed to write corpusDependency: %s", err)
	}
}

func (ss *Server) writeStatisticsToDisk() {
	out, err := proto.Marshal(ss.stat)
	if err != nil {
		log.Fatalf("Failed to encode coverage: %s", err)
	}
	path := "statistics.bin"
	_ = os.Remove(path)
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		log.Fatalf("Failed to write coverage: %s", err)
	}
}

// CheckPath ...
func CheckPath(newPath []uint32, unstablePath []uint32) (int, int, int) {

	var l = 0
	var newPathIdx = 0
	var unstablePathIdx = 0
	var idx = 0

	l1 := len(newPath)
	l2 := len(unstablePath)

	if l1 < l2 {
		l = l1
	} else {
		l = l2
	}

	for i := 0; i < l; i++ {
		if newPath[i] == unstablePath[i] {
			newPathIdx = i
			unstablePathIdx = i
			break
		} else {
			for j := 0; j < i; j++ {
				if newPath[i] == unstablePath[j] {
					newPathIdx = i
					unstablePathIdx = j
					break
				} else if unstablePath[i] == newPath[j] {
					unstablePathIdx = i
					newPathIdx = j
					break
				}
			}
		}

		if newPath[newPathIdx] == unstablePath[unstablePathIdx] {
			break
		}
	}
	if newPathIdx == 0 && unstablePathIdx == 0 && newPath[0] != unstablePath[0] {
		log.Logf(0, "newPath : %x\n", newPath)
		log.Logf(0, "unstablePath : %x\n", unstablePath)
		log.Fatalf("checkPath : can not find the address")
	}

	l1 = l1 - newPathIdx
	l2 = l2 - unstablePathIdx

	if l1 < l2 {
		l = l1
	} else {
		l = l2
	}

	for i := 0; i < l; i++ {
		if newPath[i+newPathIdx] != unstablePath[i+unstablePathIdx] {
			idx = i
			break
		}
	}
	return newPathIdx, unstablePathIdx, idx
}
