set(PROTO DependencyManager.proto)

set(SOURCES
        DependencyManager.grpc.pb.cc
        DependencyManager.pb.cc
        )

add_custom_command(
        OUTPUT ${SOURCES}
        COMMAND protoc -I ${CMAKE_CURRENT_SOURCE_DIR} --cpp_out=. --grpc_out=. --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin ${CMAKE_CURRENT_SOURCE_DIR}/${PROTO}
        DEPENDS ${PROTO}
)

add_library(RPC ${SOURCES})