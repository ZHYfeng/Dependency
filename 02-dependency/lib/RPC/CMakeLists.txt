set(PROTOC
        DependencyRPC.grpc.pb.cc
        DependencyRPC.pb.cc
        )

set(APROTOC
        a2l.pb.cc
        )

set(SOURCES
        ${PROTOC}
        ${APROTOC}
        DependencyRPCClient.cpp)

add_library(RPC ${SOURCES})