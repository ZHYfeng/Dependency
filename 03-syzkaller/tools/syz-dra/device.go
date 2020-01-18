package main

import "path/filepath"

type device struct {
	path     string
	dirName  string
	baseName string
}

func (d *device) read(path string) {
	d.path = path
	d.dirName = filepath.Dir(path)
	d.baseName = filepath.Base(path)

}
