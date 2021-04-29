#!/bin/bash

PATH_GIT=$PWD
PATH_PROTO=$PATH_GIT/05-proto
PATH_DRA=$PATH_GIT/02-dependency
PATH_SYZKALLER=$PATH_GIT/03-syzkaller

export III=$HOME/data/2018-Dependency/build

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

cd $PATH_PROTO || exit
bash ./build_protoc.bash
cd ..
cd $PATH_DRA || exit
bash ./build.bash
cd ..
cd $PATH_SYZKALLER || exit
bash ./build.bash
cd ..

# if [ -d $HOME/data/work ]; then
#     echo "[*] $HOME/data/work exist"
# else
#     echo "[*] Trying to generate work dir"
#     python3 04-script/python/main.py generate ~/data/work
# fi

# sudo apt-get install autoconf automake libtool curl make g++ unzip