#!/bin/bash

PATH_GIT=$PWD
PATH_THIRD=$PATH_GIT/third_party
PATH_PROTO=$PATH_GIT/05-proto
PATH_DRA=$PATH_GIT/02-dependency
PATH_SYZKALLER=$PATH_GIT/03-syzkaller
PATH_SCRIPT=$PATH_GIT/04-script

export III=$HOME/data/build

cd $PATH_PROTO || exit
bash ./remove.bash
cd ..
git pull

if [ -d $III ]
then
    echo "[*] $III exist"
else 
    mkdir $III
fi

if ! [ -x "$(command -v grpc_cpp_plugin)" ]; then
    git submodule update --init --recursive
    cd $PATH_THIRD || exit
    bash ./build.bash
    cd ..
    ldconfig -C $PATH_GIT/ld.so.cache
fi

cd $PATH_PROTO || exit
bash ./build_protoc.bash
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
    python3 04-script/python/main.py generate ~/data/work
fi

# sudo apt-get install autoconf automake libtool curl make g++ unzip