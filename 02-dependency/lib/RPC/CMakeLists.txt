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
target_link_libraries(RPC PUBLIC
        ${_REFLECTION}
        ${_GRPC_GRPCPP}
        ${_PROTOBUF_LIBPROTOBUF}
        )