#!/bin/bash
export LLVM_DIR=/home/yuh/data/build/llvm-7.0.0.src/build
git pull
cd ./lib/RPC/
./build-protoc.sh
cd ../../
echo "[*] Trying to Run Cmake"
mkdir build
cd build
cmake -DCMAKE_BUILD_TYPE=Debug -G "CodeBlocks - Unix Makefiles" ..
echo "[*] Trying to make"
make -j8
