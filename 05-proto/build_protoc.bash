#!/bin/bash
echo "[*] Trying to generate protoc"

PROTO=$PWD

protoc --cpp_out=../02-dependency/lib/RPC ./*.proto
protoc --grpc_out=../02-dependency/lib/RPC --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` ./DependencyRPC.proto

protoc --go_out=../03-syzkaller/pkg/dra --go_opt=paths=source_relative \
    --go-grpc_out=../03-syzkaller/pkg/dra --go-grpc_opt=paths=source_relative \
 --go_opt=MBase.proto=$PROTO Base.proto \
 --go_opt=MInput.proto=$PROTO Input.proto \
 --go_opt=MStatistic.proto=$PROTO Statistic.proto \
 --go_opt=MTask.proto=$PROTO Task.proto \
 --go_opt=MDependency.proto=$PROTO Dependency.proto \
 --go_opt=MDependencyRPC.proto=$PROTO DependencyRPC.proto