#!/bin/bash
export LLVM_DIR=$HOME/data/build/llvm-7.0.0.src/build
export CXX=$HOME/data/build/llvm-7.0.0.src/build/bin/clang
echo "[*] Trying to Run Cmake"
rm -rf build
mkdir build
cd build || exit
cmake ..
echo "[*] Trying to make"
make -j8
