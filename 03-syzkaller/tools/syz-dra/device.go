package main

import (
	"path/filepath"
)

type device struct {
	path     string
	dirName  string
	baseName string

	base              *result
	resultsWithDra    *results
	resultsWithoutDra *results
}

func (d *device) read(path string) {
	d.path = path
	d.dirName = filepath.Dir(path)
	d.baseName = filepath.Base(path)

	d.resultsWithDra.read(filepath.Join(d.path, nameWithDra))
	d.resultsWithoutDra.read(filepath.Join(d.path, nameWithoutDra))
}
