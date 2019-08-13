package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	pb "./dra"
	"github.com/golang/protobuf/proto"
)

func main() {

	statBin := os.Args[1]
	statTxt := os.Args[2]
	in, err := ioutil.ReadFile(statBin)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	stat := &pb.Statistics{}
	if err := proto.Unmarshal(in, stat); err != nil {
		log.Fatalln("Failed to parse stat:", err)
	}

	fmt.Printf("stat : %v\n", stat.Stat)

	_ = os.Remove(statTxt)
	f, _ := os.OpenFile(statTxt, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	_, _ = f.WriteString(fmt.Sprintf("%s %s\n", statBin, statBin))

}
