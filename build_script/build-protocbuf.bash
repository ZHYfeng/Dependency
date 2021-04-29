#!/bin/bash
./autogen.sh
./configure --prefix=/home/yu/data/2018-Dependency/build  --disable-shared
make clean
make -j12
make install