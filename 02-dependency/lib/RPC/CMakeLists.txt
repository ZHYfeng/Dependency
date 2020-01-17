set(PROTOC
        DependencyRPC.grpc.pb.cc
        DependencyRPC.pb.cc
        Base.pb.cc
        Dependency.pb.cc
        Task.pb.cc
        Input.pb.cc
        Statistic.pb.cc
        )

set(SOURCES
        ${PROTOC}
        DependencyRPCClient.cpp)

add_library(RPC ${SOURCES})