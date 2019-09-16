#!/bin/bash

git pull

PATH_GIT=.
PATH_PROTO=$PATH_GIT/05-proto
PATH_DRA=$PATH_GIT/02-dependency
PATH_SYZKALLER=$PATH_GIT/03-syzkaller

cd $PATH_PROTO
bash ./build-protoc.bash
cd ..
cd $PATH_DRA || exit
bash ./build.bash
cd ..
cd $PATH_SYZKALLER || exit
bash ./build.bash
cd ..
