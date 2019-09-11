package dra

import (
	"context"
	"github.com/golang/protobuf/proto"
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

// SendNewInput ...
func (d *DRPCClient) SendNewInput(input *Input) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err := d.c.SendNewInput(ctx, input, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Logf(0, "Dependency gRPC could not SendNewInput: %v", err)
		//log.Fatalf("Dependency gRPC could not SendNewInput: %v", err)
	}
}

func (d *DRPCClient) GetTasks(name string) *Tasks {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	request := &Empty{
		Name: name,
	}

	replay, err := d.c.GetTasks(ctx, request, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not SendNewInput: %v", err)
	}
	res := proto.Clone(replay).(*Tasks)
	return res
}

func (d *DRPCClient) ReturnTasks(task *Tasks) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	task.Name = *d.name
	_, err := d.c.ReturnTasks(ctx, task, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not SendNewInput: %v", err)
	}
	return
}

func (d *DRPCClient) SendUnstableInput(unstableInput *UnstableInput) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err := d.c.SendUnstableInput(ctx, unstableInput, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Logf(0, "Dependency gRPC could not SendNewInput: %v", err)
		//log.Fatalf("Dependency gRPC could not SendNewInput: %v", err)
	}
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
		//reply.input[i.Sig] = CloneInput(i)
		reply.Input = append(reply.Input, proto.Clone(i).(*Input))
	}
	return reply
}

// SendDependencyInput is
func (d *DRPCClient) ReturnDependencyInput(input *Input) (*Empty, error) {
	request := &Dependencytask{
		Input: input,
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

// SendNewInput ...
func (d *DRPCClient) SendLog(log string) {
	// Contact the server and print out its response.
	d.logMu.Lock()
	timeStr := time.Now().Format("2006/01/02 15:04:05 ")
	d.log = d.log + timeStr + *(d.name) + " syz_fuzzer : " + log + "\n"
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

func (d *DRPCClient) SendStat(stat *Statistic) (*Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	log.Logf(1, "Dependency gRPC sendStat stat : %v", stat)
	reply, err := d.c.SendStat(ctx, stat, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not sendStat: %v", err)
	}
	log.Logf(1, "Dependency gRPC sendStat reply.Name : %v", reply.Name)
	return reply, nil
}
