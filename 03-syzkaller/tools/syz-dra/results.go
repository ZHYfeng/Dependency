package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type results struct {
	path     string
	dirName  string
	baseName string
	result   []*result

	maxCoverage map[uint32]uint32
}

func (r *results) read(path string) {
	r.path = path
	r.dirName = filepath.Dir(path)
	r.baseName = filepath.Base(path)

	r.result = []*result{}
	err := filepath.Walk(r.path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == nameStatistics {
				temp := &result{}
				r.result = append(r.result, temp)

				temp.read(filepath.Dir(path))
				fmt.Println(info.Name(), info.Size())
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	r.maxCoverage = map[uint32]uint32{}
	for _, rr := range r.result {
		for a := range rr.statistics.Coverage.Coverage {
			r.maxCoverage[a]++
		}
	}
}
