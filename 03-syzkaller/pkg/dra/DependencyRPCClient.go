package dra

import (
	"context"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/log"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// DRPCClient : the RPC client
type DRPCClient struct {
	name           *string
	c              DependencyRPCClient
	MuDependency   *sync.RWMutex
	DataDependency *DataDependency

	log   string
	logMu sync.RWMutex
}

// RunDependencyRPCClient : run the client
func (d *DRPCClient) RunDependencyRPCClient(address, name *string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dependency gRPC did not connect: %v", err)
	}
	d.name = name
	d.c = NewDependencyRPCClient(conn)
	d.Connect(name)
	if CheckCondition {
		d.MuDependency = &sync.RWMutex{}
		d.MuDependency.Lock()
		d.DataDependency = &DataDependency{
			Input: map[string]*Input{},
		}
		d.MuDependency.Unlock()
	}
}

// Connect : connect to syz-manager
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

func (d *DRPCClient) GetDataDependency() {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	replay, err := d.c.GetDataDependency(ctx, &Empty{Name: *d.name})
	if err != nil {
		log.Fatalf("Dependency gRPC could not Connect: %v", err)
	}
	d.MuDependency.Lock()
	d.DataDependency = proto.Clone(replay).(*DataDependency)
	if d.DataDependency.Input == nil {
		d.DataDependency.Input = map[string]*Input{}
	}
	d.MuDependency.Unlock()
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
	}
}

// GetTasks : get task from syz-manager
func (d *DRPCClient) GetTasks(name string) *Tasks {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	request := &Empty{
		Name: name,
	}

	replay, err := d.c.GetTasks(ctx, request, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not GetTasks: %v", err)
	}
	res := proto.Clone(replay).(*Tasks)
	return res
}

// GetBootTasks ...
func (d *DRPCClient) GetBootTasks(name string) *Tasks {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	request := &Empty{
		Name: name,
	}
	replay, err := d.c.GetBootTasks(ctx, request, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not GetBootTasks: %v", err)
	}
	res := proto.Clone(replay).(*Tasks)
	return res
}

// ReturnTasks : return the task to syz-manager
func (d *DRPCClient) ReturnTasks(task *Tasks) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	task.Name = *d.name
	_, err := d.c.ReturnTasks(ctx, task, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not ReturnTasks: %v", err)
	}
	return
}

// SendBootInput ...
func (d *DRPCClient) SendBootInput(input *Input) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err := d.c.SendBootInput(ctx, input, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Logf(0, "Dependency gRPC could not SendNewInput: %v", err)
	}
}

// SendUnstableInput : send unstable input to syz-manager
func (d *DRPCClient) SendUnstableInput(unstableInput *UnstableInput) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err := d.c.SendUnstableInput(ctx, unstableInput, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Logf(0, "Dependency gRPC could not SendUnstableInput: %v", err)
	}
}

// SendLog ...
func (d *DRPCClient) SendLog(log string) {
	// Contact the server and print out its response.
	d.logMu.Lock()
	timeStr := time.Now().Format("2006/01/02 15:04:05 ")
	d.log = d.log + timeStr + *(d.name) + " syz_fuzzer : " + log + "\n"
	d.logMu.Unlock()
	return
}

// SSendLog : real send log to syz-manager
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

// SendStat : send stat information to syz-manager
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

// GetNeed : to know whether syz-manager needs input
func (d *DRPCClient) GetNeed() (*Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	request := &Empty{}
	r, err := d.c.GetNeed(ctx, request, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Fatalf("Dependency gRPC could not sendStat: %v", err)
	}
	reply := proto.Clone(r).(*Empty)
	return reply, nil
}

// SendNeedInput : send random input to syz -manager
func (d *DRPCClient) SendNeedInput(input *Input) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err := d.c.SendNeedInput(ctx, input, grpc.MaxCallSendMsgSize(0x7fffffffffffffff))
	if err != nil {
		log.Logf(0, "Dependency gRPC could not SendNeedInput: %v", err)
		//log.Fatalf("Dependency gRPC could not SendNeedInput: %v", err)
	}
}
