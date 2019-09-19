#!/bin/bash
export LLVM_DIR=$HOME/data/build/llvm-7.0.0.src/build
export CXX=$HOME/data/build/llvm-7.0.0.src/build/bin/clang
echo "[*] Trying to Run Cmake"
rm -rf build
mkdir build
cd build || exit
cmake -DCMAKE_BUILD_TYPE=Debug -G "CodeBlocks - Unix Makefiles" ..
echo "[*] Trying to make"
make -j8
