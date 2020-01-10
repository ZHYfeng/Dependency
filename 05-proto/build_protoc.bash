#!/bin/bash
echo "[*] Trying to generate protoc"
protoc --cpp_out=../02-dependency/lib/RPC ./a2l.proto
protoc --cpp_out=../02-dependency/lib/RPC ./DependencyRPC.proto

protoc --grpc_out=../02-dependency/lib/RPC --plugin=protoc-gen-grpc="$HOME"/data/build/bin/grpc_cpp_plugin ./DependencyRPC.proto
# protoc --grpc_out=../02-dependency/lib/RPC --plugin=protoc-gen-grpc=/use/local/bin/grpc_cpp_plugin ./DependencyRPC.proto

protoc -I=. --python_out=../04-script/default ./DependencyRPC.proto

protoc --go_out=plugins=grpc:../03-syzkaller/pkg/dra ./DependencyRPC.proto