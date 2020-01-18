package main

import "path/filepath"

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
	filepath.Walk(r.path)
}
