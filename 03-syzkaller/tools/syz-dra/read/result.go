package main

import (
	"fmt"
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type result struct {
	path           string
	dirName        string
	baseName       string
	dataDependency *pb.DataDependency
	dataResult     *pb.DataResult
	dataRunTime    *pb.DataRunTime
	statistics     *pb.Statistics

	uncoveredAddressInput      map[uint32]*pb.UncoveredAddress
	uncoveredAddressDependency map[uint32]*pb.UncoveredAddress
	coveredAddressInput        map[uint32]*pb.UncoveredAddress
	coveredAddressDependency   map[uint32]*pb.UncoveredAddress
}

func (r *result) read(path string) {
	r.path = path
	r.dirName = filepath.Dir(path)
	r.baseName = filepath.Base(path)

	fileName := filepath.Join(r.path, pb.NameDataDependency)
	in, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	r.dataDependency = &pb.DataDependency{}
	if err := proto.Unmarshal(in, r.dataDependency); err != nil {
		log.Fatalln("Failed to parse data:", err)
	}

	fileName = filepath.Join(r.path, pb.NameDataResult)
	in, err = ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	r.dataResult = &pb.DataResult{}
	if err := proto.Unmarshal(in, r.dataResult); err != nil {
		log.Fatalln("Failed to parse data:", err)
	}

	fileName = filepath.Join(r.path, pb.NameDataRunTime)
	in, err = ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	r.dataRunTime = &pb.DataRunTime{}
	if err := proto.Unmarshal(in, r.dataRunTime); err != nil {
		log.Fatalln("Failed to parse data:", err)
	}

	fileName = filepath.Join(r.path, pb.NameStatistics)
	in, err = ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	r.statistics = &pb.Statistics{}
	if err := proto.Unmarshal(in, r.statistics); err != nil {
		log.Fatalln("Failed to parse statistics:", err)
	}

	r.getUncoveredAddress()
}

func (r *result) getUncoveredAddress() {
	r.uncoveredAddressInput = make(map[uint32]*pb.UncoveredAddress)
	r.uncoveredAddressDependency = make(map[uint32]*pb.UncoveredAddress)
	for _, ua := range r.dataDependency.UncoveredAddress {
		if ua.Kind == pb.UncoveredAddressKind_InputRelated {
			r.uncoveredAddressInput[ua.UncoveredAddress] = ua
		} else if ua.Kind == pb.UncoveredAddressKind_DependencyRelated {
			r.uncoveredAddressDependency[ua.UncoveredAddress] = ua
		} else {

		}
	}
	for _, a := range r.dataResult.CoveredAddress {
		if a.Kind == pb.UncoveredAddressKind_InputRelated {
			r.coveredAddressInput[a.UncoveredAddress] = a
		} else if a.Kind == pb.UncoveredAddressKind_DependencyRelated {
			r.coveredAddressDependency[a.UncoveredAddress] = a
		} else {

		}
	}
}

func (r *result) checkTasks() {

}

func (r *result) checkUncoveredAddress(uncoveredAddress uint32) string {

	ua, ok := r.dataDependency.UncoveredAddress[uncoveredAddress]
	if ok {

	} else {
		return ""
	}
	ua.RunTimeDate.TaskStatus = pb.TaskStatus_not_find_input

	res := ""
	res += "*******************************************\n"
	res += "condition address 						: " + fmt.Sprintf("0xffffffff%x", ua.ConditionAddress-5) + "\n"
	res += "uncovered address 					 	: " + fmt.Sprintf("0xffffffff%x", ua.UncoveredAddress-5) + "\n"
	res += "number_arrive_basic_blocks 			 	: " + fmt.Sprintf("%d", ua.NumberArriveBasicblocks) + "\n"
	res += "number_dominator_instructions(using) 	: " + fmt.Sprintf("%d", ua.NumberDominatorInstructions) + "\n"
	res += "*******************************************\n"

	res += "*******************************************\n"
	res += "# input : " + fmt.Sprintf("%d", len(ua.Input)) + "\n"
	for sig, indexBits := range ua.Input {
		res += "-------------------------------------------\n"
		if input, ok := r.dataDependency.Input[sig]; ok {
			res += "sig : " + input.Sig + "\n"
			res += "index : " + fmt.Sprintf("%b", indexBits) + "\n"
			res += fmt.Sprintf("%s", input.Program) + "\n"
		} else {
			res += "not find input : " + input.Sig + "\n"
			ua.RunTimeDate.TaskStatus = pb.TaskStatus_not_find_input
		}
	}
	res += "*******************************************\n"

	res += "*******************************************\n"
	ua.WriteAddressStatus = map[uint32]pb.TaskStatus{}
	res += "# write : " + fmt.Sprintf("%d", len(ua.WriteAddress)) + "\n"
	if len(ua.WriteAddress) == 0 {
		res += "not find write address of " + fmt.Sprintf("0xffffffff%x", ua.UncoveredAddress-5) + "\n"
	} else {
		for wa, waa := range ua.WriteAddress {
			ua.WriteAddressStatus[wa] = pb.TaskStatus_not_find_write_address
			res += "-------------------------------------------\n"
			res += "## write address : " + fmt.Sprintf("0xffffffff%x", wa-5) + "\n"
			res += "Repeat 		: " + fmt.Sprintf("%d", waa.Repeat) + "\n"
			res += "Priority 	: " + fmt.Sprintf("%d", waa.Prio) + "\n"
			if waaa, ok := r.dataDependency.WriteAddress[wa]; ok {
				if len(waaa.Input) == 0 {
					res += "not find write input : " + fmt.Sprintf("0xffffffff%x", wa-5) + "\n"
					if ua.RunTimeDate.TaskStatus < pb.TaskStatus_not_find_write_input {
						ua.RunTimeDate.TaskStatus = pb.TaskStatus_not_find_write_input
					}
				} else {
					for sig, indexBits := range waaa.Input {
						res += "-------------------------------------------\n"
						if input, ok := r.dataDependency.Input[sig]; ok {
							res += "sig : " + input.Sig + "\n"
							res += "index : " + fmt.Sprintf("%b", indexBits) + "\n"
							res += fmt.Sprintf("%s", input.Program) + "\n"
							if ua.RunTimeDate.TaskStatus < pb.TaskStatus_untested {
								ua.RunTimeDate.TaskStatus = pb.TaskStatus_untested
							}
						} else {
							res += "not find input : " + input.Sig + "\n"
							if ua.RunTimeDate.TaskStatus < pb.TaskStatus_not_find_write_input {
								ua.RunTimeDate.TaskStatus = pb.TaskStatus_not_find_write_input
							}
						}
					}

				}
			} else {
				res += "not find write address : " + fmt.Sprintf("0xffffffff%x", wa-5) + "\n"
				if ua.RunTimeDate.TaskStatus < pb.TaskStatus_not_find_write_address {
					ua.RunTimeDate.TaskStatus = pb.TaskStatus_not_find_write_address
				}
			}
		}
	}
	res += "*******************************************\n"

	res += "*******************************************\n"
	ua.TasksCount = map[int32]uint32{}
	ua.WriteCount = map[int32]uint32{}
	tasks := &pb.Tasks{
		Name:      "",
		Kind:      0,
		TaskMap:   map[string]*pb.Task{},
		TaskArray: []*pb.Task{},
	}
	for _, t := range r.dataRunTime.Tasks.TaskArray {
		if _, ok := t.UncoveredAddress[uncoveredAddress]; ok {
			tasks.AddTask(t)
		} else if _, ok := t.CoveredAddress[uncoveredAddress]; ok {
			tasks.AddTask(t)
		}
	}
	for _, t := range tasks.TaskArray {
		res += "*******************************************\n"
		res += "task_hash 		: " + t.ComputeHash() + "\n"
		res += "task_status 	: " + t.TaskStatus.String() + "\n"
		res += "task priority 	: " + fmt.Sprintf("%d", t.Priority) + "\n"
		priority := uint32(0)
		for _, ua := range t.UncoveredAddress {
			priority += ua.Priority
		}
		res += "uncovered address priority : " + fmt.Sprintf("%d", priority) + "\n"
		res += "condition program : " + fmt.Sprintf("%d", t.Index) + " : " + t.Sig + "\n"
		res += fmt.Sprintf("%s", t.Program) + "\n"
		res += "write program : " + fmt.Sprintf("%d", t.WriteIndex) + " : " + t.WriteSig + "\n"
		res += fmt.Sprintf("%s", t.WriteProgram) + "\n"

		ua.TasksCount[int32(t.TaskStatus)]++
		ua.RunTimeDate.RecursiveCount += t.Count
		res += "check 						: " + fmt.Sprintf("%t", t.Check) + "\n"
		res += "len TaskRunTimeData	 		: " + fmt.Sprintf("%d", len(t.TaskRunTimeData)) + "\n"

		count := 0
		for _, rTD := range t.UncoveredAddress {
			if rTD.TaskStatus > pb.TaskStatus_untested {
				count++
			}
		}
		res += "uncovered address 			: " + fmt.Sprintf("%d", len(t.UncoveredAddress)) + "\n"
		res += "tested uncovered address 	: " + fmt.Sprintf("%d", count) + "\n"

		res += "-------------------------------------------\n"
		if rTD, ok := t.UncoveredAddress[uncoveredAddress]; ok {
			res += "task_status 		: " + rTD.TaskStatus.String() + "\n"
			res += "write address		: " + fmt.Sprintf("0xffffffff%x", rTD.WriteAddress-5) + "\n"
			res += "condition address 	: " + fmt.Sprintf("0xffffffff%x", rTD.ConditionAddress-5) + "\n"
			res += "uncovered address 	: " + fmt.Sprintf("0xffffffff%x", rTD.Address-5) + "\n"
			res += "check write 		: " + fmt.Sprintf("%t", rTD.CheckWrite) + "\n"
			res += "check condition 	: " + fmt.Sprintf("%t", rTD.CheckCondition) + "\n"
			res += "check address 		: " + fmt.Sprintf("%t", rTD.CheckAddress) + "\n"

			if rTD.TaskStatus == pb.TaskStatus_untested {
				continue
			}

			if len(rTD.RightBranchAddress) == len(rTD.CheckRightBranchAddress) {
				for i, b := range rTD.CheckRightBranchAddress {
					res += "check branch address: " + fmt.Sprintf("0xffffffff%x", rTD.RightBranchAddress[i]-5) + "\n"
					res += "check branch 		: " + fmt.Sprintf("%t", b) + "\n"
				}
			} else {
				res += "len(rTD.RightBranchAddress) != len(rTD.CheckRightBranchAddress)\n"
			}

			if rTD.CheckWrite {
				if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_stable_write {
					ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_stable_write
				}
			} else {
				res += "unstable write address" + "\n"
				if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_unstable_write {
					ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_unstable_write
				}
			}

			if rTD.CheckCondition {
				if rTD.CheckAddress {
					res += "error in rdd.CheckCondition" + "\n"
					if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_covered {
						ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_covered
					}
				} else {
					if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_stable_condition && ua.WriteAddressStatus[rTD.WriteAddress] >= pb.TaskStatus_stable_write {
						ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_stable_condition
					}
				}
			} else {
				res += "unstable condition address" + "\n"
				if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_unstable_condition && ua.WriteAddressStatus[rTD.WriteAddress] >= pb.TaskStatus_stable_write {
					ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_unstable_condition
				}
			}

			if t.Check {
				for _, trd := range t.TaskRunTimeData {
					if rdd, ok := trd.UncoveredAddress[uncoveredAddress]; ok {
						res += "-------------------------------------------\n"
						res += "insert task_status 			: " + rdd.TaskStatus.String() + "\n"
						res += "check insert write address 	: " + fmt.Sprintf("%t", rdd.CheckWrite) + "\n"
						res += "check condition 			: " + fmt.Sprintf("%t", rdd.CheckCondition) + "\n"
						res += "check address 				: " + fmt.Sprintf("%t", rdd.CheckAddress) + "\n"
						if rdd.TaskStatus == pb.TaskStatus_untested {
							continue
						}

						if len(rdd.RightBranchAddress) == len(rdd.CheckRightBranchAddress) {
							for i, b := range rdd.CheckRightBranchAddress {
								res += "check branch address: " + fmt.Sprintf("0xffffffff%x", rdd.RightBranchAddress[i]-5) + "\n"
								res += "check branch 		: " + fmt.Sprintf("%t", b) + "\n"
							}
						} else {
							res += "len(rdd.RightBranchAddress) != len(rdd.CheckRightBranchAddress)\n"
						}

						if rdd.CheckWrite {
							if rdd.CheckCondition {
								if rdd.CheckAddress {
									res += "error in rdd.CheckCondition" + "\n"
									if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_covered {
										ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_covered
									}
								} else {
									res += "useless write address or FP" + "\n"
									if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_tested {
										ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_tested
									}
								}
							} else {
								res += "unstable insert condition address" + "\n"
								if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_unstable_insert_condition && ua.WriteAddressStatus[rTD.WriteAddress] >= pb.TaskStatus_stable_insert_write {
									ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_unstable_insert_condition
								}
							}
						} else {
							res += "unstable insert write address" + "\n"
							if ua.WriteAddressStatus[rTD.WriteAddress] < pb.TaskStatus_unstable_insert_write && ua.WriteAddressStatus[rTD.WriteAddress] >= pb.TaskStatus_stable_condition {
								ua.WriteAddressStatus[rTD.WriteAddress] = pb.TaskStatus_unstable_insert_write
							}
						}

					} else if _, ok := trd.CoveredAddress[uncoveredAddress]; ok {
						res += "uncoveredAddress in trd.CoveredAddress" + "\n"
					} else {
						res += "no test" + "\n"
					}
				}

			} else {

			}

		} else if _, ok := t.CoveredAddress[uncoveredAddress]; ok {
			res += "uncoveredAddress in t.covered_address" + "\n"
			if ua.RunTimeDate.TaskStatus < pb.TaskStatus_covered {
				ua.RunTimeDate.TaskStatus = pb.TaskStatus_covered
			}
		} else {
		}
	}

	res += "-------------------------------------------\n"
	for _, ts := range ua.WriteAddressStatus {
		ua.WriteCount[int32(ts)]++
	}
	res += "tasksCount" + "\n"
	for ts, c := range ua.TasksCount {
		res += pb.TaskStatus_name[ts] + " : " + fmt.Sprintf("%d", c) + "\n"
	}
	res += "writeCount" + "\n"
	for ts, c := range ua.WriteCount {
		res += pb.TaskStatus_name[ts] + " : " + fmt.Sprintf("%d", c) + "\n"
		if c > 0 && ua.RunTimeDate.TaskStatus < pb.TaskStatus(ts) {
			ua.RunTimeDate.TaskStatus = pb.TaskStatus(ts)
		}
	}
	res += "ua.RunTimeDate.TaskStatus : " + ua.RunTimeDate.TaskStatus.String() + "\n"

	res += "*******************************************\n"
	return res
}

func (r *result) checkStatistic() {

	res := ""
	res += fmt.Sprintf("singal number 		: %d\n", r.statistics.SignalNum)
	res += fmt.Sprintf("basic block number	: %d\n", r.statistics.BasicBlockNumber)
	res += fmt.Sprintf("coverage			: %d\n", len(r.statistics.Coverage.Coverage))

	for _, s := range r.statistics.Stat {
		res += "-------------------------------------------\n"
		res += fmt.Sprintf("stat name				: %s\n", s.Name.String())
		res += fmt.Sprintf("execute number			: %d\n", s.ExecuteNum)
		res += fmt.Sprintf("time					: %f\n", s.Time)
		res += fmt.Sprintf("new test case number 	: %d\n", s.NewTestCaseNum)
		res += fmt.Sprintf("new address number 		: %d\n", s.NewAddressNum)
		res += "-------------------------------------------\n"
	}

	f, _ := os.OpenFile(filepath.Join(r.path, pb.NameData), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	_, _ = f.WriteString(res)
	_ = f.Close()
}
