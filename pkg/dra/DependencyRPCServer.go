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
	corpusDI map[string]*DependencyInput
}

// server is used to implement dra.DependencyServer.
type Server struct {
	address uint32
	Dport   int
	//corpusDC []*Input
	corpusDC map[string]*Input
	corpusDI map[string]*DependencyInput
	fmu      *sync.Mutex
	fuzzers  map[string]*fuzzer
	mu       *sync.Mutex
	corpus   *map[string]rpctype.RPCInput
}

func (ss Server) SendLog(ctx context.Context, request *Empty) (*Empty, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	log.Logff(1, request.Name)
	reply := &Empty{}
	return reply, nil
}

func (ss Server) Connect(ctx context.Context, request *Empty) (*Empty, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if _, ok := ss.fuzzers[request.Name]; !ok {
		ss.fuzzers[request.Name] = &fuzzer{
			corpusDI: map[string]*DependencyInput{},
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

func (ss Server) GetNewInput(context.Context, *Empty) (*NewInput, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	reply := &NewInput{}
	i := 0
	for s, c := range ss.corpusDC {
		if i < 50 {
			reply.Input = append(reply.Input, cloneInput(c))
			i++
			delete(ss.corpusDC, s)
		} else {
		}

	}
	return reply, nil
}

func (ss Server) SendDependencyInput(ctx context.Context, request *DependencyInput) (*Empty, error) {
	reply := &Empty{}
	cd := cloneDependencyInput(request)
	sig := cd.Sig
	if inp, ok := (*ss.corpus)[sig]; ok {
		for _, p := range inp.Prog {
			cd.Prog = append(cd.Prog, p)
		}
	} else {
		reply.Name = "dependency sig error : " + sig
		return reply, nil
	}
	for _, u := range cd.UncoveredAddress {
		for _, a := range u.RelatedAddress {
			for _, r := range a.RelatedInput {
				if rinp, ok := (*ss.corpus)[r.Sig]; ok {
					for _, p := range rinp.Prog {
						r.Prog = append(r.Prog, p)
					}
				} else {
					reply.Name = "related input sig error : " + r.Sig
					return reply, nil
				}
			}
		}
	}
	ss.fmu.Lock()
	for _, f := range ss.fuzzers {
		f.corpusDI[sig] = cloneDependencyInput(cd)
		reply.Address = uint32(len(f.corpusDI))
	}
	reply.Name = "success"
	ss.fmu.Unlock()
	return reply, nil
}

func (ss Server) GetDependencyInput(ctx context.Context, request *Empty) (*NewDependencyInput, error) {
	ss.fmu.Lock()
	reply := &NewDependencyInput{}
	if f, ok := ss.fuzzers[request.Name]; ok {
		i := 0
		for s, c := range f.corpusDI {
			if i < 10 {
				reply.DependencyInput = append(reply.DependencyInput, cloneDependencyInput(c))
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
	input := cloneInput(request)

	//ss.corpusDC = append(ss.corpusDC, input)
	ss.corpusDC[input.Sig] = input
	return reply, nil
}

func cloneDependencyInput(d *DependencyInput) *DependencyInput {
	cd := new(DependencyInput)
	cd.Sig = d.Sig
	for _, p := range d.Prog {
		cd.Prog = append(cd.Prog, p)
	}
	for _, u := range d.UncoveredAddress {
		u1 := new(UncoveredAddress)
		u1.Address = u.Address
		u1.Idx = u.Idx
		u1.ConditionAddress = u.ConditionAddress
		for _, a := range u.RelatedAddress {
			a1 := &RelatedAddress{
				Address: a.Address,
				Repeat:  a.Repeat,
				Prio:    a.Prio,
			}

			for _, i := range a.RelatedInput {
				i1 := &RelatedInput{
					Sig: i.Sig,
				}
				for _, p := range i.Prog {
					i1.Prog = append(cd.Prog, p)
				}
				a1.RelatedInput = append(a1.RelatedInput, i1)
			}

			for _, s := range a.RelatedSyscall {
				s1 := &RelatedSyscall{
					Name:   s.Name,
					Number: s.Number,
				}
				a1.RelatedSyscall = append(a1.RelatedSyscall, s1)
			}
			u1.RelatedAddress = append(u1.RelatedAddress, a1)
		}
		cd.UncoveredAddress = append(cd.UncoveredAddress, u1)
	}
	return cd
}

func cloneInput(d *Input) *Input {
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
	return ci
}

func (ss *Server) SetAddress(address uint32) {
	ss.address = address
}

// RunDependencyRPCServer
func (ss *Server) RunDependencyRPCServer(corpus *map[string]rpctype.RPCInput) {

	//ss.corpusDC = []*Input{}
	ss.corpusDC = make(map[string]*Input)
	ss.corpusDI = make(map[string]*DependencyInput)
	ss.fuzzers = make(map[string]*fuzzer)
	ss.mu = &sync.Mutex{}
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
