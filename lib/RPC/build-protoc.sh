#!/bin/bash

protoc --grpc_out=. --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin ./DependencyManger.proto
protoc --cpp_out=.  ./DependencyManger.proto