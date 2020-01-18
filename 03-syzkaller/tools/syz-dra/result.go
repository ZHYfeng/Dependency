package main

import (
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"path/filepath"
)

type result struct {
	path      string
	dirName   string
	baseName  string
	data      *pb.Data
	statistic *pb.Statistic
}

func (r *result) read(path string) {
	r.path = path
	r.dirName = filepath.Dir(path)
	r.baseName = filepath.Base(path)

	fileName := filepath.Join(r.path, nameData)
	in, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	if err := proto.Unmarshal(in, r.data); err != nil {
		log.Fatalln("Failed to parse data:", err)
	}

	fileName = filepath.Join(r.path, nameStatistic)
	in, err = ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	if err := proto.Unmarshal(in, r.statistic); err != nil {
		log.Fatalln("Failed to parse statistic:", err)
	}
}
