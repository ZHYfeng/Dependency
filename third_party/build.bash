#!/bin/bash

sudo apt install -y python3-pip python3-tk autoconf automake libtool curl make g++ unzip build-essential autoconf libtool pkg-config libgflags-dev libgtest-dev libc++-dev
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u google.golang.org/grpc python3-protobuf protobuf scipy

PATH_PROTOBUF=protobuf
PATH_GRPC=grpc

cd $PATH_PROTOBUF || exit
bash ../build-protocbuf.bash
cd ..
cd $PATH_GRPC || exit
bash ../build-grpc.bash
cd ..
