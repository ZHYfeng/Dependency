llvm_map_components_to_libnames(llvm_libs bitwriter core ipo irreader instcombine instrumentation target linker analysis scalaropts support)

add_executable(a2i a2i.cpp)
target_link_libraries(a2i DCC DRA RPC JSON STA ${llvm_libs} "-L/usr/local/lib -lprotobuf -lpthread -lgrpc++ -lgrpc")

#pkg-config --cflags --libs protobuf grpc++ grpc