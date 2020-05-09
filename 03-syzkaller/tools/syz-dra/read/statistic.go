package main

import (
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
)

type statistic struct {
	Kind string
	Name string
	tag  []string
	data []uint32
}

func (s *statistic) output(dir string) {
	path := filepath.Join(dir, s.Name, s.Name + ".txt")
	fmt.Printf("statistic path : %s\n", path)
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	_, _ = f.WriteString(fmt.Sprintf("%s", s.Name))
	for _, v := range s.data {
		_, _ = f.WriteString(fmt.Sprintf("##%d", v))
	}
	_, _ = f.WriteString(fmt.Sprintf("\n"))
	_ = f.Close()
}

func average(ss []*statistic) *statistic {

	res := &statistic{
		Kind: "",
		Name: "",
		tag:  []string{},
		data: []uint32{},
	}
	if len(ss) > 0 {
		res.Kind = ss[0].Kind
		res.Name = ss[0].Name
		for _, t := range ss[0].tag {
			res.tag = append(res.tag, t)
		}
		res.data = make([]uint32, len(res.tag))
	} else {
		return nil
	}
	for _, s := range ss {
		if s.Kind != res.Kind {
			return nil
		} else {
			for i, d := range s.data {
				res.data[i] += d
			}
		}
	}
	for i := range res.data {
		res.data[i] /= uint32(len(ss))
	}
	return res
}

func prevalent(r *result) *statistic {
	res := &statistic{
		Kind: "prevalent",
		Name: "",
		tag: []string{
			"NumberBasicBlockReal",
			"NumberCovered",
			"NumberUncovered",
			"NumberUnresolvedConditions",
			"NumberNotDependency",
			"NumberDependency",
			"NumberInstructions",
			"NumberInstructionsDominator",
		},
		data: nil,
	}
	res.data = make([]uint32, len(res.tag))
	index := 0

	index = 0
	res.data[index+0] = r.statistics.NumberBasicBlockReal
	fmt.Printf("r.statistics.NumberBasicBlockReal : %d\n", r.statistics.NumberBasicBlockReal)
	res.data[index+1] = r.statistics.NumberBasicBlockCovered
	res.data[index+2] = res.data[index+0] - res.data[index+1]
	res.data[index+3] = uint32(len(r.dataDependency.UncoveredAddress))

	index = 4
	res.data[index+0] = 0
	res.data[index+1] = 0
	res.data[index+2] = 0
	res.data[index+3] = 0
	for _, u := range r.dataDependency.UncoveredAddress {
		if u.Kind == pb.UncoveredAddressKind_UncoveredAddressInputRelated {
			res.data[index+0]++
		} else if u.Kind == pb.UncoveredAddressKind_UncoveredAddressDependencyRelated {
			res.data[index+1]++
			res.data[index+2] += u.NumberArriveBasicblocks
			res.data[index+3] += u.NumberDominatorInstructions
		}
	}
	res.data[index+2] /= res.data[index+1]
	res.data[index+3] /= res.data[index+1]
	return res
}

func write_statement(r *result) *statistic {
	res := &statistic{
		Kind: "write_statement",
		Name: "",
		tag: []string{
			"NumberConditions",
			"NumberWriteStatement",
			"NumberConstant",
			"NumberExpression",
			"NumberUseful",
		},
		data: nil,
	}
	res.data = make([]uint32, len(res.tag))
	index := 0

	index = 0
	res.data[index+0] = 0
	res.data[index+1] = 0
	res.data[index+2] = 0
	res.data[index+3] = 0
	for _, ua := range r.dataDependency.UncoveredAddress {
		if len(ua.WriteAddress) > 0 {
			res.data[index+0] += 1
			res.data[index+1] += uint32(len(ua.WriteAddress))
			for wa := range ua.WriteAddress {
				if ws, ok := r.dataDependency.WriteAddress[wa]; ok {
					if ws.Kind == pb.WriteStatementKind_WriteStatementConstant {
						res.data[index+2]++
					} else if ws.Kind == pb.WriteStatementKind_WriteStatementNonconstant {
						res.data[index+3]++
					} else {

					}
				}
			}
		}
	}

	res.data[index+1] /= res.data[index+0]
	res.data[index+2] /= res.data[index+0]
	res.data[index+3] /= res.data[index+0]
	return res
}

func unstable(r *result) *statistic {
	res := &statistic{
		Kind: "unstable",
		Name: "",
		tag: []string{
			"NumberInput",
			"NumberCondition",
			"NumberDependency",
			"NumberTaskCondition",
			"NumberStable",
			"NumberUnstable",
			"NumberTaskWrite",
			"NumberStable",
			"NumberUnstable",
			"NumberCombination",
			"NumberStable",
			"NumberUnstable",
		},
		data: nil,
	}
	res.data = make([]uint32, len(res.tag))
	index := 0

	index = 0
	res.data[index+0] = 0
	res.data[index+1] = 0
	res.data[index+2] = 0
	for _, i := range r.dataDependency.Input {
		if len(i.UncoveredAddress) > 0 {
			res.data[index+0] += 1
			res.data[index+1] += uint32(len(i.UncoveredAddress))
			// fmt.Printf("len(i.UncoveredAddress) : %d\n", len(i.UncoveredAddress))
			for address := range i.UncoveredAddress {
				if ua, ok := r.dataDependency.UncoveredAddress[address]; ok {
					if ua.Kind == pb.UncoveredAddressKind_UncoveredAddressDependencyRelated {
						res.data[index+2] += 1
					}
				}
			}
		}
	}
	res.data[index+1] /= res.data[index+0]
	res.data[index+2] /= res.data[index+0]

	index = 3
	for _, t := range r.dataRunTime.Tasks.TaskArray {
		for _, ua := range t.UncoveredAddress {
			if ua.CheckCondition {
				res.data[index+1]++
			} else {
				res.data[index+2]++
			}
		}
	}
	res.data[index+0] = res.data[index+1] + res.data[index+2]

	index = 6
	for _, t := range r.dataRunTime.Tasks.TaskArray {
		for _, ua := range t.UncoveredAddress {
			if ua.CheckWrite {
				res.data[index+1]++
			} else {
				res.data[index+2]++
			}
		}
	}
	res.data[index+0] = res.data[index+1] + res.data[index+2]

	index = 9
	for _, t := range r.dataRunTime.Tasks.TaskArray {
		for _, tr := range t.TaskRunTimeData {
			for _, t := range tr.UncoveredAddress {
				if t.CheckWrite && t.CheckCondition {
					res.data[index+1]++
				} else {
					res.data[index+2]++
				}
			}
		}
	}
	res.data[index+0] = res.data[index+1] + res.data[index+2]

	return res
}
