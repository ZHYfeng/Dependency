llvm_map_components_to_libnames(llvm_libs support)

add_executable(a2l a2l.cpp)
target_link_libraries(a2l RPC ${llvm_libs} "-L/usr/local/lib -lprotobuf -lpthread")

#pkg-config --cflags --libs protobuf grpc++ grpc