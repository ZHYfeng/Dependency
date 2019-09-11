package dra

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/syzkaller/pkg/log"
	"github.com/google/syzkaller/pkg/rpctype"
	"google.golang.org/grpc"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

const (
	startTime  = 21600
	taskNum    = 100
	DebugLevel = 2
)

type syzFuzzer struct {
	taskMu     *sync.Mutex
	newTask    Tasks
	returnTask Tasks
}

type newStats struct {
	newStat []*Statistic
}

type Dependencys struct {
	newDependency []*Dependency
}

// server is used to implement dra.DependencyServer.
type Server struct {
	address uint32
	Port    int
	Address string

	taskIndex        int
	corpusDependency *Corpus
	stat             *Statistics

	fuzzerMu *sync.Mutex
	fuzzers  map[string]*syzFuzzer

	corpus    *map[string]rpctype.RPCInput
	timeStart time.Time

	logMu *sync.Mutex
	log   *Empty

	statMu  *sync.Mutex
	newStat *newStats

	dependencyMu  *sync.Mutex
	newDependency *Dependencys

	newInputMu *sync.Mutex
	newInput   *Inputs

	inputMu *sync.Mutex
	input   *Inputs

	coveredInputMu *sync.Mutex
	coveredInput   *Inputs
}

func (ss Server) GetVmOffsets(context.Context, *Empty) (*Empty, error) {
	reply := &Empty{}
	reply.Address = ss.address
	return reply, nil
}

func (ss Server) SendBasicBlockNumber(ctx context.Context, request *Empty) (*Empty, error) {
	ss.stat.BasicBlockNumber = request.Address
	reply := &Empty{}
	return reply, nil
}

func (ss Server) GetNewInput(context.Context, *Empty) (*Inputs, error) {
	log.Logf(DebugLevel, "(ss Server) GetNewInput")

	reply := &Inputs{
		Input: []*Input{},
	}

	ss.newInputMu.Lock()
	last := len(ss.newInput.Input)
	log.Logf(DebugLevel, "(ss Server) GetNewInput len of newInput : %v", last)
	if last > 0 {
		last = last - 1
		reply.Input = append(reply.Input, ss.newInput.Input[last])
		ss.newInput.Input = ss.newInput.Input[:last]
	}
	ss.newInputMu.Unlock()

	ss.inputMu.Lock()
	ss.input.Input = append(ss.input.Input, reply.Input...)
	ss.inputMu.Unlock()

	return reply, nil
}

func (ss Server) SendDependency(ctx context.Context, request *Dependency) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendDependency")
	d := proto.Clone(request).(*Dependency)

	ss.dependencyMu.Lock()
	ss.newDependency.newDependency = append(ss.newDependency.newDependency, d)
	ss.dependencyMu.Unlock()

	reply := &Empty{}

	return reply, nil
}

func (ss Server) GetCondition(context.Context, *Empty) (*Conditions, error) {
	log.Logf(DebugLevel, "(ss Server) GetCondition")
	reply := &Conditions{
		//Condition: map[uint64]*Condition{},
		Condition: []*Condition{},
	}

	return reply, nil
}

func (ss Server) SendWriteAddress(ctx context.Context, request *WriteAddresses) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendWriteAddress")

	return &Empty{}, nil
}

func (ss Server) Connect(ctx context.Context, request *Empty) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) Connect")

	name := request.Name
	ss.fuzzerMu.Lock()
	defer ss.fuzzerMu.Unlock()

	_, ok := ss.fuzzers[name]
	if !ok {
		ss.fuzzers[name] = &syzFuzzer{
			taskMu: &sync.Mutex{},
			newTask: Tasks{
				Name: name,
				Task: []*Task{},
			},
			returnTask: Tasks{
				Name: name,
				Task: []*Task{},
			},
		}
	} else {

	}
	return &Empty{}, nil
}

func (ss Server) SendNewInput(ctx context.Context, request *Input) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendNewInput")

	reply := &Empty{}
	r := proto.Clone(request).(*Input)

	ss.newInputMu.Lock()
	ss.newInput.Input = append(ss.newInput.Input, r)
	last := len(ss.newInput.Input)
	log.Logf(DebugLevel, "(ss Server) SendNewInput len of newInput : %v", last)
	log.Logf(DebugLevel, "(ss Server) SendNewInput newInput : %v", r)
	ss.newInputMu.Unlock()

	ss.coveredInputMu.Lock()
	ss.coveredInput.Input = append(ss.coveredInput.Input, r)
	ss.coveredInputMu.Unlock()

	return reply, nil
}

func (ss Server) GetTasks(ctx context.Context, request *Empty) (*Tasks, error) {
	log.Logf(DebugLevel, "(ss Server) GetTasks")

	name := request.Name
	tasks := ss.pickTask(name)

	return tasks, nil
}

func (ss Server) ReturnTasks(ctx context.Context, request *Tasks) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) ReturnTasks")
	tasks := proto.Clone(request).(*Tasks)

	f, ok := ss.fuzzers[tasks.Name]
	if ok {
		f.taskMu.Lock()
		f.returnTask.Task = append(f.returnTask.Task, tasks.Task...)
		f.taskMu.Unlock()
	} else {
		log.Fatalf("ReturnTasks with error name")
	}
	reply := &Empty{}
	return reply, nil
}

func (ss Server) SendUnstableInput(ctx context.Context, request *UnstableInput) (*Empty, error) {
	ss.logMu.Lock()
	defer ss.logMu.Unlock()
	ss.log.Name = ss.log.Name + fmt.Sprintf("(ss Server) SendUnstableInput : %x\n", request.NewPath.Address)
	ss.log.Name = ss.log.Name + fmt.Sprintf("(ss Server) SendUnstableInput : %x\n", request.UnstablePath.Address)
	newPathIdx, unstablePathIdx, idx := CheckPath(request.NewPath.Address, request.UnstablePath.Address)
	ss.log.Name = ss.log.Name + fmt.Sprintf("(ss Server) SendUnstableInput newPathIdx: %v unstablePathIdx : %v idx : %v\n",
		newPathIdx, unstablePathIdx, idx)
	ss.log.Name = ss.log.Name + fmt.Sprintf("(ss Server) SendUnstableInput different address: %x\n",
		request.NewPath.Address[idx+newPathIdx])
	ss.log.Name = ss.log.Name + fmt.Sprintf("(ss Server) SendUnstableInput different address: %x\n",
		request.UnstablePath.Address[idx+unstablePathIdx])

	reply := &Empty{}
	return reply, nil
}

func (ss Server) GetDependencyInput(ctx context.Context, request *Empty) (*Inputs, error) {
	log.Logf(DebugLevel, "(ss Server) GetDependencyInput")
	reply := &Inputs{
		Input: []*Input{},
	}
	return reply, nil
}

func (ss Server) ReturnDependencyInput(ctx context.Context, request *Dependencytask) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) ReturnDependencyInput")
	reply := &Empty{}

	return reply, nil
}

func (ss Server) SendLog(ctx context.Context, request *Empty) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendLog")

	ss.logMu.Lock()
	defer ss.logMu.Unlock()

	var name = make([]uint8, len(request.Name))
	copy(name, request.Name)

	ss.log.Name = ss.log.Name + string(name)

	reply := &Empty{}
	return reply, nil
}

func (ss Server) SendStat(ctx context.Context, request *Statistic) (*Empty, error) {

	stat := proto.Clone(request).(*Statistic)
	ss.statMu.Lock()
	ss.newStat.newStat = append(ss.newStat.newStat, stat)
	ss.statMu.Unlock()

	reply := &Empty{}
	return reply, nil
}

func (ss *Server) SetAddress(address uint32) {
	ss.address = address
}

func (ss *Server) SyncSignal(signalNum uint64) {
	ss.stat.SignalNum = signalNum
}

// RunDependencyRPCServer
func (ss *Server) RunDependencyRPCServer(corpus *map[string]rpctype.RPCInput) {
	ss.taskIndex = 0

	ss.corpusDependency = &Corpus{
		Input:            map[string]*Input{},
		UncoveredAddress: map[uint32]*UncoveredAddress{},
		CoveredAddress:   map[uint32]*UncoveredAddress{},
		WriteAddress:     map[uint32]*WriteAddress{},
		IoctlCmd:         map[uint64]*IoctlCmd{},
		Tasks:            &Tasks{Name: "", Task: []*Task{}},
		NewInput:         map[string]*Input{},
	}

	ss.fuzzerMu = &sync.Mutex{}
	ss.fuzzers = make(map[string]*syzFuzzer)

	ss.stat = &Statistics{
		SignalNum:        0,
		BasicBlockNumber: 0,
		Coverage:         &Coverage{Coverage: map[uint32]uint32{}, Time: []*Time{}},
		Stat:             map[int32]*Statistic{},
		UsefulInput:      []*UsefulInput{},
	}

	ss.corpus = corpus
	ss.timeStart = time.Now()

	ss.logMu = &sync.Mutex{}
	ss.log = &Empty{
		Address: 0,
		Name:    "",
	}

	ss.statMu = &sync.Mutex{}
	ss.newStat = &newStats{newStat: []*Statistic{}}

	ss.dependencyMu = &sync.Mutex{}
	ss.newDependency = &Dependencys{newDependency: []*Dependency{}}

	ss.newInputMu = &sync.Mutex{}
	ss.newInput = &Inputs{Input: []*Input{}}

	ss.inputMu = &sync.Mutex{}
	ss.input = &Inputs{Input: []*Input{}}

	ss.coveredInputMu = &sync.Mutex{}
	ss.coveredInput = &Inputs{Input: []*Input{}}

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

func (ss *Server) Update() {

	// deal covered input
	ss.coveredInputMu.Lock()
	coveredInput := append([]*Input{}, ss.coveredInput.Input...)
	ss.coveredInput = &Inputs{Input: []*Input{}}
	ss.coveredInputMu.Unlock()

	for _, i := range coveredInput {
		ss.addCoveredAddress(i)
	}

	// deal new input
	ss.inputMu.Lock()
	input := append([]*Input{}, ss.input.Input...)
	ss.input = &Inputs{Input: []*Input{}}
	ss.inputMu.Unlock()

	for _, i := range input {
		ss.addInput(i)
	}

	// deal dependency
	for _, d := range ss.newDependency.newDependency {
		for _, wa := range d.WriteAddress {
			ss.addWriteAddress(wa)
		}
		ss.addUncoveredAddress(d.UncoveredAddress)
		ss.addInput(d.Input)
		ss.addInputTask(d.Input)
	}

	// deal retrun tasks
	var returnTask []*Task
	for _, f := range ss.fuzzers {
		f.taskMu.Lock()
		returnTask = append(returnTask, f.returnTask.Task...)
		f.returnTask.Task = []*Task{}
		f.taskMu.Unlock()
	}
	for _, task := range returnTask {
		for _, t := range ss.corpusDependency.Tasks.Task {
			if t.Sig == task.Sig && t.Index == task.Index &&
				t.WriteSig == task.WriteSig && t.WriteIndex == task.WriteIndex {
				t.MergeTask(task)
				break
			}
		}
	}

	sort.Slice(ss.corpusDependency.Tasks.Task, func(i, j int) bool {
		return ss.corpusDependency.Tasks.Task[i].Priority > ss.corpusDependency.Tasks.Task[j].Priority
	})

	// get new tasks
	t := time.Now()
	elapsed := t.Sub(ss.timeStart)
	if elapsed.Seconds() > startTime {
		var task []*Task
		for _, t := range ss.corpusDependency.Tasks.Task {
			if (t.TaskStatus == TaskStatus_untested || t.TaskStatus == TaskStatus_testing) && len(t.UncoveredAddress) > 0 {
				t.TaskStatus = TaskStatus_testing
				task = append(task, t)
			}
		}
		for _, f := range ss.fuzzers {
			f.taskMu.Lock()
			f.newTask.Task = append([]*Task{}, task...)
			f.taskMu.Unlock()
		}
	}

	ss.logMu.Lock()
	var templog = ss.log.Name
	ss.log.Name = ""
	ss.logMu.Unlock()
	f, _ := os.OpenFile("./dependency.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	_, _ = f.WriteString(string(templog))
	f.Close()

	ss.statMu.Lock()
	newStat := append([]*Statistic{}, ss.newStat.newStat...)
	ss.newStat = &newStats{newStat: []*Statistic{}}
	ss.statMu.Unlock()

	for _, stat := range newStat {
		s, ok := ss.stat.Stat[int32(stat.Name)]
		if ok {
			s.MergeStatistic(stat)
		} else {
			ss.stat.Stat[int32(stat.Name)] = stat
		}
	}

	ss.writeStatisticsToDisk()
	ss.writeCorpusToDisk()

}
