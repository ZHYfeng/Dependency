#!/bin/bash

protoc --cpp_out=. --grpc_out=. --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin ./DependencyRPC.proto ./a2l.proto