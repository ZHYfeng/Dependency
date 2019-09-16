// Copyright 2017 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

//go:generate go run fidlgen/main.go

package fuchsia

import (
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/prog"
	"github.com/ZHYfeng/2018_dependency/03-syzkaller/sys/targets"
)

func InitTarget(target *prog.Target) {
	target.MakeMmap = targets.MakeSyzMmap(target)
}
