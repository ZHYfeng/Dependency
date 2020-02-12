package main

import (
	"fmt"
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	main2 "github.com/ZHYfeng/2018_dependency/03-syzkaller/tools/syz-dra/read"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) == 2 {
		read(os.Args[1])
	}
}

func read(path string) {
	baseName := filepath.Base(path)
	if strings.HasPrefix(baseName, pb.NameDevice) {
		fmt.Printf("nameDevice\n")
		d := &main2.device{}
		d.read(path)
	} else if strings.HasPrefix(baseName, pb.NameWithDra) || strings.HasPrefix(baseName, pb.NameWithoutDra) {
		fmt.Printf("nameWithDra or NameWithoutDra\n")
		r := &main2.results{}
		r.read(path)
	} else {
		fmt.Printf("readUnstableInput\n")
		readUnstableInput(path)
	}
}

func readUnstableInput(path string) {

	//in, err := ioutil.ReadFile(path)
	//if err != nil {
	//	log.Fatalln("Error reading file:", err)
	//}
	//ui := &pb.UnstableInput{}
	//if err := proto.Unmarshal(in, ui); err != nil {
	//	log.Fatalln("Failed to parse ui:", err)
	//}
	//fmt.Printf("0xffffffff%x\n", ui.Address-5)
	//fmt.Printf("%d\n", ui.Idx)
	//fmt.Printf("%s\n", ui.Program)
	//for idx, path := range ui.NewPath {
	//	fmt.Printf("check %d NewPath\n", idx)
	//	for index, p := range path.Path {
	//		fmt.Printf("check %d path\n", index)
	//		for _, a := range p.Address {
	//			if a == ui.Address {
	//				fmt.Printf("find address in %d newpath %d path\n", idx, index)
	//			}
	//		}
	//	}
	//}
	//
	//fmt.Printf("check UnstablePath\n")
	//for index, p := range ui.UnstablePath {
	//	fmt.Printf("check %d path\n", index)
	//	for _, a := range p.Address {
	//		if a == ui.Address {
	//			fmt.Printf("find address in UnstablePath %d path\n", index)
	//		}
	//	}
	//}
}
