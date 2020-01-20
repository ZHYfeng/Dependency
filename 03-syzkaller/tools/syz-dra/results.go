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
}

func (r *results) read(path string) {
	r.path = path
	r.dirName = filepath.Dir(path)
	r.baseName = filepath.Base(path)

	err := filepath.Walk(r.path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == nameStatistic {
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
}
