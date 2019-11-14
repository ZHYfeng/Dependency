package dra

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/log"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/rpctype"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
)

// useful const
const (
	// startTime  = 21600
	startTime  = 0
	newTime    = 600
	bootTime   = 300
	taskNum    = 20
	DebugLevel = 2
)

type syzFuzzer struct {
	taskMu         *sync.Mutex
	bootTasks      *Tasks
	highTasks      *Tasks
	newTask        *Tasks
	returnTask     *Tasks
	returnBootTask *Tasks
}

type newStats struct {
	newStat []*Statistic
}

type dependencys struct {
	newDependency []*Dependency
}

// Server is used to implement dra.DependencyServer.
type Server struct {
	address    uint32
	Port       int
	Address    string
	Dependency bool

	taskIndex        int
	corpusDependency *Corpus
	stat             *Statistics

	fuzzerMu *sync.Mutex
	fuzzers  map[string]*syzFuzzer

	corpus           *map[string]rpctype.RPCInput
	timeStart        time.Time
	timeNew          time.Time
	needWriteaddress bool
	needboot         bool

	logMu *sync.Mutex
	log   *Empty

	statMu  *sync.Mutex
	newStat *newStats

	dependencyMu  *sync.Mutex
	newDependency *dependencys

	newInputMu *sync.Mutex
	newInput   *Inputs

	// not new input, it is random picked inputs which used as new input.
	needInputMu *sync.Mutex
	needInput   *Inputs

	inputMu *sync.Mutex
	input   *Inputs

	coveredInputMu *sync.Mutex
	coveredInput   *Inputs
}

// GetVMOffsets is to send the offset address in vmlinux to dra
func (ss Server) GetVMOffsets(context.Context, *Empty) (*Empty, error) {
	reply := &Empty{}
	reply.Address = ss.address
	return reply, nil
}

// SendBasicBlockNumber is to get the basic block number from dra
func (ss Server) SendBasicBlockNumber(ctx context.Context, request *Empty) (*Empty, error) {
	ss.stat.BasicBlockNumber = request.Address
	reply := &Empty{}
	return reply, nil
}

// GetNewInput is to send new input ro dra
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
	for _, i := range reply.Input {
		ss.input.Input = append(ss.input.Input, proto.Clone(i).(*Input))
	}
	ss.inputMu.Unlock()

	return reply, nil
}

// SendDependency is to get depednency information from dra
func (ss Server) SendDependency(ctx context.Context, request *Dependency) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendDependency")
	d := proto.Clone(request).(*Dependency)

	ss.dependencyMu.Lock()
	ss.newDependency.newDependency = append(ss.newDependency.newDependency, d)
	ss.dependencyMu.Unlock()

	reply := &Empty{}

	return reply, nil
}

// GetCondition is to send condition to dra
func (ss Server) GetCondition(context.Context, *Empty) (*Conditions, error) {
	log.Logf(DebugLevel, "(ss Server) GetCondition")
	reply := &Conditions{
		//Condition: map[uint64]*Condition{},
		Condition: []*Condition{},
	}

	return reply, nil
}

// SendWriteAddress is to get write address for the condition from dra
func (ss Server) SendWriteAddress(ctx context.Context, request *WriteAddresses) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendWriteAddress")

	return &Empty{}, nil
}

// Connect is to connect with syz-fuzzer
func (ss Server) Connect(ctx context.Context, request *Empty) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) Connect")

	name := request.Name
	ss.fuzzerMu.Lock()
	defer ss.fuzzerMu.Unlock()

	_, ok := ss.fuzzers[name]
	if !ok {
		ss.fuzzers[name] = &syzFuzzer{
			taskMu: &sync.Mutex{},
			bootTasks: &Tasks{
				Name:  name,
				Task:  map[string]*Task{},
				Tasks: []*Task{},
			},
			highTasks: &Tasks{
				Name:  name,
				Task:  map[string]*Task{},
				Tasks: []*Task{},
			},
			newTask: &Tasks{
				Name:  name,
				Task:  map[string]*Task{},
				Tasks: []*Task{},
			},
			returnTask: &Tasks{
				Name:  name,
				Task:  map[string]*Task{},
				Tasks: []*Task{},
			},
			returnBootTask: &Tasks{
				Name:  name,
				Task:  map[string]*Task{},
				Tasks: []*Task{},
			},
		}
	} else {

	}
	return &Empty{}, nil
}

// SendNewInput is get new input from syz-fuzzer
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

// GetTasks ...
func (ss Server) GetTasks(ctx context.Context, request *Empty) (*Tasks, error) {
	log.Logf(DebugLevel, "(ss Server) GetTasks")

	name := request.Name
	tasks := ss.pickTask(name)

	return tasks, nil
}

// GetBootTasks for the tasks need to be tested when boot
func (ss Server) GetBootTasks(ctx context.Context, request *Empty) (*Tasks, error) {
	log.Logf(DebugLevel, "(ss Server) GetTasks")

	name := request.Name
	tasks := ss.pickBootTask(name)

	return tasks, nil
}

// ReturnTasks is to retrun the tasks from syz-fuzzer
func (ss Server) ReturnTasks(ctx context.Context, request *Tasks) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) ReturnTasks")
	tasks := proto.Clone(request).(*Tasks)

	f, ok := ss.fuzzers[tasks.Name]
	if ok {
		if tasks.Kind == TaskKind_Normal || tasks.Kind == TaskKind_High {
			f.taskMu.Lock()
			f.returnTask.addTasks(tasks)
			f.taskMu.Unlock()
		} else if tasks.Kind == TaskKind_Boot {
			f.taskMu.Lock()
			f.returnBootTask.addTasks(tasks)
			f.taskMu.Unlock()
		}
	} else {
		log.Fatalf("ReturnTasks with error name")
	}
	reply := &Empty{}
	return reply, nil
}

// SendBootInput is get new input from syz-fuzzer
func (ss Server) SendBootInput(ctx context.Context, request *Input) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendBootInput")
	reply := &Empty{}
	r := proto.Clone(request).(*Input)
	ss.coveredInputMu.Lock()
	ss.coveredInput.Input = append(ss.coveredInput.Input, r)
	ss.coveredInputMu.Unlock()
	return reply, nil
}

// SendUnstableInput is get unstable input from syz-fuzzer
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

// SendLog is to get log from syz-fuzzer
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

// SendStat is to get stat from suz-fuzzer
func (ss Server) SendStat(ctx context.Context, request *Statistic) (*Empty, error) {

	stat := proto.Clone(request).(*Statistic)
	ss.statMu.Lock()
	ss.newStat.newStat = append(ss.newStat.newStat, stat)
	ss.statMu.Unlock()

	reply := &Empty{}
	return reply, nil
}

// GetNeed is to random get input from syz-fuzzer, not new input but used as new input.
func (ss Server) GetNeed(ctx context.Context, request *Empty) (*Empty, error) {

	reply := &Empty{}
	if ss.needWriteaddress {
		reply.Address = 1
	} else {
		reply.Address = 0
	}
	return reply, nil
}

// SendNeedInput is to random get input from syz-fuzzer, not new input but used as new input.
func (ss Server) SendNeedInput(ctx context.Context, request *Input) (*Empty, error) {
	reply := &Empty{}
	r := proto.Clone(request).(*Input)

	ss.needInputMu.Lock()
	ss.needInput.Input = append(ss.needInput.Input, r)
	ss.needInputMu.Unlock()

	return reply, nil
}

// SetAddress is to set the port address for Server
func (ss *Server) SetAddress(address uint32) {
	ss.address = address
}

// SyncSignal is to sync the number of signal
func (ss *Server) SyncSignal(signalNum uint64) {
	ss.stat.SignalNum = signalNum
}

// RunDependencyRPCServer : run the server
func (ss *Server) RunDependencyRPCServer(corpus *map[string]rpctype.RPCInput) {
	ss.taskIndex = 0

	ss.corpusDependency = &Corpus{
		Input:            map[string]*Input{},
		UncoveredAddress: map[uint32]*UncoveredAddress{},
		CoveredAddress:   map[uint32]*UncoveredAddress{},
		WriteAddress:     map[uint32]*WriteAddress{},
		FileOperations:   map[string]*FileOperations{},
		Tasks:            &Tasks{Name: "", Task: map[string]*Task{}, Tasks: []*Task{}},
		HighTask:         &Tasks{Name: "", Task: map[string]*Task{}, Tasks: []*Task{}},
		BootTask:         &Tasks{Name: "", Task: map[string]*Task{}, Tasks: []*Task{}},
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
	ss.timeNew = time.Now()
	ss.needWriteaddress = false

	ss.logMu = &sync.Mutex{}
	ss.log = &Empty{
		Address: 0,
		Name:    "",
	}

	ss.statMu = &sync.Mutex{}
	ss.newStat = &newStats{newStat: []*Statistic{}}

	ss.dependencyMu = &sync.Mutex{}
	ss.newDependency = &dependencys{newDependency: []*Dependency{}}

	ss.newInputMu = &sync.Mutex{}
	ss.newInput = &Inputs{Input: []*Input{}}

	ss.needInputMu = &sync.Mutex{}
	ss.needInput = &Inputs{Input: []*Input{}}

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

// Update : update the information in the server
func (ss *Server) Update() {

	// deal covered input
	ss.coveredInputMu.Lock()
	coveredInput := append([]*Input{}, ss.coveredInput.Input...)
	ss.coveredInput = &Inputs{Input: []*Input{}}
	ss.coveredInputMu.Unlock()
	for _, i := range coveredInput {
		ss.addCoveredAddress(i)
	}
	coveredInput = nil

	// deal new input
	ss.inputMu.Lock()
	input := append([]*Input{}, ss.input.Input...)
	ss.input = &Inputs{Input: []*Input{}}
	ss.inputMu.Unlock()
	for _, i := range input {
		ss.addInput(i)
	}
	if len(input) == 0 {
		t := time.Now()
		elapsed := t.Sub(ss.timeNew)
		if elapsed.Seconds() > newTime {
			ss.needWriteaddress = true
		}
		if elapsed.Seconds() > bootTime {
			ss.needboot = true
		}
	} else {
		ss.timeNew = time.Now()
		ss.needWriteaddress = false
		ss.needboot = false
	}
	input = nil

	// reboot the qemu
	if ss.needboot {
		ss.needboot = false
	}

	// deal need input
	ss.needInputMu.Lock()
	needInput := append([]*Input{}, ss.needInput.Input...)
	ss.needInput = &Inputs{Input: []*Input{}}
	ss.needInputMu.Unlock()
	for _, i := range needInput {
		ss.addInput(i)
	}
	needInput = nil

	// deal Dependency
	ss.dependencyMu.Lock()
	newDependency := append([]*Dependency{}, ss.newDependency.newDependency...)
	ss.newDependency.newDependency = []*Dependency{}
	ss.dependencyMu.Unlock()
	for _, d := range newDependency {
		for _, wa := range d.WriteAddress {
			ss.addWriteAddress(wa)
		}
		ss.addUncoveredAddress(d.UncoveredAddress)
		ss.addInput(d.Input)
		ss.addInputTask(d.Input)
	}
	newDependency = nil

	// deal return tasks
	returnTask := &Tasks{Name: "", Task: map[string]*Task{}, Tasks: []*Task{}}
	for _, f := range ss.fuzzers {
		f.taskMu.Lock()
		returnTask.addTasks(f.returnTask)
		f.returnTask.emptyTask()
		f.taskMu.Unlock()
	}
	for hash, task := range returnTask.Task {
		if t, ok := ss.corpusDependency.Tasks.Task[hash]; ok {
			if task.TaskStatus == TaskStatus_unstable && t.TaskStatus == TaskStatus_testing {
				ss.reducePriority(task)
			}
			if task.TaskStatus == TaskStatus_tested {
				ss.reducePriority(task)
			}
			t.mergeTask(task)
			for u := range t.UncoveredAddress {
				_, ok := ss.corpusDependency.UncoveredAddress[u]
				if ok {

				} else {
					delete(t.UncoveredAddress, u)
				}
			}
		}
	}
	sort.Slice(ss.corpusDependency.Tasks.Tasks, func(i, j int) bool {
		return ss.corpusDependency.Tasks.Tasks[i].getRealPriority() > ss.corpusDependency.Tasks.Tasks[j].getRealPriority()
	})
	returnTask = nil

	// get new tasks
	if ss.Dependency {
		t := time.Now()
		elapsed := t.Sub(ss.timeStart)
		if elapsed.Seconds() > startTime {
			if len(ss.corpusDependency.HighTask.Task) != 0 {
				var task []*Task
				for _, t := range ss.corpusDependency.HighTask.Task {
					for u := range t.UncoveredAddress {
						_, ok := ss.corpusDependency.UncoveredAddress[u]
						if ok {

						} else {
							delete(t.UncoveredAddress, u)
						}
					}
					if len(t.UncoveredAddress) > 0 {
						task = append(task, t)
					}
				}
				ss.corpusDependency.HighTask.emptyTask()
				task = []*Task{}
				for _, f := range ss.fuzzers {
					f.taskMu.Lock()
					for _, t := range task {
						f.highTasks.addTask(proto.Clone(t).(*Task))
					}
					f.taskMu.Unlock()
				}
				task = nil
			} else {
				var task []*Task
				for _, t := range ss.corpusDependency.Tasks.Task {
					for u := range t.UncoveredAddress {
						_, ok := ss.corpusDependency.UncoveredAddress[u]
						if ok {

						} else {
							delete(t.UncoveredAddress, u)
						}
					}
					if len(t.UncoveredAddress) > 0 {
						if t.TaskStatus == TaskStatus_untested {
							t.TaskStatus = TaskStatus_testing
							task = append(task, t)
						} else if t.TaskStatus == TaskStatus_testing {
							task = append(task, t)
						} else if t.TaskStatus == TaskStatus_unstable {
							task = append(task, t)
						}
						if len(task) > taskNum {
							break
						}
					}
				}
				for _, f := range ss.fuzzers {
					f.taskMu.Lock()
					for _, t := range task {
						f.newTask.addTask(proto.Clone(t).(*Task))
					}
					f.taskMu.Unlock()
				}
				task = nil
			}
		}
	} else {
		ss.corpusDependency.HighTask.emptyTask()
		ss.corpusDependency.Tasks.emptyTask()
	}

	// deal return boot tasks
	returnBootTask := &Tasks{Name: "", Task: map[string]*Task{}, Tasks: []*Task{}}
	for _, f := range ss.fuzzers {
		f.taskMu.Lock()
		returnBootTask.addTasks(f.returnBootTask)
		f.returnBootTask.emptyTask()
		f.taskMu.Unlock()
	}
	for hash, task := range returnBootTask.Task {
		if t, ok := ss.corpusDependency.BootTask.Task[hash]; ok {
			if task.TaskStatus == TaskStatus_covered {
				t.mergeTask(task)
			} else {
				t.TaskStatus = TaskStatus_tested
			}
			t.mergeTask(task)
			for u := range t.UncoveredAddress {
				_, ok := ss.corpusDependency.UncoveredAddress[u]
				if ok {

				} else {
					delete(t.UncoveredAddress, u)
				}
			}
		}
	}
	sort.Slice(ss.corpusDependency.BootTask.Tasks, func(i, j int) bool {
		return ss.corpusDependency.BootTask.Tasks[i].getRealPriority() > ss.corpusDependency.BootTask.Tasks[j].getRealPriority()
	})
	returnBootTask = nil

	// get boot tasks
	if ss.Dependency {
		t := time.Now()
		elapsed := t.Sub(ss.timeStart)
		if elapsed.Seconds() > startTime {
			if len(ss.corpusDependency.BootTask.Task) != 0 {
				var task []*Task
				for _, t := range ss.corpusDependency.BootTask.Task {
					if t.TaskStatus == TaskStatus_untested {
						for u := range t.UncoveredAddress {
							_, ok := ss.corpusDependency.UncoveredAddress[u]
							if ok {

							} else {
								delete(t.UncoveredAddress, u)
							}
						}
						if len(t.UncoveredAddress) > 0 {
							task = append(task, t)
							t.TaskStatus = TaskStatus_testing
						}
					}
				}
				for _, f := range ss.fuzzers {
					f.taskMu.Lock()
					for _, t := range task {
						f.bootTasks.addTask(proto.Clone(t).(*Task))
					}
					f.taskMu.Unlock()
				}
				task = nil
			}
		}
	} else {
		ss.corpusDependency.BootTask.emptyTask()
	}

	ss.logMu.Lock()
	var templog = ss.log.Name
	ss.log.Name = ""
	ss.logMu.Unlock()
	f, _ := os.OpenFile("./Dependency.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	_, _ = f.WriteString(string(templog))
	_ = f.Close()

	ss.statMu.Lock()
	newStat := append([]*Statistic{}, ss.newStat.newStat...)
	ss.newStat = &newStats{newStat: []*Statistic{}}
	ss.statMu.Unlock()
	for _, stat := range newStat {
		s, ok := ss.stat.Stat[int32(stat.Name)]
		if ok {
			s.mergeStatistic(stat)
		} else {
			ss.stat.Stat[int32(stat.Name)] = stat
		}
	}
	newStat = nil

	ss.writeStatisticsToDisk()
	ss.writeCorpusToDisk()
}
