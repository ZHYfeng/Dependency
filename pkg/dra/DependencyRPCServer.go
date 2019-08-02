package dra

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/syzkaller/pkg/log"
	"github.com/google/syzkaller/pkg/rpctype"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"sync"
)

const (
	taskNum = 4
)

type syzFuzzer struct {
	task Tasks
}

// server is used to implement dra.DependencyServer.
type Server struct {
	address uint32
	Port    int
	Address string

	mu               *sync.Mutex
	inputmu          *sync.Mutex
	coveragemu       *sync.Mutex
	taskmu           *sync.Mutex
	corpusDependency *Corpus

	fuzzermu *sync.Mutex
	fuzzers  map[string]*syzFuzzer

	corpusmu *sync.Mutex
	corpus   *map[string]rpctype.RPCInput

	tmu   *sync.Mutex
	logmu *sync.Mutex
}

func (ss Server) ReturnTasks(ctx context.Context, request *Tasks) (*Empty, error) {
	tasks := CloneTasks(request)

	ss.taskmu.Lock()
	defer ss.taskmu.Unlock()
	for _, task := range tasks.Task {
		for _, t := range ss.corpusDependency.Tasks.Task {
			if t.Sig == task.Sig && t.Index == task.Index &&
				t.WriteSig == task.WriteSig && t.WriteIndex == task.WriteIndex {
				t.MergeTask(task)
			}
		}
	}

	reply := &Empty{}

	return reply, nil
}

func (ss Server) SendDependency(ctx context.Context, request *Dependency) (*Empty, error) {
	log.Logf(1, "(ss Server) SendDependency")
	d := CloneDependency(request)

	ss.mu.Lock()
	defer ss.mu.Unlock()

	for _, wa := range d.WriteAddress {
		ss.addWriteAddress(wa)
	}
	ss.addUncoveredAddress(d.UncoveredAddress)
	ss.addInput(d.Input)
	ss.addInputTask(d.Input)

	reply := &Empty{}

	return reply, nil
}

func (ss Server) GetTasks(context.Context, *Empty) (*Tasks, error) {
	log.Logf(1, "(ss Server) GetTasks")

	ss.taskmu.Lock()
	defer ss.taskmu.Unlock()
	tasks := ss.pickTask()

	return tasks, nil
}

func (ss Server) pickTask() *Tasks {
	log.Logf(1, "(ss Server) GetTasks")

	var i uint
	tasks := &Tasks{
		Name: "",
		Task: []*Task{},
	}
	i = 0
	for _, t := range ss.corpusDependency.Tasks.Task {
		if t.TaskStatus == TaskStatus_untested && len(t.UncoveredAddress) > 0 {
			i++
			tasks.Task = append(tasks.Task, t)
			if i > taskNum {
				break
			}
		}
	}
	return tasks
}

func (ss Server) ReturnDependencyInput(ctx context.Context, request *Dependencytask) (*Empty, error) {
	log.Logf(1, "(ss Server) ReturnDependencyInput")
	//input := CloneInput(request.Input)
	ss.fuzzermu.Lock()
	defer ss.fuzzermu.Unlock()
	//if f, ok := ss.fuzzers[request.Name]; ok {
	//if _, ok := f.task.Task[request.Input.Sig]; ok {
	//	delete(f.corpusDI, request.Input.Sig)
	//	ss.mu.Lock()
	//	defer ss.mu.Unlock()
	//	if ok := ss.checkDependencyInput(input); ok {
	//		//ss.corpusDependency.Input[input.Sig] = input
	//	} else {
	//		//ss.corpusDependency.Input[input.Sig] = input
	//	}
	//}
	//} else {
	//	log.Fatalf("ReturnDependencyInput : ", request.Name)
	//}
	reply := &Empty{}
	ss.writeToDisk()

	return reply, nil
}

func (ss Server) GetCondition(context.Context, *Empty) (*Conditions, error) {
	log.Logf(1, "(ss Server) GetCondition")
	reply := &Conditions{
		//Condition: map[uint64]*Condition{},
		Condition: []*Condition{},
	}
	//for _, wa := range ss.corpusDependency.WriteAddress {
	//if len(wa.) == 0 {
	//	//reply.Condition[wa.Condition.ConditionAddress] = CloneCondition(wa.Condition)
	//	reply.Condition = append(reply.Condition, CloneCondition(wa.Condition))
	//	return reply, nil
	//}
	//}
	//ss.writeToDisk()

	return reply, nil
}

func (ss Server) SendWriteAddress(ctx context.Context, request *WriteAddresses) (*Empty, error) {
	log.Logf(1, "(ss Server) SendWriteAddress")
	ss.mu.Lock()
	defer ss.mu.Unlock()
	//a := request.Condition.ConditionAddress<<32 + request.Condition.Successor
	//if wa, ok := ss.corpusDependency.WriteAddress[a]; ok {
	//	for _, wwa := range request.WriteAddress {
	//		//wa.WriteAddress[wwa.WriteAddress] = CloneWriteAddress(wwa)
	//		wa.WriteAddress = append(wa.WriteAddress, CloneWriteAddress(wwa))
	//	}
	//	for sig, i := range ss.corpusDependency.CorpusRecursiveInput {
	//		if ok := ss.checkDependencyInput(i); ok {
	//			ss.corpusDependency.CorpusRecursiveInput[sig] = CloneInput(i)
	//		} else {
	//		}
	//	}
	//} else {
	//	log.Fatalf("SendWriteAddress : ", request.Condition.ConditionAddress)
	//}
	ss.writeToDisk()

	return &Empty{}, nil
}

func (ss Server) SendLog(ctx context.Context, request *Empty) (*Empty, error) {
	log.Logf(1, "(ss Server) SendLog")
	ss.logmu.Lock()
	defer ss.logmu.Unlock()

	f, _ := os.OpenFile("./dependency.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	_, _ = f.WriteString(fmt.Sprintf(request.Name))

	reply := &Empty{}
	return reply, nil
}

func (ss Server) Connect(ctx context.Context, request *Empty) (*Empty, error) {
	log.Logf(1, "(ss Server) Connect")
	ss.fuzzermu.Lock()
	defer ss.fuzzermu.Unlock()

	if _, ok := ss.fuzzers[request.Name]; !ok {
		ss.fuzzers[request.Name] = &syzFuzzer{
			task: Tasks{Task: []*Task{}},
		}
	}
	return &Empty{}, nil
}

func (ss Server) GetVmOffsets(context.Context, *Empty) (*Empty, error) {
	reply := &Empty{}
	reply.Address = ss.address
	return reply, nil
}

func (ss Server) GetNewInput(context.Context, *Empty) (*Inputs, error) {
	log.Logf(1, "(ss Server) GetNewInput")

	reply := &Inputs{
		//Input: map[string]*Input{},
		Input: []*Input{},
	}

	i := 0
	ss.inputmu.Lock()
	defer ss.inputmu.Unlock()
	for s, c := range ss.corpusDependency.NewInput {
		if i < 1 {
			//reply.Input[c.Sig] = CloneInput(c)
			reply.Input = append(reply.Input, CloneInput(c))
			i++
			ss.mu.Lock()
			defer ss.mu.Unlock()
			ss.addInput(c)
			delete(ss.corpusDependency.NewInput, s)
		} else {
		}
	}
	return reply, nil
}

func (ss Server) SendDependencyInput(ctx context.Context, request *Input) (*Empty, error) {
	log.Logf(1, "(ss Server) SendDependencyInput")
	reply := &Empty{}

	if len(request.Program) == 0 {
		reply.Name = "dependency Program error : " + request.Sig
		return reply, nil
	} else if len(request.Sig) == 0 {
		reply.Name = "dependency Sig error : " + string(request.Program)
		return reply, nil
	}

	ss.mu.Lock()
	defer ss.mu.Unlock()

	ci := CloneInput(request)
	ss.addInput(ci)

	reply.Name = "success"
	return reply, nil
}

//
func (ss Server) GetDependencyInput(ctx context.Context, request *Empty) (*Inputs, error) {
	log.Logf(1, "(ss Server) GetDependencyInput")
	reply := &Inputs{
		//Input: map[string]*Input{},
		Input: []*Input{},
	}
	//if f, ok := ss.fuzzers[request.Name]; ok {
	//
	//	ss.fuzzermu.Lock()
	//	defer ss.fuzzermu.Unlock()
	//if len(f.corpusDI) > 0 {
	//for s, c := range f.corpusDI {
	//ss.corpusDependency.CorpusErrorInput[s] = c
	//delete(f.corpusDI, s)
	//}
	//}
	//
	//	ss.mu.Lock()
	//	defer ss.mu.Unlock()
	//i := 0
	//for s, c := range ss.corpusDependency.CorpusDependencyInput {
	//	if i < taskNum {
	//		i++
	//
	//		reply.Input[c.Sig] = CloneInput(c)
	//reply.Input = append(reply.Input, CloneInput(c))
	//f.corpusDI[s] = c
	//delete(ss.corpusDependency.CorpusDependencyInput, s)
	//return reply, nil
	//} else {
	//}
	//}
	//} else {
	//	log.Fatalf("syz_fuzzer %v is not connected", request.Name)
	//}

	//ss.writeToDisk()

	//for i := 0; i < 50 && len(f.corpusDependencyInput) > 0; i++ {
	//	last := len(f.corpusDependencyInput) - 1
	//	reply.DependencyInput = append(reply.DependencyInput, cloneDependencyInput(&f.corpusDependencyInput[last]))
	//	f.corpusDependencyInput = f.corpusDependencyInput[:last]
	//}
	//if len(f.corpusDependencyInput) == 0 {
	//	f.corpusDependencyInput = nil
	//}
	return reply, nil
}

func (ss Server) SendNewInput(ctx context.Context, request *Input) (*Empty, error) {
	log.Logf(1, "(ss Server) SendNewInput")

	reply := &Empty{}
	input := CloneInput(request)

	ss.inputmu.Lock()
	defer ss.inputmu.Unlock()
	ss.addNewInput(input)

	ss.coveragemu.Lock()
	defer ss.coveragemu.Unlock()
	ss.addCoveredAddress(input)

	return reply, nil
}

//
//func (ss Server) checkDependencyInput(request *Input) (res bool) {
//	res = false
//	for _, u := range request.UncoveredAddress {
//		if u.RunTimeDate.TaskStatus == RunTimeData_recursive {
//			for _, wa := range u.WriteAddress {
//				res = res || ss.checkWriteAddress(wa)
//			}
//		}
//	}
//	return res
//}
//
//func (ss Server) checkWriteAddress(wa *WriteAddress) (res bool) {
//	res = false
//	if wa.RunTimeDate.TaskStatus == RunTimeData_recursive {
//		for _, wc := range wa.WriteSyscall {
//			if wc.RunTimeDate.TaskStatus == RunTimeData_recursive {
//				if len(wc.WriteAddress) != 0 {
//					for _, wwa := range wc.WriteAddress {
//						res = res || ss.checkWriteAddress(wwa)
//					}
//				} else {
//					res = res || ss.checkCondition(wc)
//				}
//			}
//		}
//	}
//	return res
//}
//
//func (ss Server) checkCondition(wc *Syscall) (res bool) {
//	res = false
//	condition := wc.RunTimeDate.ConditionAddress
//	if cc, ok := wc.CriticalCondition[condition]; ok {
//		a := cc.ConditionAddress<<32 + cc.Successor
//		if wa, ok := ss.corpusDependency.WriteAddress[a]; ok {
//			if len(wa.WriteAddress) > 0 {
//				res = true
//				for _, wwa := range wa.WriteAddress {
//					temp := CloneWriteAddress(wwa)
//					temp.RunTimeDate = CloneRunTimeData(wc.RunTimeDate)
//
//					//wc.WriteAddress[temp.WriteAddress] = temp
//					wc.WriteAddress = append(wc.WriteAddress, temp)
//
//					for _, wwc := range temp.WriteSyscall {
//						wwc.RunTimeDate = CloneRunTimeData(wc.RunTimeDate)
//						wwc.RunTimeDate.Address = wwa.WriteAddress
//					}
//				}
//			}
//		} else {
//			ss.corpusDependency.WriteAddress[a] = &WriteAddresses{
//				Condition: CloneCondition(cc),
//				//WriteAddress: map[uint32]*WriteAddress{},
//				WriteAddress: []*WriteAddress{},
//			}
//		}
//	}
//	return res
//}

func (ss *Server) addNewInput(s *Input) {
	if i, ok := ss.corpusDependency.NewInput[s.Sig]; ok {
		i.MergeInput(s)
	} else {
		ss.corpusDependency.NewInput[s.Sig] = s
	}

	return
}

func (ss *Server) addInput(s *Input) {
	if i, ok := ss.corpusDependency.Input[s.Sig]; ok {
		i.MergeInput(s)
	} else {
		ss.corpusDependency.Input[s.Sig] = s
	}

	ss.addWriteAddressMapInput(s)
	ss.addUncoveredAddressMapInput(s)

	return
}

func (ss *Server) addWriteAddressMapInput(s *Input) {
	sig := s.Sig
	input := ss.corpusDependency.Input[sig]
	for index, call := range s.Call {
		indexBits := uint32(1 << index)
		for a := range call.Address {
			if wa, ok := ss.corpusDependency.WriteAddress[a]; ok {
				waIndex, ok := wa.Input[sig]
				if ok {
					//
					if (waIndex|indexBits)^waIndex > 0 {
						ss.addWriteAddressTask(wa, sig, (waIndex|indexBits)^waIndex)
						wa.Input[sig] = waIndex | indexBits
					}
				} else {
					ss.addWriteAddressTask(wa, sig, indexBits)
					wa.Input[sig] = indexBits
				}

				if iIndex, ok := input.WriteAddress[a]; ok {
					input.WriteAddress[a] = iIndex | indexBits
				} else {
					input.WriteAddress[a] = indexBits
				}

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
				u2.Input[sig] = i1
			}
		}
	}
	return
}

func (ss *Server) checkUncoveredAddress(uncoveredAddress uint32) bool {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	_, ok := ss.corpusDependency.UncoveredAddress[uncoveredAddress]
	if !ok {
		return false
	} else {
		ss.deleteUncoveredAddress(uncoveredAddress)
	}
	return true
}

func (ss *Server) deleteUncoveredAddress(uncoveredAddress uint32) {
	u, ok := ss.corpusDependency.UncoveredAddress[uncoveredAddress]
	if !ok {
		return
	}

	for sig, _ := range u.Input {
		input, ok := ss.corpusDependency.Input[sig]
		if !ok {
			log.Fatalf("deleteUncoveredAddress not find sig")
			continue
		} else {
			delete(u.Input, sig)
		}
		_, ok1 := input.UncoveredAddress[uncoveredAddress]
		if !ok1 {
			log.Fatalf("deleteUncoveredAddress input not find uncoveredAddress")
		} else {
			delete(input.UncoveredAddress, uncoveredAddress)

		}
	}

	for wa, _ := range u.WriteAddress {
		waa, ok := ss.corpusDependency.WriteAddress[wa]
		if !ok {
			log.Fatalf("deleteUncoveredAddress not find wa")
			continue
		} else {
			delete(u.WriteAddress, wa)
		}
		_, ok1 := waa.UncoveredAddress[uncoveredAddress]
		if !ok1 {
			log.Fatalf("deleteUncoveredAddress write address not find uncoveredAddress")
		} else {
			delete(waa.UncoveredAddress, uncoveredAddress)
		}
	}

	delete(ss.corpusDependency.UncoveredAddress, uncoveredAddress)

	return
}

func (ss *Server) addCoveredAddress(input *Input) {
	var isDependency uint32
	if input.Dependency {
		isDependency = 1
	} else {
		isDependency = 0
	}

	for _, call := range input.Call {
		for a := range call.Address {
			ss.corpusDependency.Coverage.Coverage[a] = isDependency
			ss.checkUncoveredAddress(a)
		}
	}

	out, err := proto.Marshal(ss.corpusDependency.Coverage)
	if err != nil {
		log.Fatalf("Failed to encode coverage:", err)
	}
	path := "coverage.bin"
	_ = os.Remove(path)
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		log.Fatalf("Failed to write coverage:", err)
	}
	return
}

func CloneDependency(s *Dependency) *Dependency {
	d := &Dependency{
		Input:            CloneInput(s.Input),
		UncoveredAddress: CloneUncoverAddress(s.UncoveredAddress),
		WriteAddress:     []*WriteAddress{},
	}

	for _, wa := range s.WriteAddress {
		d.WriteAddress = append(d.WriteAddress, CloneWriteAddress(wa))
	}

	return d
}

func CloneInput(s *Input) *Input {
	d := &Input{
		Sig:              s.Sig,
		Program:          []byte{},
		Call:             make(map[uint32]*Call),
		Dependency:       s.Dependency,
		UncoveredAddress: map[uint32]uint32{},
		WriteAddress:     map[uint32]uint32{},
	}

	for _, c := range s.Program {
		d.Program = append(d.Program, c)
	}

	for i, u := range s.Call {
		u1 := &Call{
			//Address: []uint32{},
			Address: make(map[uint32]uint32),
			Idx:     u.Idx,
		}
		for aa := range u.Address {
			u1.Address[aa] = 0
		}
		d.Call[i] = u1
	}

	for i, c := range s.UncoveredAddress {
		d.UncoveredAddress[i] = c
	}

	for i, c := range s.WriteAddress {
		d.WriteAddress[i] = c
	}

	return d
}

func (m *Input) MergeInput(d *Input) {

	for i, u := range d.Call {
		var call *Call
		if c, ok := m.Call[i]; ok {
			call = c
		} else {
			call = &Call{
				//Address: []uint32{},
				Address: make(map[uint32]uint32),
				Idx:     u.Idx,
			}
			d.Call[i] = call
		}

		for a := range u.Address {
			call.Address[a] = 0
		}
	}

	for i, c := range d.UncoveredAddress {
		if index, ok := m.UncoveredAddress[i]; ok {
			m.UncoveredAddress[i] = index | c
		} else {
			m.UncoveredAddress[i] = c
		}
	}

	for i, c := range d.WriteAddress {
		if index, ok := m.WriteAddress[i]; ok {
			m.WriteAddress[i] = index | c
		} else {
			m.WriteAddress[i] = c
		}
	}

	return
}

func (ss *Server) addUncoveredAddress(s *UncoveredAddress) {
	if i, ok := ss.corpusDependency.UncoveredAddress[s.UncoveredAddress]; ok {
		i.MergeUncoveredAddress(s)
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

func CloneUncoverAddress(s *UncoveredAddress) *UncoveredAddress {
	d := &UncoveredAddress{
		ConditionAddress:   s.ConditionAddress,
		UncoveredAddress:   s.UncoveredAddress,
		RightBranchAddress: []uint32{},
		Input:              map[string]uint32{},
		WriteAddress:       map[uint32]*WriteAddressAttributes{},
		RunTimeDate:        CloneRunTimeData(s.RunTimeDate),
	}

	for i, c := range s.Input {
		d.Input[i] = c
	}

	for i, c := range s.WriteAddress {
		d.WriteAddress[i] = CloneWriteAddressAttributes(c)
	}

	return d
}

func (m *UncoveredAddress) MergeUncoveredAddress(d *UncoveredAddress) {

	for i, c := range d.Input {
		if index, ok := m.Input[i]; ok {
			m.Input[i] = index | c
		} else {
			m.Input[i] = c
		}
	}

	for i, c := range d.WriteAddress {
		if _, ok := m.WriteAddress[i]; ok {

		} else {
			m.WriteAddress[i] = CloneWriteAddressAttributes(c)
		}
	}

	return
}

func CloneWriteAddressAttributes(s *WriteAddressAttributes) *WriteAddressAttributes {
	d := &WriteAddressAttributes{
		WriteAddress: s.WriteAddress,
		Repeat:       s.Repeat,
		Prio:         s.Prio,
	}
	return d
}

func (ss *Server) addWriteAddress(s *WriteAddress) {
	if i, ok := ss.corpusDependency.WriteAddress[s.WriteAddress]; ok {
		if len(s.Input) != len(i.Input) {
			log.Fatalf("addWriteAddress len(s.Input) != len(i.Input)")
		}
		i.MergeWriteAddress(s)
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
				waInput.WriteAddress[s.WriteAddress] = indexBits1
			}
		} else {
			log.Fatalf("addWriteAddress not find sig")
		}
	}
}

func CloneWriteAddress(s *WriteAddress) *WriteAddress {
	d := &WriteAddress{
		WriteAddress:     s.WriteAddress,
		ConditionAddress: s.ConditionAddress,
		UncoveredAddress: map[uint32]*WriteAddressAttributes{},
		IoctlCmd:         map[uint64]uint32{},
		Input:            map[string]uint32{},

		RunTimeDate: CloneRunTimeData(s.RunTimeDate),
	}

	for i, c := range s.UncoveredAddress {
		d.UncoveredAddress[i] = CloneWriteAddressAttributes(c)
	}

	for i, c := range s.IoctlCmd {
		d.IoctlCmd[i] = c
	}

	for i, c := range s.Input {
		d.Input[i] = c
	}
	return d
}

func (m *WriteAddress) MergeWriteAddress(d *WriteAddress) {

	for i, c := range d.UncoveredAddress {
		if _, ok := m.UncoveredAddress[i]; ok {

		} else {
			m.UncoveredAddress[i] = CloneWriteAddressAttributes(c)
		}
	}

	for i, c := range d.IoctlCmd {
		if ii, ok := m.IoctlCmd[i]; ok {
			m.IoctlCmd[i] = ii | c
		} else {
			m.IoctlCmd[i] = c
		}
	}

	for i, c := range d.Input {
		if index, ok := m.Input[i]; ok {
			m.Input[i] = index | c
		} else {
			m.Input[i] = c
		}
	}

	return
}

func CloneIoctlCmdInput(s *IoctlCmdInput) *IoctlCmdInput {
	d := &IoctlCmdInput{
		Sig:          s.Sig,
		Index:        s.Index,
		Cmd:          s.Cmd,
		WriteAddress: s.WriteAddress,
	}
	return d
}

func CloneIoctlCmd(s *IoctlCmd) *IoctlCmd {
	d := &IoctlCmd{
		Name: s.Name,
		Cmd:  s.Cmd,
		//CriticalCondition: map[uint32]*Condition{},
		RunTimeDate: CloneRunTimeData(s.RunTimeDate),

		WriteAddress: map[uint32]uint32{},
	}

	//for i, c := range s.CriticalCondition {
	//	d.CriticalCondition[i] = CloneCondition(c)
	//}

	for i, c := range s.WriteAddress {
		d.WriteAddress[i] = c
	}

	return d
}

func CloneCondition(c *Condition) *Condition {
	c1 := &Condition{
		ConditionAddress:          c.ConditionAddress,
		SyzkallerConditionAddress: c.SyzkallerConditionAddress,
		UncoveredAddress:          c.UncoveredAddress,
		SyzkallerUncoveredAddress: c.SyzkallerUncoveredAddress,
		Idx:                       c.Idx,
		Successor:                 c.Successor,
		//RightBranchAddress:          map[uint64]uint32{},
		//SyzkallerRightBranchAddress: map[uint32]uint32{},
		//WrongBranchAddress:          map[uint64]uint32{},
		//SyzkallerWrongBranchAddress: map[uint32]uint32{},
		RightBranchAddress:          []uint64{},
		SyzkallerRightBranchAddress: []uint32{},
		//WrongBranchAddress:          []uint64{},
		//SyzkallerWrongBranchAddress: []uint32{},
	}

	for _, a := range c.RightBranchAddress {
		//c1.RightBranchAddress[a] = 0
		c1.RightBranchAddress = append(c1.RightBranchAddress, a)
	}

	for _, a := range c.SyzkallerRightBranchAddress {
		//c1.SyzkallerRightBranchAddress[a] = 0
		c1.SyzkallerRightBranchAddress = append(c1.SyzkallerRightBranchAddress, a)
	}

	//for _, a := range c.WrongBranchAddress {
	//	//c1.WrongBranchAddress[a] = 0
	//	c1.WrongBranchAddress = append(c1.WrongBranchAddress, a)
	//}

	//for _, a := range c.SyzkallerWrongBranchAddress {
	//	//c1.SyzkallerWrongBranchAddress[a] = 0
	//	c1.SyzkallerWrongBranchAddress = append(c1.SyzkallerWrongBranchAddress, a)
	//}

	return c1
}

func CloneRunTimeData(d *RunTimeData) *RunTimeData {
	d1 := &RunTimeData{
		Program:                 []byte{},
		TaskStatus:              d.TaskStatus,
		RcursiveCount:           d.RcursiveCount,
		Priority:                d.Priority,
		Idx:                     d.Idx,
		CheckCondition:          d.CheckCondition,
		ConditionAddress:        d.ConditionAddress,
		CheckAddress:            d.CheckAddress,
		Address:                 d.Address,
		CheckRightBranchAddress: d.CheckRightBranchAddress,
		//RightBranchAddress:      map[uint32]uint32{},
		RightBranchAddress: []uint32{},
	}

	for _, c := range d.Program {
		d1.Program = append(d1.Program, c)
	}

	for _, a := range d.RightBranchAddress {
		d1.RightBranchAddress = append(d1.RightBranchAddress, a)
		//d1.RightBranchAddress[a] = 0
	}

	return d1
}

func (m *RunTimeData) MergeRunTimeData(d *RunTimeData) {
	if d == nil {
		return
	}

	return
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
				ss.addTasks(sig, inputIndexBits, writeSig, indexBits, w, u)
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
			ss.addTasks(sig, inputIndexBits, writeSig, indexBits, wa.WriteAddress, u)
		}
	}
}

func (ss *Server) addTasks(sig string, indexBits uint32, writeSig string,
	writeIndexBits uint32, writeAddress uint32, uncoveredAddress uint32) {
	var i uint32
	var index []uint32
	var writeIndex []uint32
	for i = 0; i < 32; i++ {
		if (1<<i)&indexBits == 1 {
			index = append(index, i)
		}
	}
	for i = 0; i < 32; i++ {
		if (1<<i)&writeIndexBits == 1 {
			writeIndex = append(writeIndex, i)
		}
	}

	for _, i := range index {
		for _, wi := range writeIndex {
			ss.addTask(ss.getTask(sig, i, writeSig, wi, writeAddress, uncoveredAddress))
		}
	}
	return
}

func (ss *Server) getTask(sig string, index uint32, writeSig string, writeIndex uint32,
	writeAddress uint32, uncoveredAddress uint32) *Task {
	task := &Task{
		Sig:              sig,
		Index:            index,
		Program:          []byte{},
		WriteSig:         writeSig,
		WriteIndex:       writeIndex,
		WriteProgram:     []byte{},
		WriteAddress:     writeAddress,
		UncoveredAddress: map[uint32]*RunTimeData{},
		TaskStatus:       TaskStatus_untested,
	}

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

func (ss *Server) addTask(task *Task) {
	ss.taskmu.Lock()
	defer ss.taskmu.Unlock()

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

	for _, t := range ss.corpusDependency.Tasks.Task {
		if t.Sig == task.Sig && t.Index == task.Index &&
			t.WriteSig == task.WriteSig && t.WriteIndex == task.WriteIndex {
			if r, ok := t.UncoveredAddress[uncoveredAddress]; ok {
				t.UncoveredAddress[uncoveredAddress].Priority = ss.updatePriority(r.Priority, r.Priority)
			} else {
				t.UncoveredAddress[uncoveredAddress] = CloneRunTimeData(dr)
			}
		}
	}
}

func (ss *Server) updatePriority(p1 uint32, p2 uint32) uint32 {
	priority := p1 + p2
	return priority
}

func (ss *Server) getPriority(writeAddress uint32, uncoveredAddress uint32) uint32 {
	priority := writeAddress + uncoveredAddress
	return priority
}

func (ss *Server) SetAddress(address uint32) {
	ss.address = address
}

func CloneTasks(s *Tasks) *Tasks {
	d := &Tasks{
		Name: s.Name,
		Task: []*Task{},
	}
	for _, t := range s.Task {
		d.Task = append(d.Task, CloneTask(t))
	}
	return d
}

func CloneTask(s *Task) *Task {
	d := &Task{
		Sig:              s.Sig,
		Index:            s.Index,
		Program:          []byte{},
		WriteSig:         s.WriteSig,
		WriteIndex:       s.WriteIndex,
		WriteProgram:     []byte{},
		WriteAddress:     s.WriteAddress,
		UncoveredAddress: map[uint32]*RunTimeData{},
		TaskStatus:       s.TaskStatus,
		CoveredAddress:   map[uint32]*RunTimeData{},
	}

	for _, c := range s.Program {
		d.Program = append(d.Program, c)
	}

	for _, c := range s.WriteProgram {
		d.WriteProgram = append(d.WriteProgram, c)
	}

	for u, p := range s.UncoveredAddress {
		d.UncoveredAddress[u] = CloneRunTimeData(p)
	}

	for u, p := range s.CoveredAddress {
		d.CoveredAddress[u] = CloneRunTimeData(p)
	}

	return d
}

func (m *Task) MergeTask(s *Task) {
	for u, p := range s.CoveredAddress {
		m.CoveredAddress[u] = CloneRunTimeData(p)
	}

	for u := range m.UncoveredAddress {
		_, ok := m.CoveredAddress[u]
		if ok {
			delete(m.UncoveredAddress, u)
		}
	}

	return
}

// RunDependencyRPCServer
func (ss *Server) RunDependencyRPCServer(corpus *map[string]rpctype.RPCInput) {

	ss.corpusDependency = &Corpus{
		Input:            map[string]*Input{},
		UncoveredAddress: map[uint32]*UncoveredAddress{},
		CoveredAddress:   map[uint32]*UncoveredAddress{},
		WriteAddress:     map[uint32]*WriteAddress{},
		IoctlCmd:         map[uint64]*IoctlCmd{},
		Tasks:            &Tasks{Task: []*Task{}},
		Coverage:         &Coverage{Coverage: map[uint32]uint32{}},
		NewInput:         map[string]*Input{},
	}
	ss.fuzzers = make(map[string]*syzFuzzer)

	ss.mu = &sync.Mutex{}
	ss.inputmu = &sync.Mutex{}
	ss.coveragemu = &sync.Mutex{}
	ss.taskmu = &sync.Mutex{}
	ss.fuzzermu = &sync.Mutex{}
	ss.corpusmu = &sync.Mutex{}
	ss.tmu = &sync.Mutex{}
	ss.logmu = &sync.Mutex{}

	ss.corpus = corpus

	lis, err := net.Listen("tcp", ss.Address)
	log.Logf(0, "drpc on tcp : %s", ss.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ss.Port = lis.Addr().(*net.TCPAddr).Port
	s := grpc.NewServer(grpc.MaxRecvMsgSize(0x7fffffffffffffff), grpc.MaxSendMsgSize(0x7fffffffffffffff))
	RegisterDependencyRPCServer(s, ss)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (ss *Server) writeToDisk() {
	// Write the new back to disk.
	out, err := proto.Marshal(ss.corpusDependency)
	if err != nil {
		log.Fatalf("Failed to encode address:", err)
	}
	ss.tmu.Lock()
	defer ss.tmu.Unlock()
	path := "data.bin"
	_ = os.Remove(path)
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		log.Fatalf("Failed to write corpusDependency:", err)
	}
	// [END marshal_proto]
}
