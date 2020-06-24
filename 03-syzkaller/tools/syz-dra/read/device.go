package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
)

type device struct {
	path     string
	dirName  string
	baseName string
	dataPath string
	a2i      bool

	base              *result
	resultsWithDra    *results
	resultsWithoutDra *results

	uniqueCoverageWithDra    map[uint32]uint32
	uniqueCoverageWithoutDra map[uint32]uint32
	unionCoverage            map[uint32]uint32
	intersectionCoverage     map[uint32]uint32
}

func (d *device) read(path string, a2i bool) {
	d.path = path
	d.dirName = filepath.Dir(path)
	d.baseName = filepath.Base(path)
	d.dataPath = filepath.Join(path, pb.NameData)
	d.a2i = a2i

	pathBase := filepath.Join(d.path, pb.NameBase)
	if _, err := os.Stat(pathBase); os.IsNotExist(err) {
		fmt.Printf(pb.NameBase + " does not exist\n")
	} else {
		d.base = &result{}
		d.base.read(pathBase)
	}

	d.resultsWithDra = &results{}
	d.resultsWithDra.read(filepath.Join(d.path, pb.NameWithDra))
	d.resultsWithoutDra = &results{}
	d.resultsWithoutDra.read(filepath.Join(d.path, pb.NameWithoutDra))

	_ = os.Remove(filepath.Join(d.path, pb.NameData))
	fmt.Printf("remove %s\n", pb.NameData)

	d.checkStatistic()
	d.checkCoverage()
	d.checkUncoveredAddress()

}

func (d *device) checkCoverage() {
	d.uniqueCoverageWithDra = map[uint32]uint32{}
	d.uniqueCoverageWithoutDra = map[uint32]uint32{}
	d.unionCoverage = map[uint32]uint32{}
	d.intersectionCoverage = map[uint32]uint32{}

	for a, c := range d.resultsWithDra.maxCoverage {
		if cc, ok := d.resultsWithoutDra.maxCoverage[a]; ok {
			d.intersectionCoverage[a] = c + cc
		} else {
			d.uniqueCoverageWithDra[a] = c
		}
		d.unionCoverage[a] += c
	}
	for a, c := range d.resultsWithoutDra.maxCoverage {
		if _, ok := d.resultsWithDra.maxCoverage[a]; ok {

		} else {
			d.uniqueCoverageWithoutDra[a] = c
		}
		d.unionCoverage[a] += c
	}

	res := ""
	res += "*******************************************\n"
	res += "coverage : " + "\n"
	res += "uniqueCoverageWithDra    : " + fmt.Sprintf("%5d", len(d.uniqueCoverageWithDra)) + "\n"
	res += "uniqueCoverageWithoutDra : " + fmt.Sprintf("%5d", len(d.uniqueCoverageWithoutDra)) + "\n"
	res += "unionCoverage            : " + fmt.Sprintf("%5d", len(d.unionCoverage)) + "\n"
	res += "intersectionCoverage     : " + fmt.Sprintf("%5d", len(d.intersectionCoverage)) + "\n"
	res += "*******************************************\n"

	solvedCondition := map[uint32]*pb.RunTimeData{}
	for _, r := range d.resultsWithDra.result {
		for _, t := range r.dataRunTime.Tasks.TaskArray {
			for ca, rt := range t.CoveredAddress {
				solvedCondition[ca] = rt
			}
		}
	}
	stableSolvedCondition := map[uint32]*pb.RunTimeData{}
	unStableSolvedCondition := map[uint32]*pb.RunTimeData{}
	for a, rt := range solvedCondition {
		if _, ok := d.resultsWithDra.maxCoverage[a]; ok {
			stableSolvedCondition[a] = rt
		} else {
			unStableSolvedCondition[a] = rt
		}
	}
	res += "*******************************************\n"
	res += "solvedCondition         : " + fmt.Sprintf("%5d", len(solvedCondition)) + "\n"
	res += "stableSolvedCondition   : " + fmt.Sprintf("%5d", len(stableSolvedCondition)) + "\n"
	res += "unStableSolvedCondition : " + fmt.Sprintf("%5d", len(unStableSolvedCondition)) + "\n"
	res += "*******************************************\n"

	f, _ := os.OpenFile(d.dataPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	_, _ = f.WriteString(res)
	_ = f.Close()

}

func (d *device) checkUncoveredAddress() {

	res := ""

	UA := map[uint32]uint32{}
	UAD := map[uint32]uint32{}
	UADCW := map[uint32]uint32{}
	UADCWD := map[uint32]uint32{}
	UADCWDU := map[uint32]uint32{}
	UADCWO := map[uint32]uint32{}
	if d.base != nil {
		for a := range d.base.uncoveredAddressDependency {
			UA[a] = 0
			UAD[a] = 0
			if c, ok := d.resultsWithDra.maxCoverage[a]; ok {
				UADCW[a] = c
				if c > 0 {
					UADCWD[a] = c
					if _, ok := d.uniqueCoverageWithDra[a]; ok {
						UADCWDU[a] = c
					}
				}
			}
			if c, ok := d.resultsWithoutDra.maxCoverage[a]; ok {
				UADCWO[a] = c
			}
		}
	} else {
		for _, r := range d.resultsWithDra.result {
			for a := range r.uncoveredAddressDependency {
				UA[a] = 0
				UAD[a] = 0
				if c, ok := d.resultsWithDra.maxCoverage[a]; ok {
					UADCW[a] = c
					if c > 0 {
						res += "uncovered address covered by dependency : " + fmt.Sprintf("0xffffffff%x", a-5) + "\n"
						UADCWD[a] = c
						if _, ok := d.uniqueCoverageWithDra[a]; ok {
							UADCWDU[a] = c
						}
					}
				}
				if c, ok := d.resultsWithoutDra.maxCoverage[a]; ok {
					UADCWO[a] = c
				}
			}

			for a := range r.dataResult.CoveredAddress {
				UA[a] = 0
				UAD[a] = 0
				if c, ok := d.resultsWithDra.maxCoverage[a]; ok {
					UADCW[a] = c
					if c > 0 {
						res += "covered address covered by dependency : " + fmt.Sprintf("0xffffffff%x", a-5) + "\n"
						UADCWD[a] = c
						if _, ok := d.uniqueCoverageWithDra[a]; ok {
							UADCWDU[a] = c
						}
					}
				}
				if c, ok := d.resultsWithoutDra.maxCoverage[a]; ok {
					UADCWO[a] = c
				}
			}
		}
	}
	res += "*******************************************\n"
	res += "number of uncovered address      : " + fmt.Sprintf("%5d", len(UA)) + "\n"
	res += "related to dependency            : " + fmt.Sprintf("%5d", len(UAD)) + "\n"
	res += "covered by syzkaller with dra    : " + fmt.Sprintf("%5d", len(UADCW)) + "\n"
	res += "covered by dependency mutate     : " + fmt.Sprintf("%5d", len(UADCWD)) + "\n"
	res += "unique one of them               : " + fmt.Sprintf("%5d", len(UADCWDU)) + "\n"
	res += "covered by syzkaller without dra : " + fmt.Sprintf("%5d", len(UADCWO)) + "\n"
	res += "*******************************************\n"

	for _, r := range d.resultsWithoutDra.result {

		r.checkStatistic()

	}

	resultSize := uint32(len(d.resultsWithDra.result))

	writeAddressCount := map[uint32]uint32{}
	tempWA := map[uint32]*pb.WriteAddress{}
	for _, r := range d.resultsWithDra.result {

		for address, writeAddress := range r.dataDependency.WriteAddress {
			if len(writeAddress.Input) == 0 {
				tempWA[address] = writeAddress
				if c, ok := writeAddressCount[address]; ok {
					writeAddressCount[address] = c + 1
				} else {
					writeAddressCount[address] = 1
				}
			}
		}
	}

	allWriteAddress := map[uint32]*pb.WriteAddress{}
	for address, count := range writeAddressCount {
		if count == resultSize {
			allWriteAddress[address] = tempWA[address]
		}
	}

	uncoveringAddressCount := map[uint32]uint32{}
	tempUA := map[uint32]*pb.UncoveredAddress{}
	for _, r := range d.resultsWithDra.result {

		r.checkStatistic()

		for address, uncoveringAddress := range r.dataDependency.UncoveredAddress {
			tempUA[address] = uncoveringAddress
			if c, ok := uncoveringAddressCount[address]; ok {
				uncoveringAddressCount[address] = c + 1
			} else {
				uncoveringAddressCount[address] = 1
			}
		}
	}

	allUncoveringAddress := map[uint32]*pb.UncoveredAddress{}
	for address, count := range uncoveringAddressCount {
		if count == resultSize {
			allUncoveringAddress[address] = tempUA[address]
		}
	}

	for address, uncoveringAddress := range allUncoveringAddress {
		ress := ""
		for _, r := range d.resultsWithDra.result {
			ress = r.checkUncoveredAddress(address)
		}

		path := filepath.Join(d.path, fmt.Sprintf("0xffffffff%x.txt", uncoveringAddress.ConditionAddress-5))
		if _, err := os.Stat(path); err == nil {
			// path/to/whatever exists

		} else if os.IsNotExist(err) {
			// path/to/whatever does *not* exist
			ff, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			_, _ = ff.WriteString(ress)
			_ = ff.Close()

		} else {
			// Schrodinger: file may or may not exist. See err for details.
			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		}

	}

	uaStatus := map[pb.TaskStatus]uint32{}
	for _, uaa := range allUncoveringAddress {

		uaStatus[uaa.RunTimeDate.TaskStatus]++
	}
	res += "*******************************************\n"
	for ts, c := range uaStatus {
		res += ts.String() + fmt.Sprintf("%5d", c) + "\n"
	}
	res += "*******************************************\n"

	f, _ := os.OpenFile(d.dataPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	_, _ = f.WriteString(res)
	_ = f.Close()

	res = ""
	sort.Slice(allUncoveringAddress, func(i, j int) bool {
		return allUncoveringAddress[uint32(i)].NumberDominatorInstructions < allUncoveringAddress[uint32(j)].NumberDominatorInstructions
	})
	res += "UncoveringAddress @Inst@Input@WA @task @Tested@ Count "
	res += fmt.Sprintf("@%25s", "Kind")
	for _, uaa := range allUncoveringAddress {
		res += fmt.Sprintf("0xffffffff%x", uaa.UncoveredAddress-5)
		res += fmt.Sprintf("@%4d", uaa.NumberDominatorInstructions)
		res += fmt.Sprintf("@%5d", len(uaa.Input))
		res += fmt.Sprintf("@%3d", len(uaa.WriteAddress))
		count := uint32(0)
		for _, c := range uaa.TasksCount {
			count += c
		}
		res += fmt.Sprintf("@%5d", count)
		count -= uaa.TasksCount[int32(pb.TaskStatus_untested)]
		res += fmt.Sprintf("@%6d", count)
		if uaa.RunTimeDate == nil {

		} else {
			res += fmt.Sprintf("@%7d", uaa.RunTimeDate.RecursiveCount)
			res += fmt.Sprintf("@%25s", uaa.RunTimeDate.TaskStatus.String())

		}
		res += "\n"
	}
	f, _ = os.OpenFile(filepath.Join(d.path, "uncovering_more.txt"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	_, _ = f.WriteString(res)
	_ = f.Close()

	if d.a2i {
		_ = os.Remove(filepath.Join(d.path, fmt.Sprintf("write.txt")))
		f, _ := os.OpenFile(filepath.Join(d.path, fmt.Sprintf("write.txt")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

		for address := range allWriteAddress {
			_, _ = f.WriteString(fmt.Sprintf("0xffffffff%x\n", address-5))
		}
		_ = f.Close()

		_ = os.Remove(filepath.Join(d.path, fmt.Sprintf("uncovering.txt")))
		f, _ = os.OpenFile(filepath.Join(d.path, fmt.Sprintf("uncovering.txt")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

		for _, uncoveringAddress := range allUncoveringAddress {
			_, _ = f.WriteString(fmt.Sprintf("0xffffffff%x&0xffffffff%x\n", uncoveringAddress.ConditionAddress-5, uncoveringAddress.UncoveredAddress-5))
		}
		_ = f.Close()

		_ = os.Chdir(d.path)
		err := filepath.Walk(d.path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.HasPrefix(info.Name(), "0x") {
					_ = os.Remove(path)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
		for address, uncoveringAddress := range allUncoveringAddress {
			ff, _ := os.OpenFile(filepath.Join(d.path, fmt.Sprintf("0xffffffff%x.txt", uncoveringAddress.ConditionAddress-5)), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			for _, r := range d.resultsWithDra.result {
				ress := r.checkUncoveredAddress(address)
				_, _ = ff.WriteString(ress)
			}
			_ = ff.Close()
		}

		_ = os.Chdir(d.path)
		err = filepath.Walk(d.path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if strings.HasPrefix(info.Name(), "condition") {
					_ = os.Remove(path)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
		cmd := exec.Command(pb.PathA2i, "-asm="+pb.FileAsm, "-objdump="+pb.FileVmlinuxObjdump, "-bc="+pb.FileBc, pb.FileDRAConfig)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		log.Println("cmd : ")
		log.Println(cmd.String())
		err = cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}

}

func (d *device) checkStatistic() {

	name := d.baseName

	path := filepath.Join(d.path, "statistic.txt")
	_ = os.Remove(path)

	f := func(gs func(r *result) *statistic) {
		var ss []*statistic
		for _, r := range d.resultsWithDra.result {
			tempS := gs(r)
			tempS.output(filepath.Join(d.path, tempS.Name+".txt"))
			ss = append(ss, tempS)
		}
		s := average(ss)
		s.Name = name
		s.output(path)
	}

	f(prevalent)
	f(writeStatement)
	f(controlFlow)
	f(unstable)
	f(recursive)
}
