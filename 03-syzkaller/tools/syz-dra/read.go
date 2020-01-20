package main

import (
	"path/filepath"
	"strings"
)

const (
	nameDevive     = "dev_"
	nameBase       = "base"
	nameWithDra    = "01-result-with-dra"
	nameWithoutDra = "02-result-without-dra"
	nameData       = "data.bin"
	nameStatistic  = "statistics.bin"
	nameDataResult = "data.txt"
)

func main() {
	//fileName := "/home/yhao016/data/work/dev_kvm/01-result-with-dra/0/data.bin"
	//in, err := ioutil.ReadFile(fileName)
	//if err != nil {
	//	log.Fatalln("Error reading file:", err)
	//}
	//corpus := &pb.Data{}
	//if err := proto.Unmarshal(in, corpus); err != nil {
	//	log.Fatalln("Failed to parse corpus:", err)
	//}
	//fmt.Printf("size : %d\n", len(corpus.BootTask.TaskArray))
	//fmt.Printf("corpus.Tasks.Tasks size : %d\n", len(corpus.Tasks.TaskArray))
	//fmt.Printf("size : %d\n", len(corpus.HighTask.TaskArray))
	//fmt.Printf("size : %d\n", len(corpus.Input))
	//fmt.Printf("size : %d\n", len(corpus.NewInput))
	//fmt.Printf("size : %d\n", len(corpus.UncoveredAddress))
	//fmt.Printf("size : %d\n", len(corpus.WriteAddress))

	tes := results{}
	tes.read(".")
}

func read(path string) {
	baseName := filepath.Base(path)
	if strings.HasPrefix(baseName, nameDevive) {

	}
}
