package dra

import (
	"context"
	"github.com/google/syzkaller/pkg/log"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type DRPCClient struct {
	c    DependencyRPCClient
	I    []*Input
	name *string

	log   string
	logMu sync.RWMutex
}

func (d *DRPCClient) RunDependencyRPCClient(address, name *string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dependency gRPC did not connect: %v", err)
	}
	d.c = NewDependencyRPCClient(conn)
	d.name = name
	d.Connect(name)
}

func (d *DRPCClient) Connect(name *string) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	_, err := d.c.Connect(ctx, &Empty{Name: *name})
	if err != nil {
		log.Fatalf("Dependency gRPC could not Connect: %v", err)
	}
	return
}

func (d *DRPCClient) GetDependencyInput(name string) *Inputs {
	// Contact the server and print out its response.
	request := &Empty{
		Name: name,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	dInputs, err := d.c.GetDependencyInput(ctx, request, grpc.MaxCallRecvMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not GetDependencyInput: %v", err)
	}
	reply := &Inputs{}
	for _, i := range dInputs.Input {
		reply.Input = append(reply.Input, CloneInput(i))
	}
	return reply
}

// SendDependencyInput is
func (d *DRPCClient) ReturnDependencyInput(input *Input) (*Empty, error) {
	request := &Task{
		Input: CloneInput(input),
		Name:  *d.name,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	reply, err := d.c.ReturnDependencyInput(ctx, request, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not ReturnDependencyInput: %v", err)
	}
	log.Logf(1, "Dependency gRPC ReturnDependencyInput reply.Name : %v", reply.Name)
	return reply, nil
}

// SendInput ...
func (d *DRPCClient) SendInput(input *Input) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err := d.c.SendNewInput(ctx, CloneInput(input), grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not SendInput: %v", err)
	}
	log.Logf(1, "Dependency gRPC SendInput sig : v%", input.Sig)
	for idx, cc := range input.Call {
		log.Logf(1, "Dependency gRPC SendInput idx: %v address : %x", idx, cc.Address)
	}
}

// SendInput ...
func (d *DRPCClient) SendLog(log string) {
	// Contact the server and print out its response.
	d.logMu.Lock()
	timeStr := time.Now().Format("2006/01/02 15:04:05 ")
	d.log = d.log + timeStr + *(d.name) + " fuzzer : " + log + "\n"
	d.logMu.Unlock()
	return
}

func (d *DRPCClient) SSendLog() {
	// Contact the server and print out its response.
	d.logMu.Lock()
	request := &Empty{
		Name: d.log,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	_, _ = d.c.SendLog(ctx, request, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	d.log = ""
	d.logMu.Unlock()
	return
}
