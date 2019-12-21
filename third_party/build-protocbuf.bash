#!/bin/bash

./autogen.sh
./configure --prefix=$III  --disable-shared
make clean
make -j12
make install