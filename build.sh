#!/bin/bash
export LLVM_DIR=/home/yuh/data/build/llvm-7.0.0.src/build
cd ./lib/RPC/
./build.sh
cd ../../
echo "[*] Trying to Run Cmake"
mkdir build
cd build
cmake ..
echo "[*] Trying to make"
make -j8
