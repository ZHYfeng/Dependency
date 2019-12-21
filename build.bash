#!/bin/bash

PATH_GIT=$PWD
PATH_THIRD=$PATH_GIT/third_party
PATH_PROTO=$PATH_GIT/05-proto
PATH_DRA=$PATH_GIT/02-dependency
PATH_SYZKALLER=$PATH_GIT/03-syzkaller
PATH_SCRIPT=$PATH_GIT/04-script

export III=$HOME/data/build

git pull

if [ -d $III ]; then
    echo "[*] $III exist"
else 
    mkdir $III
    git submodule update --init --recursive
    cd $PATH_THIRD || exit
    bash ./build.bash
    cd ..
fi

ldconfig -C $PATH_GIT/ld.so.cache

cd $PATH_PROTO || exit
bash ./build-protoc.bash
cd ..
cd $PATH_DRA || exit
bash ./build.bash
cd ..
cd $PATH_SYZKALLER || exit
bash ./build.bash
cd ..
cd $PATH_SCRIPT || exit
bash ./build.bash
cd ..

if [ -d $HOME/data/work ]; then
    echo "[*] $HOME/data/work exist"
else
    echo "[*] Trying to generate work dir"
    python3 04-script/main.py generate ~/data/work
fi
