#!/bin/bash

protoc --cpp_out=../02-dependency/lib/RPC ./a2l.proto
protoc --cpp_out=../02-dependency/lib/RPC ./DependencyRPC.proto

protoc --grpc_out=../02-dependency/lib/RPC --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin ./DependencyRPC.proto

protoc -I=. --python_out=../04-script/result ./DependencyRPC.proto

protoc --go_out=plugins=grpc:../03-syzkaller/pkg/dra ./DependencyRPC.proto