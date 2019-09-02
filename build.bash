#!/bin/bash
export LLVM_DIR=/home/yuh/data/build/llvm-7.0.0.src/build
git pull
cd ./lib/RPC/ || exit
./build-protoc.sh
cd ../../
echo "[*] Trying to Run Cmake"
rm -rf build
mkdir build
cd build || exit
cmake -DCMAKE_BUILD_TYPE=Debug -G "CodeBlocks - Unix Makefiles" ..
echo "[*] Trying to make"
make -j8
