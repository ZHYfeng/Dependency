#!/bin/bash
echo "[*] Trying to generate protoc"

PROTO=$PWD
export PATH=/home/yu/data/2018-Dependency/build/bin:$PATH

protoc --cpp_out=../02-dependency/lib/RPC ./*.proto
# protoc --grpc_out=../02-dependency/lib/RPC --plugin=protoc-gen-grpc=/home/yu/data/2018-Dependency/build/bin/grpc_cpp_plugin ./DependencyRPC.proto
protoc --grpc_out=../02-dependency/lib/RPC --plugin=protoc-gen-grpc=/home/yhao016/data/18-Dependency/install/bin/grpc_cpp_plugin ./DependencyRPC.proto

protoc --go_out=plugins=grpc:../03-syzkaller/pkg/dra \
 --go_opt=MBase.proto=$PROTO Base.proto \
 --go_opt=MInput.proto=$PROTO Input.proto \
 --go_opt=MStatistic.proto=$PROTO Statistic.proto \
 --go_opt=MTask.proto=$PROTO Task.proto \
 --go_opt=MDependency.proto=$PROTO Dependency.proto \
 --go_opt=MDependencyRPC.proto=$PROTO DependencyRPC.proto