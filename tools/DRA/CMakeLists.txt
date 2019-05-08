llvm_map_components_to_libnames(llvm_libs bitwriter core ipo irreader instcombine instrumentation target linker analysis scalaropts support)

add_executable(dra dra.cpp)
target_link_libraries(dra DCC DRA RPC JSON STA ${llvm_libs} "-L/usr/local/lib -lprotobuf -lpthread -lgrpc++ -lgrpc")

