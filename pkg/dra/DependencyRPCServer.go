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
	port    = ":50051"
	taskNum = 1
)

type fuzzer struct {
	corpusDI map[string]*Input
}

// server is used to implement dra.DependencyServer.
type Server struct {
	address uint32
	Dport   int
	//corpusDC []*Input

	imu      *sync.Mutex
	corpusDC map[string]*Input

	mu               *sync.Mutex
	corpusDependency *Corpus

	fmu     *sync.Mutex
	fuzzers map[string]*fuzzer

	cmu    *sync.Mutex
	corpus *map[string]rpctype.RPCInput
}

func (ss Server) ReturnDependencyInput(ctx context.Context, request *Task) (*Empty, error) {
	ss.fmu.Lock()
	defer ss.fmu.Unlock()
	if f, ok := ss.fuzzers[request.Name]; ok {
		if _, ok := f.corpusDI[request.Input.Sig]; ok {
			delete(f.corpusDI, request.Input.Sig)
			ss.mu.Lock()
			defer ss.mu.Unlock()
			if ok := ss.checkDependencyInput(request.Input); ok {
				ss.corpusDependency.CorpusDependencyInput[request.Input.Sig] = CloneInput(request.Input)
			} else {
				ss.corpusDependency.CorpusRecursiveInput[request.Input.Sig] = CloneInput(request.Input)
			}
		}
	} else {
		log.Fatalf("ReturnDependencyInput : ", request.Name)
	}
	reply := &Empty{}

	ss.writeToDisk()

	return reply, nil
}

func (ss Server) GetCondition(context.Context, *Empty) (*Conditions, error) {
	reply := &Conditions{}
	for _, wa := range ss.corpusDependency.WriteAddress {
		if len(wa.WriteAddress) == 0 {
			reply.Condition = append(reply.Condition, CloneCondition(wa.Condition))
			return reply, nil
		}
	}
	ss.writeToDisk()

	return reply, nil
}

func (ss Server) SendWriteAddress(ctx context.Context, request *WriteAddresses) (*Empty, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	a := request.Condition.ConditionAddress<<32 + request.Condition.Successor
	if wa, ok := ss.corpusDependency.WriteAddress[a]; ok {
		for _, wwa := range request.WriteAddress {
			wa.WriteAddress = append(wa.WriteAddress, CloneWriteAddress(wwa))
		}
		for sig, i := range ss.corpusDependency.CorpusRecursiveInput {
			if ok := ss.checkDependencyInput(i); ok {
				delete(ss.corpusDependency.CorpusRecursiveInput, sig)
				ss.corpusDependency.CorpusRecursiveInput[sig] = CloneInput(i)
			} else {
			}
		}
	} else {
		log.Fatalf("SendWriteAddress : ", request.Condition.ConditionAddress)
	}

	ss.writeToDisk()

	return &Empty{}, nil
}

func (ss Server) SendLog(ctx context.Context, request *Empty) (*Empty, error) {
	ss.fmu.Lock()
	defer ss.fmu.Unlock()

	f, _ := os.OpenFile("./dependency.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	_, _ = f.WriteString(fmt.Sprintf(request.Name))

	reply := &Empty{}
	return reply, nil
}

func (ss Server) Connect(ctx context.Context, request *Empty) (*Empty, error) {
	if _, ok := ss.fuzzers[request.Name]; !ok {
		ss.fuzzers[request.Name] = &fuzzer{
			corpusDI: map[string]*Input{},
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
	reply := &Inputs{}
	i := 0
	for s, c := range ss.corpusDC {
		if i < 1 {
			reply.Input = append(reply.Input, CloneInput(c))
			i++
			delete(ss.corpusDC, s)
		} else {
		}

	}
	return reply, nil
}

func (ss Server) SendDependencyInput(ctx context.Context, request *Input) (*Empty, error) {
	reply := &Empty{}
	cd := CloneInput(request)

	if len(cd.Program) == 0 {
		reply.Name = "dependency Program error : " + cd.Sig
		return reply, nil
	} else if len(cd.Sig) == 0 {
		reply.Name = "dependency Sig error : " + string(cd.Program)
		return reply, nil
	}

	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.corpusDependency.CorpusDependencyInput[request.Sig] = cd

	ss.writeToDisk()

	reply.Name = "success"
	return reply, nil
}

//
func (ss Server) GetDependencyInput(ctx context.Context, request *Empty) (*Inputs, error) {
	reply := &Inputs{}
	if f, ok := ss.fuzzers[request.Name]; ok {

		ss.mu.Lock()
		defer ss.mu.Unlock()

		if len(f.corpusDI) > 0 {
			for s, c := range f.corpusDI {
				ss.corpusDependency.CorpusErrorInput[s] = c
				delete(f.corpusDI, s)
			}
		}

		i := 0
		for s, c := range ss.corpusDependency.CorpusDependencyInput {
			if i < taskNum {
				i++

				reply.Input = append(reply.Input, CloneInput(c))
				f.corpusDI[s] = c
				delete(ss.corpusDependency.CorpusDependencyInput, s)
				return reply, nil
			} else {
			}
		}
	} else {
		log.Fatalf("fuzzer %v is not connected", request.Name)
	}

	ss.writeToDisk()

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
	ss.mu.Lock()
	defer ss.mu.Unlock()
	reply := &Empty{}
	input := CloneInput(request)

	//ss.corpusDC = append(ss.corpusDC, input)
	ss.corpusDC[input.Sig] = input

	return reply, nil
}

func (ss Server) checkDependencyInput(request *Input) (res bool) {
	res = false
	for _, u := range request.UncoveredAddress {
		if u.RunTimeDate.TaskStatus == RunTimeData_recursive {
			for _, wa := range u.WriteAddress {
				res = res || ss.checkWriteAddress(wa)
			}
		}
	}
	return res
}

func (ss Server) checkWriteAddress(wa *WriteAddress) (res bool) {
	res = false
	if wa.RunTimeDate.TaskStatus == RunTimeData_recursive {
		for _, wc := range wa.WriteSyscall {
			if wc.RunTimeDate.TaskStatus == RunTimeData_recursive {
				if len(wc.WriteAddress) != 0 {
					for _, wwa := range wc.WriteAddress {
						res = res || ss.checkWriteAddress(wwa)
					}
				} else {
					res = res || ss.checkCondition(wc)
				}
			}
		}
	}
	return res
}

func (ss Server) checkCondition(wc *Syscall) (res bool) {
	res = false
	condition := wc.RunTimeDate.ConditionAddress
	if cc, ok := wc.CriticalCondition[condition]; ok {
		a := cc.ConditionAddress<<32 + cc.Successor
		if wa, ok := ss.corpusDependency.WriteAddress[a]; ok {
			if len(wa.WriteAddress) > 0 {
				res = true
				for _, wwa := range wa.WriteAddress {
					temp := CloneWriteAddress(wwa)
					wc.WriteAddress = append(wc.WriteAddress, temp)

					temp.RunTimeDate = CloneRunTimeData(wc.RunTimeDate)
					for _, wwc := range temp.WriteSyscall {
						wwc.RunTimeDate = CloneRunTimeData(temp.RunTimeDate)
						wwc.RunTimeDate.Address = wwa.WriteAddress
					}
				}
			}
		} else {
			ss.corpusDependency.WriteAddress[a] = new(WriteAddresses)
			ss.corpusDependency.WriteAddress[a].Condition = CloneCondition(cc)
		}
	}
	return res
}

// not finish
func (m *Input) Merge(i *Input) {

}

func CloneInput(input *Input) *Input {
	if input == nil {
		return &Input{}
	}
	inputClone := &Input{
		Sig:              input.Sig,
		Program:          []byte{},
		Call:             make(map[uint32]*Call),
		Dependency:       input.Dependency,
		UncoveredAddress: []*UncoveredAddress{},
		WriteAddress:     input.WriteAddress,
		Idx:              input.Idx,
	}

	for _, c := range input.Program {
		inputClone.Program = append(inputClone.Program, c)
	}

	for i, u := range input.Call {
		u1 := &Call{
			Address: make(map[uint32]uint32),
			Idx:     u.Idx,
		}
		for aa := range u.Address {
			u1.Address[aa] = 0
		}
		inputClone.Call[i] = u1
	}

	for _, u := range input.UncoveredAddress {
		u1 := &UncoveredAddress{
			ConditionAddress: u.ConditionAddress,
			UncoveredAddress: u.UncoveredAddress,
			RunTimeDate:      CloneRunTimeData(u.RunTimeDate),
			WriteAddress:     []*WriteAddress{},
		}
		for _, wa := range u.WriteAddress {
			u1.WriteAddress = append(u1.WriteAddress, CloneWriteAddress(wa))
		}
		inputClone.UncoveredAddress = append(inputClone.UncoveredAddress, u1)
	}

	return inputClone
}

func CloneWriteAddress(a *WriteAddress) *WriteAddress {
	if a == nil {
		return &WriteAddress{}
	}
	a1 := &WriteAddress{
		Repeat:     a.Repeat,
		RealRepeat: a.RealRepeat,
		Prio:       a.Prio,

		WriteAddress: a.WriteAddress,

		ConditionAddress: a.ConditionAddress,
		WriteInput:       []*Input{},
		WriteSyscall:     []*Syscall{},

		RunTimeDate: CloneRunTimeData(a.RunTimeDate),
	}

	for _, i := range a.WriteInput {
		a1.WriteInput = append(a1.WriteInput, CloneInput(i))
	}

	for _, s := range a.WriteSyscall {
		a1.WriteSyscall = append(a1.WriteSyscall, CloneSyscall(s))
	}
	return a1
}

func CloneSyscall(s *Syscall) *Syscall {
	if s == nil {
		return &Syscall{}
	}
	s1 := &Syscall{
		Name:              s.Name,
		Cmd:               s.Cmd,
		CriticalCondition: map[uint32]*Condition{},
		RunTimeDate:       CloneRunTimeData(s.RunTimeDate),
		WriteAddress:      []*WriteAddress{},
	}

	for i, c := range s.CriticalCondition {
		s1.CriticalCondition[i] = CloneCondition(c)
	}

	for _, wa := range s.WriteAddress {
		s1.WriteAddress = append(s1.WriteAddress, CloneWriteAddress(wa))
	}

	return s1
}

func CloneCondition(c *Condition) *Condition {
	if c == nil {
		return &Condition{}
	}
	c1 := &Condition{
		ConditionAddress:            c.ConditionAddress,
		SyzkallerConditionAddress:   c.SyzkallerConditionAddress,
		UncoveredAddress:            c.UncoveredAddress,
		SyzkallerUncoveredAddress:   c.SyzkallerUncoveredAddress,
		Idx:                         c.Idx,
		Successor:                   c.Successor,
		RightBranchAddress:          []uint64{},
		SyzkallerRightBranchAddress: []uint32{},
		WrongBranchAddress:          []uint64{},
		SyzkallerwrongBranchAddress: []uint32{},
	}

	for _, a := range c.RightBranchAddress {
		c1.RightBranchAddress = append(c1.RightBranchAddress, a)
	}

	for _, a := range c.SyzkallerRightBranchAddress {
		c1.SyzkallerRightBranchAddress = append(c1.SyzkallerRightBranchAddress, a)
	}

	for _, a := range c.WrongBranchAddress {
		c1.WrongBranchAddress = append(c1.WrongBranchAddress, a)
	}

	for _, a := range c.SyzkallerwrongBranchAddress {
		c1.SyzkallerwrongBranchAddress = append(c1.SyzkallerwrongBranchAddress, a)
	}

	return c1
}

// not finish
func (m *RunTimeData) MargeRunTimeData(d *RunTimeData) {
	if d == nil {
		return
	}

	return
}

func CloneRunTimeData(d *RunTimeData) *RunTimeData {
	if d == nil {
		return &RunTimeData{}
	}
	d1 := &RunTimeData{
		Program:                 []byte{},
		TaskStatus:              d.TaskStatus,
		RcursiveCount:           d.RcursiveCount,
		Idx:                     d.Idx,
		CheckCondition:          d.CheckCondition,
		ConditionAddress:        d.ConditionAddress,
		CheckAddress:            d.CheckAddress,
		Address:                 d.Address,
		CheckRightBranchAddress: d.CheckRightBranchAddress,
		RightBranchAddress:      []uint32{},
	}

	for _, c := range d.Program {
		d1.Program = append(d1.Program, c)
	}

	for _, a := range d.RightBranchAddress {
		d1.RightBranchAddress = append(d1.RightBranchAddress, a)
	}

	return d1
}

func (ss *Server) SetAddress(address uint32) {
	ss.address = address
}

// RunDependencyRPCServer
func (ss *Server) RunDependencyRPCServer(corpus *map[string]rpctype.RPCInput) {

	//ss.corpusDC = []*Input{}
	ss.corpusDC = make(map[string]*Input)
	ss.corpusDependency = &Corpus{
		CorpusDependencyInput: map[string]*Input{},
		CorpusRecursiveInput:  map[string]*Input{},
		CorpusErrorInput:      map[string]*Input{},
		WriteAddress:          map[uint64]*WriteAddresses{},
	}
	ss.fuzzers = make(map[string]*fuzzer)
	ss.mu = &sync.Mutex{}
	ss.fmu = &sync.Mutex{}
	ss.corpus = corpus

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ss.Dport = lis.Addr().(*net.TCPAddr).Port
	s := grpc.NewServer()
	RegisterDependencyRPCServer(s, ss)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (ss *Server) writeToDisk() {
	// Write the new address book back to disk.
	out, err := proto.Marshal(ss.corpusDependency)
	if err != nil {
		log.Fatalf("Failed to encode address:", err)
	}
	if err := ioutil.WriteFile("data.txt", out, 0644); err != nil {
		log.Fatalf("Failed to write address:", err)
	}
	// [END marshal_proto]
}
