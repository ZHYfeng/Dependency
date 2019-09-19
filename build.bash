#!/bin/bash

git pull

PATH_GIT=.
PATH_PROTO=$PATH_GIT/05-proto
PATH_DRA=$PATH_GIT/02-dependency
PATH_SYZKALLER=$PATH_GIT/03-syzkaller
PATH_SCRIPT=$PATH_GIT/04-script

cd $PATH_PROTO
bash ./build-protoc.bash
cd ..
cd $PATH_DRA || exit
bash ./build.bash
cd ..
cd $PATH_SYZKALLER || exit
bash ./build.bash
cd ..
# cd $PATH_SCRIPT || exit
# bash ./build.bash
# cd ..

echo "[*] Trying to generate work dir"
python3 04-script/main.py generate ~/data/work
