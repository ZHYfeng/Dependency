#!/bin/bash

make clean
make HAS_SYSTEM_PROTOBUF=false -j12
make install prefix=$III