#!/bin/bash

PATH_GIT=/home/yuh/data/git
PATH_DRA=$PATH_GIT/2018_dependency
PATH_SYZKALLER=$PATH_GIT/gopath/src/github.com/google/syzkaller
cd $PATH_DRA || exit
git pull
./build.bash
cd $PATH_SYZKALLER || exit
git pull
./build.bash
