package dra

import (
	"context"
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
	//startTime = 10800
	startTime       = 0
	newTime         = 3600
	bootTime        = 3600
	TimeWriteToDisk = 3600
	TimeExit        = 3600 * 24

	TaskNum             = 40
	TaskCountLimitation = 30

	DebugLevel = 0

	CollectPath     = true
	CollectUnstable = true

	// collect coverage by intersection instead of union.
	StableCoverage = true
	// check Condition address in syz-fuzzer
	CheckCondition = true
)

const (
	NameDevice         = "dev_"
	NameBase           = "base"
	NameWithDra        = "01-result-with-dra"
	NameWithoutDra     = "02-result-without-dra"
	NameData           = "data.txt"
	NameDataDependency = "dataDependency.bin"
	NameDataResult     = "dataResult.bin"
	NameDataRunTime    = "dataRunTime.bin"
	NameStatistics     = "statistics.bin"
	NameUnstable       = "unstable.bin"
	NameUnstableResult = "unstable.txt"
)

type syzFuzzer struct {
	MuDependency   *sync.RWMutex
	dataDependency *DataDependency
	MuRunTime      *sync.Mutex
	dataRunTime    *DataRunTime
}

type newStats struct {
	newStat []*Statistic
}

type dependencies struct {
	newDependency []*Dependency
}

// Server is used to implement dra.DependencyServer.
type Server struct {
	address    uint32
	Port       int
	Address    string
	Dependency bool

	dataDependency *DataDependency
	dataResult     *DataResult
	dataRunTime    *DataRunTime
	stat           *Statistics

	MuFuzzer *sync.Mutex
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
	newDependency *dependencies

	// inputs of new test cases, used to get DataDependency
	newInputMu *sync.Mutex
	newInput   *Inputs

	// inputs picked by randomly, used as new test cases.
	needInputMu *sync.Mutex
	needInput   *Inputs

	// inputs of new test cases, used to check write addresses
	inputMu *sync.Mutex
	input   *Inputs

	// inputs of new test cases, used to check coverage
	coveredInputMu *sync.Mutex
	coveredInput   *Inputs

	unstableInputMu    *sync.Mutex
	unstableInputs     *UnstableInputs
	unstableInputsData *UnstableInputs
}

// GetVMOffsets is to send the offset address in vmlinux to dra
func (ss Server) GetVMOffsets(context.Context, *Empty) (*Empty, error) {
	reply := &Empty{}
	reply.Address = ss.address
	return reply, nil
}

// SendBasicBlockNumber is to get the basic block number from dra
func (ss Server) SendBasicBlockNumber(_ context.Context, request *Empty) (*Empty, error) {
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
func (ss Server) SendDependency(_ context.Context, request *Dependency) (*Empty, error) {
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
func (ss Server) SendWriteAddress(_ context.Context, _ *WriteAddresses) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendWriteAddress")

	return &Empty{}, nil
}

// Connect is to connect with syz-fuzzer
func (ss Server) Connect(_ context.Context, request *Empty) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) Connect")

	name := request.Name
	ss.MuFuzzer.Lock()
	defer ss.MuFuzzer.Unlock()

	_, ok := ss.fuzzers[name]
	if !ok {
		ss.fuzzers[name] = &syzFuzzer{
			MuDependency:   &sync.RWMutex{},
			MuRunTime:      &sync.Mutex{},
			dataDependency: &DataDependency{},
			dataRunTime: &DataRunTime{
				Tasks:      &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
				Return:     &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
				HighTask:   &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
				BootTask:   &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
				ReturnBoot: &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
			},
		}
	} else {

	}
	return &Empty{}, nil
}

func (ss Server) GetDataDependency(_ context.Context, request *Empty) (*DataDependency, error) {

	name := request.Name
	replay := &DataDependency{
		Input: map[string]*Input{},
	}
	ss.MuFuzzer.Lock()
	f, ok := ss.fuzzers[name]
	ss.MuFuzzer.Unlock()
	if ok {
		f.MuDependency.RLock()
		replay = proto.Clone(f.dataDependency).(*DataDependency)
		f.MuDependency.RUnlock()
	} else {
		for n := range ss.fuzzers {
			log.Logf(DebugLevel, "GetDataDependency name : %s", n)
		}
		log.Fatalf("GetDataDependency with error name : %s", name)
	}
	return replay, nil
}

// SendNewInput is get new input from syz-fuzzer
func (ss Server) SendNewInput(_ context.Context, request *Input) (*Empty, error) {
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
func (ss Server) GetTasks(_ context.Context, request *Empty) (*Tasks, error) {
	log.Logf(DebugLevel, "(ss Server) GetTasks")

	name := request.Name
	tasks := ss.pickTask(name)

	return tasks, nil
}

// GetBootTasks for the tasks need to be tested when boot
func (ss Server) GetBootTasks(_ context.Context, request *Empty) (*Tasks, error) {
	log.Logf(DebugLevel, "(ss Server) GetTasks")
	name := request.Name
	tasks := ss.pickBootTask(name)
	return tasks, nil
}

// ReturnTasks is to retrun the tasks from syz-fuzzer
func (ss Server) ReturnTasks(_ context.Context, request *Tasks) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) ReturnTasks")
	tasks := proto.Clone(request).(*Tasks)

	f, ok := ss.fuzzers[tasks.Name]
	if ok {
		if tasks.Kind == TaskKind_Normal || tasks.Kind == TaskKind_High || tasks.Kind == TaskKind_Ckeck {
			f.MuRunTime.Lock()
			f.dataRunTime.Return.AddTasks(tasks)
			f.MuRunTime.Unlock()
		} else if tasks.Kind == TaskKind_Boot {
			f.MuRunTime.Lock()
			f.dataRunTime.ReturnBoot.AddTasks(tasks)
			f.MuRunTime.Unlock()
		} else {
			log.Fatalf("ReturnTasks with error kind")
		}
	} else {
		log.Fatalf("ReturnTasks with error name")
	}
	reply := &Empty{}
	return reply, nil
}

// SendBootInput is get new input from syz-fuzzer
func (ss Server) SendBootInput(_ context.Context, request *Input) (*Empty, error) {
	log.Logf(DebugLevel, "(ss Server) SendBootInput")
	reply := &Empty{}
	r := proto.Clone(request).(*Input)
	ss.coveredInputMu.Lock()
	ss.coveredInput.Input = append(ss.coveredInput.Input, r)
	ss.coveredInputMu.Unlock()
	return reply, nil
}

// SendUnstableInput is get unstable input from syz-fuzzer
func (ss Server) SendUnstableInput(_ context.Context, request *UnstableInput) (*Empty, error) {
	if CollectUnstable {
		ui := proto.Clone(request).(*UnstableInput)
		ss.unstableInputMu.Lock()
		defer ss.unstableInputMu.Unlock()
		if u, ok := ss.unstableInputs.UnstableInput[ui.Sig]; ok {
			u.mergeUnstableInput(ui)
		} else {
			ss.unstableInputs.UnstableInput[ui.Sig] = ui
		}
	}

	reply := &Empty{}
	return reply, nil
}

// SendLog is to get log from syz-fuzzer
func (ss Server) SendLog(_ context.Context, request *Empty) (*Empty, error) {
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
func (ss Server) SendStat(_ context.Context, request *Statistic) (*Empty, error) {

	stat := proto.Clone(request).(*Statistic)
	ss.statMu.Lock()
	ss.newStat.newStat = append(ss.newStat.newStat, stat)
	ss.statMu.Unlock()

	reply := &Empty{}
	return reply, nil
}

// GetNeed is to random get input from syz-fuzzer, not new input but used as new input.
func (ss Server) GetNeed(context.Context, *Empty) (*Empty, error) {

	reply := &Empty{}
	if ss.needWriteaddress {
		reply.Address = 1
	} else {
		reply.Address = 0
	}
	return reply, nil
}

// SendNeedInput is to random get input from syz-fuzzer, not new input but used as new input.
func (ss Server) SendNeedInput(_ context.Context, request *Input) (*Empty, error) {
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

	ss.dataDependency = &DataDependency{
		Input:            map[string]*Input{},
		UncoveredAddress: map[uint32]*UncoveredAddress{},
		WriteAddress:     map[uint32]*WriteAddress{},
		OtherInput:       map[string]*Input{},
	}

	ss.dataResult = &DataResult{
		CoveredAddress: map[uint32]*UncoveredAddress{},
		FileOperations: map[string]*FileOperations{},
	}

	ss.dataRunTime = &DataRunTime{
		Tasks:      &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
		Return:     &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
		HighTask:   &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
		BootTask:   &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
		ReturnBoot: &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}},
	}

	ss.MuFuzzer = &sync.Mutex{}
	ss.fuzzers = map[string]*syzFuzzer{}

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
	ss.newDependency = &dependencies{newDependency: []*Dependency{}}

	ss.newInputMu = &sync.Mutex{}
	ss.newInput = &Inputs{Input: []*Input{}}

	ss.needInputMu = &sync.Mutex{}
	ss.needInput = &Inputs{Input: []*Input{}}

	ss.inputMu = &sync.Mutex{}
	ss.input = &Inputs{Input: []*Input{}}

	ss.coveredInputMu = &sync.Mutex{}
	ss.coveredInput = &Inputs{Input: []*Input{}}

	if CollectUnstable {
		ss.unstableInputMu = &sync.Mutex{}
		ss.unstableInputs = &UnstableInputs{
			UnstableInput: map[string]*UnstableInput{},
		}
		ss.unstableInputsData = &UnstableInputs{
			UnstableInput: map[string]*UnstableInput{},
		}
	}

	lis, err := net.Listen("tcp", ss.Address)
	log.Logf(DebugLevel, "drpc on tcp : %s", ss.Address)
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

	if CollectPath {

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

	if ss.needWriteaddress {
		ss.needWriteaddress = false
	}
	// reboot the qemu
	if ss.needboot {
		ss.needboot = false
	}

	// deal need input
	ss.needInputMu.Lock()
	//needInput := append([]*Input{}, ss.needInput.Input...)
	ss.needInput = &Inputs{Input: []*Input{}}
	ss.needInputMu.Unlock()
	//for _, i := range needInput {
	//	ss.addInput(i)
	//}
	//needInput = nil

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
	var returnTask []*Task
	for _, f := range ss.fuzzers {
		f.MuRunTime.Lock()
		for _, t := range f.dataRunTime.Return.TaskArray {
			returnTask = append(returnTask, t)
		}
		f.dataRunTime.Return.emptyTask()
		f.MuRunTime.Unlock()
	}
	for _, task := range returnTask {
		if t, ok := ss.dataRunTime.Tasks.TaskMap[task.Hash]; ok {
			t.mergeTask(task)
			for u := range t.UncoveredAddress {
				_, ok := ss.dataDependency.UncoveredAddress[u]
				if ok {

				} else {
					delete(t.UncoveredAddress, u)
				}
			}
		} else {
			ss.dataRunTime.Tasks.AddTask(task)
		}
		ss.updateUncoveredAddress(task)
	}
	sort.Slice(ss.dataRunTime.Tasks.TaskArray, func(i, j int) bool {
		return ss.dataRunTime.Tasks.TaskArray[i].getRealPriority() < ss.dataRunTime.Tasks.TaskArray[j].getRealPriority()
	})
	returnTask = nil

	// get new tasks
	if ss.Dependency {
		t := time.Now()
		elapsed := t.Sub(ss.timeStart)
		if elapsed.Seconds() > startTime {
			if len(ss.dataRunTime.HighTask.TaskArray) != 0 {
				var task []*Task
				for _, t := range ss.dataRunTime.HighTask.TaskArray {
					for u := range t.UncoveredAddress {
						_, ok := ss.dataDependency.UncoveredAddress[u]
						if ok {

						} else {
							delete(t.UncoveredAddress, u)
						}
					}
					if len(t.UncoveredAddress) > 0 {
						task = append(task, t)
					}
				}
				ss.dataRunTime.HighTask.emptyTask()
				task = []*Task{}
				for _, t := range task {
					t.TaskRunTimeData = []*TaskRunTimeData{}
				}
				for _, f := range ss.fuzzers {
					f.MuRunTime.Lock()
					for _, t := range task {
						f.dataRunTime.HighTask.AddTask(proto.Clone(t).(*Task))
					}
					f.MuRunTime.Unlock()
				}
				task = nil
			}
			{
				var task []*Task
				for _, t := range ss.dataRunTime.Tasks.TaskArray {
					for u := range t.UncoveredAddress {
						_, ok := ss.dataDependency.UncoveredAddress[u]
						if ok {

						} else {
							delete(t.UncoveredAddress, u)
						}
					}
					if len(t.UncoveredAddress) > 0 && t.Count < TaskCountLimitation {
						//if len(t.UncoveredAddress) > 0 {
						if t.TaskStatus == TaskStatus_untested {
							t.TaskStatus = TaskStatus_testing
							t.reducePriority()
							task = append(task, proto.Clone(t).(*Task))
						} else if t.TaskStatus < TaskStatus_tested {
							t.reducePriority()
							task = append(task, t)
						}
						if len(task) > TaskNum {
							break
						}
					}
				}
				for _, t := range task {
					t.TaskRunTimeData = []*TaskRunTimeData{}
				}
				for _, f := range ss.fuzzers {
					f.MuRunTime.Lock()
					f.dataRunTime.Tasks.emptyTask()
					for _, t := range task {
						f.dataRunTime.Tasks.AddTask(proto.Clone(t).(*Task))
					}
					f.MuRunTime.Unlock()
				}
				task = nil
			}
		}
	} else {
		ss.dataRunTime.HighTask.emptyTask()
		ss.dataRunTime.Tasks.emptyTask()
	}

	// deal return boot tasks
	returnBootTask := &Tasks{Name: "", TaskMap: map[string]*Task{}, TaskArray: []*Task{}}
	for _, f := range ss.fuzzers {
		f.MuRunTime.Lock()
		returnBootTask.AddTasks(f.dataRunTime.ReturnBoot)
		f.dataRunTime.ReturnBoot.emptyTask()
		f.MuRunTime.Unlock()
	}
	for hash, task := range returnBootTask.TaskMap {
		if t, ok := ss.dataRunTime.BootTask.TaskMap[hash]; ok {
			if task.TaskStatus == TaskStatus_covered {
				t.mergeTask(task)
			} else {
				t.TaskStatus = TaskStatus_tested
			}
			t.mergeTask(task)
			for u := range t.UncoveredAddress {
				_, ok := ss.dataDependency.UncoveredAddress[u]
				if ok {

				} else {
					delete(t.UncoveredAddress, u)
				}
			}
		}
	}
	sort.Slice(ss.dataRunTime.BootTask.TaskArray, func(i, j int) bool {
		return ss.dataRunTime.BootTask.TaskArray[i].getRealPriority() < ss.dataRunTime.BootTask.TaskArray[j].getRealPriority()
	})
	returnBootTask = nil

	// get boot tasks
	if ss.Dependency {
		t := time.Now()
		elapsed := t.Sub(ss.timeStart)
		if elapsed.Seconds() > startTime {
			if len(ss.dataRunTime.BootTask.TaskArray) != 0 {
				var task []*Task
				for _, t := range ss.dataRunTime.BootTask.TaskArray {
					if t.TaskStatus == TaskStatus_untested {
						for u := range t.UncoveredAddress {
							_, ok := ss.dataDependency.UncoveredAddress[u]
							if ok {

							} else {
								delete(t.UncoveredAddress, u)
							}
						}
						if len(t.UncoveredAddress) > 0 {
							task = append(task, t)
						}
					}
				}
				for _, t := range task {
					t.TaskRunTimeData = []*TaskRunTimeData{}
				}
				for _, f := range ss.fuzzers {
					f.MuRunTime.Lock()
					for _, t := range task {
						f.dataRunTime.BootTask.AddTask(proto.Clone(t).(*Task))
					}
					f.MuRunTime.Unlock()
				}
				task = nil
			}
		}
	} else {
		ss.dataRunTime.BootTask.emptyTask()
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

	ss.writeMessageToDisk(ss.dataDependency, NameDataDependency)
	ss.writeMessageToDisk(ss.dataResult, NameDataResult)
	ss.writeMessageToDisk(ss.dataRunTime, NameDataRunTime)
	ss.writeMessageToDisk(ss.stat, NameStatistics)

	if CollectUnstable {
		ss.unstableInputMu.Lock()
		unstableInput := map[string]*UnstableInput{}
		for sig, ui := range ss.unstableInputs.UnstableInput {
			unstableInput[sig] = ui
		}
		ss.unstableInputs = &UnstableInputs{
			UnstableInput: map[string]*UnstableInput{},
		}
		ss.unstableInputMu.Unlock()
		for sig, ui := range unstableInput {
			if i, ok := ss.unstableInputsData.UnstableInput[sig]; ok {
				i.mergeUnstableInput(ui)
				ss.outPutUnstableInput(i)
			} else {
				ss.unstableInputsData.UnstableInput[sig] = ui
				ss.outPutUnstableInput(ui)
			}
		}
		ss.writeMessageToDisk(ss.unstableInputsData, NameUnstable)
	}

	if CheckCondition {
		d := proto.Clone(ss.dataDependency).(*DataDependency)
		for _, i := range d.Input {
			i.Paths = nil
			i.Call = nil
		}
		d.OtherInput = nil
		for _, f := range ss.fuzzers {
			f.MuDependency.Lock()
			f.dataDependency = d
			f.MuDependency.Unlock()
		}
	}

	t := time.Now()
	elapsed := t.Sub(ss.timeStart)
	if elapsed.Seconds() > TimeExit {
		os.Exit(1)
	}
}
