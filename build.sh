#!/bin/bash
export LLVM_DIR=/home/hang/dep/dr_checker/helper_scripts/drchecker_dep/llvm/build/../cmake/
echo "[*] Trying to Run Cmake"
mkdir build
cd build
cmake ..
echo "[*] Trying to make"
make -j8
