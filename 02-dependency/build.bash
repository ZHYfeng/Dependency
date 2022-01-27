#!/bin/bash

echo "[*] Trying to build dependency"
rm -rf cmake-build
mkdir cmake-build
cd cmake-build || exit
if [[ $PATH_PROJECT != "" ]]
then
    echo "[*] install in $PATH_PROJECT/install"
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=$PATH_PROJECT/install ..
    make -j
    make install
else
    echo "[*] just build, not install"
    cmake -DCMAKE_BUILD_TYPE=Release ..
    make -j
fi