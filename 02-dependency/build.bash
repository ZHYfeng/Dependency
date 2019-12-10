#!/bin/bash
rm -rf cmake-build-debug
mkdir cmake-build-debug
cd cmake-build-debug || exit
cmake -DCMAKE_BUILD_TYPE=Debug ..
echo "[*] Trying to make"
make -j8
