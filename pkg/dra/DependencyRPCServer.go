package dra

import (
	"context"
	"github.com/google/syzkaller/pkg/log"
	"github.com/google/syzkaller/pkg/rpctype"
	"google.golang.org/grpc"
	"net"
	"sync"
)

const (
	port = ":50051"
)

type fuzzer struct {
	corpusDI map[string]*Input
}

// server is used to implement dra.DependencyServer.
type Server struct {
	address uint32
	Dport   int
	//corpusDC []*Input
	corpusDC map[string]*Input
	corpusDI map[string]*Input
	fmu      *sync.Mutex
	fuzzers  map[string]*fuzzer
	mu       *sync.Mutex
	corpus   *map[string]rpctype.RPCInput
}

func (ss Server) GetCondition(context.Context, *Empty) (*Condition, error) {
	panic("implement me")
}

func (ss Server) SendWriteAddress(context.Context, *WriteAddress) (*Empty, error) {
	panic("implement me")
}

func (ss Server) SendNewInput(context.Context, *Input) (*Empty, error) {
	panic("implement me")
}

func (ss Server) SendCondition(context.Context, *Condition) (*Empty, error) {
	panic("implement me")
}

func (ss Server) GetWriteAddress(context.Context, *Empty) (*WriteAddress, error) {
	panic("implement me")
}

func (ss Server) SendLog(ctx context.Context, request *Empty) (*Empty, error) {
	//ss.mu.Lock()
	//defer ss.mu.Unlock()
	log.Logff(1, request.Name)
	reply := &Empty{}
	return reply, nil
}

func (ss Server) Connect(ctx context.Context, request *Empty) (*Empty, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if _, ok := ss.fuzzers[request.Name]; !ok {
		ss.fuzzers[request.Name] = &fuzzer{
			corpusDI: map[string]*Input{},
		}
	}
	return &Empty{}, nil
}

func (ss Server) GetVmOffsets(context.Context, *Empty) (*Empty, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	reply := &Empty{}
	reply.Address = ss.address
	return reply, nil
}

func (ss Server) GetNewInput(context.Context, *Empty) (*Input, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	reply := &Input{}
	i := 0
	for s, c := range ss.corpusDC {
		if i < 1 {
			reply = CloneInput(c)
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
		reply.Name = "dependency Sig error : " + cd.Sig
		return reply, nil
	} else if len(cd.Sig) == 0 {
		reply.Name = "dependency Prog error : " + string(cd.Program)
		return reply, nil
	}

	ss.fmu.Lock()
	for _, f := range ss.fuzzers {
		f.corpusDI[cd.Sig] = CloneInput(cd)
		reply.Address = uint32(len(f.corpusDI))
	}
	reply.Name = "success"
	ss.fmu.Unlock()
	return reply, nil
}

func (ss Server) GetDependencyInput(ctx context.Context, request *Empty) (*Input, error) {
	ss.fmu.Lock()
	reply := &Input{}
	if f, ok := ss.fuzzers[request.Name]; ok {
		i := 0
		for s, c := range f.corpusDI {
			if i < 1 {
				reply = CloneInput(c)
				i++
				delete(f.corpusDI, s)
			} else {
			}
		}
	} else {
		log.Fatalf("fuzzer %v is not connected", request.Name)
	}

	//for i := 0; i < 50 && len(f.corpusDI) > 0; i++ {
	//	last := len(f.corpusDI) - 1
	//	reply.DependencyInput = append(reply.DependencyInput, cloneDependencyInput(&f.corpusDI[last]))
	//	f.corpusDI = f.corpusDI[:last]
	//}
	//if len(f.corpusDI) == 0 {
	//	f.corpusDI = nil
	//}
	ss.fmu.Unlock()
	return reply, nil
}

func (ss Server) SendInput(ctx context.Context, request *Input) (*Empty, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	reply := &Empty{}
	input := CloneInput(request)

	//ss.corpusDC = append(ss.corpusDC, input)
	ss.corpusDC[input.Sig] = input
	return reply, nil
}

func CloneInput(d *Input) *Input {
	ci := &Input{
		Sig:        d.Sig,
		Call:       make(map[uint32]*Call),
		Dependency: d.Dependency,
	}
	for i, u := range d.Call {
		u1 := &Call{
			Address: make(map[uint32]uint32),
			Idx:     u.Idx,
		}
		for aa := range u.Address {
			u1.Address[aa] = 0
		}
		ci.Call[i] = u1
	}
	copy(ci.Program, d.Program)

	for _, u := range d.UncoveredAddress {
		u1 := new(UncoveredAddress)
		u1.ConditionAddress = u.ConditionAddress
		u1.UncoveredAddress = u.UncoveredAddress
		for _, a := range u.WriteAddress {
			a1 := &WriteAddress{
				WriteAddress:     a.WriteAddress,
				ConditionAddress: a.ConditionAddress,
				Repeat:           a.Repeat,
				Prio:             a.Prio,
			}

			for _, i := range a.WriteInput {
				i1 := CloneInput(i)
				a1.WriteInput = append(a1.WriteInput, i1)
			}

			for _, s := range a.WriteSyscall {
				s1 := &Syscall{
					Name: s.Name,
					Cmd:  s.Cmd,
				}
				for _, c := range s.CriticalCondition {
					c1 := &Condition{
						ConditionAddress: c.ConditionAddress,
						UncoveredAddress: c.UncoveredAddress,
						Idx:              c.Idx,
					}
					for _, a := range c.RightBranchAddress {
						c1.RightBranchAddress = append(c1.RightBranchAddress, a)
					}
					for _, a := range c.WrongBranchAddress {
						c1.WrongBranchAddress = append(c1.WrongBranchAddress, a)
					}
					s1.CriticalCondition = append(s1.CriticalCondition, c1)
				}
				a1.WriteSyscall = append(a1.WriteSyscall, s1)
			}
			u1.WriteAddress = append(u1.WriteAddress, a1)
		}
		ci.UncoveredAddress = append(ci.UncoveredAddress, u1)
	}

	return ci
}

func (ss *Server) SetAddress(address uint32) {
	ss.address = address
}

// RunDependencyRPCServer
func (ss *Server) RunDependencyRPCServer(corpus *map[string]rpctype.RPCInput) {

	//ss.corpusDC = []*Input{}
	ss.corpusDC = make(map[string]*Input)
	ss.corpusDI = make(map[string]*Input)
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
