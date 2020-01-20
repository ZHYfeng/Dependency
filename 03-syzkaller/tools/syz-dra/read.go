package main

import (
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
	if len(os.Args) == 1 {
		read(os.Args[1])
	}
}

func read(path string) {
	baseName := filepath.Base(path)
	if strings.HasPrefix(baseName, nameDevice) {
		d := &device{}
		d.read(path)
	} else if strings.HasPrefix(baseName, nameWithDra) || strings.HasPrefix(baseName, nameWithDra) {
		r := &results{}
		r.read(path)
	} else {

	}
}
