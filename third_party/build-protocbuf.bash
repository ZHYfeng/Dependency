#!/bin/bash

if [ -z "$III" ]
then
    III=$HOME/data/build
else
    echo "\$III is " $III
fi

./autogen.sh
./configure --prefix=$III  --disable-shared
make clean
make -j12
make install