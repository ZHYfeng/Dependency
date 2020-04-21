#!/bin/bash

if [ -z "$III" ]
then
    III=$HOME/data/build
else
    echo "\$III is " $III
fi

III=$HOME/data/build
make clean
make HAS_SYSTEM_PROTOBUF=false -j12
make install prefix=$III