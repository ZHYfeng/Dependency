package main

import (
	"fmt"
	"io/ioutil"
	"log"

	pb "./dra"
	"github.com/golang/protobuf/proto"
)

func main() {

	// name := "data.bin"
	// if len(os.Args) == 2 {
	// 	name = os.Args[1]
	// }

	// // [START unmarshal_proto]
	// in, err := ioutil.ReadFile(name)
	// if err != nil {
	// 	log.Fatalln("Error reading file:", err)
	// }
	// corpus := &pb.Corpus{}
	// if err := proto.Unmarshal(in, corpus); err != nil {
	// 	log.Fatalln("Failed to parse corpus:", err)
	// }

	// path := "./data.txt"
	// _ = os.Remove(path)
	// f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// defer f.Close()
	// _, _ = f.WriteString(fmt.Sprintf("%v", corpus))

	coverageWithDra := "./result-with-dra/coverage.bin"
	in, err := ioutil.ReadFile(coverageWithDra)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	coverage1 := &pb.Coverage{}
	if err := proto.Unmarshal(in, coverage1); err != nil {
		log.Fatalln("Failed to parse coverage1:", err)
	}

	coverageWithoutDra := "./result-without-dra/coverage.bin"
	in, err = ioutil.ReadFile(coverageWithoutDra)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	coverage2 := &pb.Coverage{}
	if err := proto.Unmarshal(in, coverage2); err != nil {
		log.Fatalln("Failed to parse coverage2:", err)
	}

	coverage3 := &pb.Coverage{
		Coverage: make(map[uint32]uint32),
	}
	fmt.Printf("size c2 : %d\n", len(coverage2.Coverage))
	for a := range coverage2.Coverage {
		if k, ok := coverage3.Coverage[a]; ok {
			coverage3.Coverage[a] = 1<<2 + k
		} else {
			coverage3.Coverage[a] = 1 << 2
		}
	}
	fmt.Printf("size c1 : %d\n", len(coverage1.Coverage))
	for a, d := range coverage1.Coverage {
		if k, ok := coverage3.Coverage[a]; ok {
			coverage3.Coverage[a] = 1<<1 + k + d
		} else {
			coverage3.Coverage[a] = 1<<1 + d
		}
	}

	both := 0
	c1 := 0
	c2 := 0
	c1d := 0
	for a, k := range coverage3.Coverage {
		if k == 6 {
			both++
		} else if k == 7 {
			// c1d++
			both++
		} else if k == 4 {
			c2++
		} else if k == 2 {
			c1++
		} else if k == 3 {
			c1++
			c1d++
			fmt.Printf("dependency : %x\n", a)
		} else {
			fmt.Printf("%v\n", a)
		}
	}

	fmt.Printf("both : %d\n", both)
	fmt.Printf("result with dra : %d\n", c1)
	fmt.Printf("result without dra : %d\n", c2)
	fmt.Printf("result with dra of d : %d\n", c1d)
}
