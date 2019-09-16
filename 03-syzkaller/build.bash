#!/bin/bash
git pull
cd ./pkg/dra/ || exit
./build-protoc.sh
cd ../../
echo "[*] Trying to make"
make -j8
