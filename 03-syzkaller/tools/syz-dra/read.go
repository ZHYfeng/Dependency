package main

import (
	"fmt"
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
)

func main() {
	fileName := "/home/yuh/data/work/dev_snd_seq/01-result-with-dra/0/data.bin"
	in, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	corpus := &pb.Corpus{}
	if err := proto.Unmarshal(in, corpus); err != nil {
		log.Fatalln("Failed to parse corpus:", err)
	}
	fmt.Printf("size : %d\n", len(corpus.BootTask.Tasks))
	fmt.Printf("corpus.Tasks.Tasks size : %d\n", len(corpus.Tasks.Tasks))
	fmt.Printf("size : %d\n", len(corpus.HighTask.Tasks))
	fmt.Printf("size : %d\n", len(corpus.Input))
	fmt.Printf("size : %d\n", len(corpus.NewInput))
	fmt.Printf("size : %d\n", len(corpus.UncoveredAddress))
	fmt.Printf("size : %d\n", len(corpus.WriteAddress))
}
