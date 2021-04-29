#!/bin/bash
git submodule update --init --recursive
make clean
make HAS_SYSTEM_PROTOBUF=false -j12
make install prefix=/home/yu/data/2018-Dependency/build