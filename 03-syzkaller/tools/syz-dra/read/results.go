package read

import (
	pb "github.com/ZHYfeng/2018_dependency/03-syzkaller/pkg/dra"
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
			if info.Name() == pb.NameStatistics {
				temp := &result{}
				r.result = append(r.result, temp)

				temp.read(filepath.Dir(path))
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	r.maxCoverage = map[uint32]uint32{}
	for _, rr := range r.result {
		for a, c := range rr.statistics.Coverage.Coverage {
			r.maxCoverage[a] += c
		}
	}
}
