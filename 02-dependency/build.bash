#!/bin/bash

echo "[*] Trying to build dependency"
rm -rf cmake-build
mkdir cmake-build
cd cmake-build || exit
if [[ $PROJECT_PATH != "" ]]
then
    echo "[*] install in $PROJECT_PATH/install"
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=$PROJECT_PATH/install ..
else
    echo "[*] install in $HOME/install"
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=$HOME/install ..
fi

make -j
make install