package main

import (
	"fmt"
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	nameDevice     = "dev_"
	nameBase       = "base"
	nameWithDra    = "01-result-with-dra"
	nameWithoutDra = "02-result-without-dra"
	nameData       = "data.bin"
	nameStatistics = "statistics.bin"
	nameDataResult = "data.txt"
)

func main() {
	if len(os.Args) == 2 {
		read(os.Args[1])
	}
}

func read(path string) {
	baseName := filepath.Base(path)
	if strings.HasPrefix(baseName, nameDevice) {
		fmt.Printf("nameDevice\n")
		d := &device{}
		d.read(path)
	} else if strings.HasPrefix(baseName, nameWithDra) || strings.HasPrefix(baseName, nameWithoutDra) {
		fmt.Printf("nameWithDra\n")
		r := &results{}
		r.read(path)
	} else {
		fmt.Printf("readUnstableInput\n")
		readUnstableInput(path)
	}
}

func readUnstableInput(path string) {

	in, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	ui := &pb.UnstableInput{}
	if err := proto.Unmarshal(in, ui); err != nil {
		log.Fatalln("Failed to parse ui:", err)
	}
	fmt.Printf("%d\n", ui.Idx)
	fmt.Printf("%s\n", ui.Program)
	for idx, path := range ui.NewPath {
		fmt.Printf("check %d NewPath\n", idx)
		for index, p := range path.Path {
			fmt.Printf("check %d path\n", index)
			for _, a := range p.Address {
				if a == ui.Address {
					fmt.Printf("find address in %d newpath %d path\n", idx, index)
				}
			}
		}
	}

	fmt.Printf("check UnstablePath\n")
	for index, p := range ui.UnstablePath {
		fmt.Printf("check %d path\n", index)
		for _, a := range p.Address {
			if a == ui.Address {
				fmt.Printf("find address in UnstablePath %d path\n", index)
			}
		}
	}
}
