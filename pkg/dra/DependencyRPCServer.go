package dra

import (
	"context"
	"github.com/google/syzkaller/pkg/log"
	"google.golang.org/grpc"
	"net"
)

const (
	port = ":50051"
)

type fuzzer struct {
	corpusDI []DependencyInput
}

// server is used to implement dra.DependencyServer.
type Server struct {
	address  uint32
	Dport    int
	corpusDC map[string]Input
	corpusDI map[string]DependencyInput
	fuzzers  map[string]*fuzzer
}

func (ss Server) Connect(ctx context.Context, request *Empty) (*Empty, error) {
	ss.fuzzers[request.Name] = &fuzzer{}
	return &Empty{}, nil
}

func (ss Server) GetVmOffsets(context.Context, *Empty) (*Empty, error) {
	reply := &Empty{}
	reply.Address = ss.address
	return reply, nil
}

func (ss Server) GetNewInput(context.Context, *Empty) (*NewInput, error) {
	reply := &NewInput{}
	for _, c := range ss.corpusDC {
		reply.Input = append(reply.Input, cloneInput(&c))
	}
	ss.corpusDI = nil
	return reply, nil
}

func (ss Server) SendDependencyInput(ctx context.Context, request *DependencyInput) (*Empty, error) {
	for _, f := range ss.fuzzers {
		f.corpusDI = append(f.corpusDI, *cloneDependencyInput(request))
	}
	return &Empty{}, nil
}

func (ss Server) GetDependencyInput(ctx context.Context, request *Empty) (*NewDependencyInput, error) {
	f := ss.fuzzers[request.Name]
	if f == nil {
		log.Fatalf("fuzzer %v is not connected", request.Name)
	}
	reply := &NewDependencyInput{}

	for i := 0; i < 100 && len(f.corpusDI) > 0; i++ {
		last := len(f.corpusDI) - 1
		reply.DependencyInput = append(reply.DependencyInput, cloneDependencyInput(&f.corpusDI[last]))
		f.corpusDI = f.corpusDI[:last]
	}
	if len(f.corpusDI) == 0 {
		f.corpusDI = nil
	}
	return reply, nil
}

func (ss Server) SendInput(ctx context.Context, request *Input) (*Empty, error) {
	reply := &Empty{}
	input := Input{
		Sig:  request.Sig,
		Call: make(map[uint32]*Call),
	}
	for i, c := range request.Call {
		cc := &Call{Idx: c.Idx, Address: make(map[uint32]uint32)}
		input.Call[i] = cc
		for _, a := range c.Address {
			cc.Address[a] = 0
		}
	}
	ss.corpusDC[request.Sig] = input
	return reply, nil
}

func cloneDependencyInput(d *DependencyInput) *DependencyInput {
	cd := &DependencyInput{
		Sig:              d.Sig,
		UncoveredAddress: []*UncoveredAddress{},
	}
	for _, u := range d.UncoveredAddress {
		u1 := &UncoveredAddress{
			Address:          u.Address,
			Idx:              u.Idx,
			ConditionAddress: u.ConditionAddress,
			RelatedInput:     []*RelatedInput{},
			RelatedSyscall:   []*RelatedSyscall{},
		}
		for _, i := range u.RelatedInput {
			i1 := &RelatedInput{
				Sig:     i.Sig,
				Address: i.Address,
			}
			u1.RelatedInput = append(u1.RelatedInput, i1)
		}
		for _, s := range u.RelatedSyscall {
			s1 := &RelatedSyscall{
				Name:    s.Name,
				Address: s.Address,
			}
			u1.RelatedSyscall = append(u1.RelatedSyscall, s1)
		}
		cd.UncoveredAddress = append(cd.UncoveredAddress, u1)
	}
	return cd
}

func cloneInput(d *Input) *Input {
	ci := &Input{
		Sig:  d.Sig,
		Call: make(map[uint32]*Call),
	}
	for i, u := range d.Call {
		u1 := &Call{
			Address: u.Address,
			Idx:     u.Idx,
		}
		ci.Call[i] = u1
	}
	return ci
}

func (ss *Server) SetAddress(address uint32) {
	ss.address = address
}

// RunDependencyRPCServer
func (ss *Server) RunDependencyRPCServer() {

	ss.corpusDC = make(map[string]Input)
	ss.corpusDI = make(map[string]DependencyInput)
	ss.fuzzers = make(map[string]*fuzzer)

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
