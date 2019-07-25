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

	name := "data.bin"
	if len(os.Args) == 2 {
		name = os.Args[1]
	}
	
	// [START unmarshal_proto]
	in, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	corpus := &pb.Corpus{}
	if err := proto.Unmarshal(in, corpus); err != nil {
		log.Fatalln("Failed to parse corpus:", err)
	}

	path := "./data.txt"
	_ = os.Remove(path)
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	_, _ = f.WriteString(fmt.Sprintf("%v", corpus))
}
