set(PROTO DependencyRPC.proto)
set(PROTOC DependencyRPC.grpc.pb.cc
        DependencyRPC.pb.cc)

set(SOURCES
        ${PROTOC}
        DependencyRPCClient.cpp)

add_custom_command(
        OUTPUT ${PROTOC}
        COMMAND protoc -I ${CMAKE_CURRENT_SOURCE_DIR} --cpp_out=. --grpc_out=. --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin ${CMAKE_CURRENT_SOURCE_DIR}/${PROTO}
        DEPENDS ${PROTO}
)

add_library(RPC ${SOURCES})