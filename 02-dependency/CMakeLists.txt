cmake_minimum_required(VERSION 3.4.3)
project(2018_dependency CXX)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_FLAGS "-std=gnu++11 -O3 -pthread")

find_package(LLVM REQUIRED CONFIG)
include_directories(${LLVM_INCLUDE_DIRS})
add_definitions(${LLVM_DEFINITIONS})
llvm_map_components_to_libnames(llvm_libs bitwriter core ipo irreader instcombine instrumentation target linker analysis scalaropts support)

find_package(PkgConfig REQUIRED)
pkg_check_modules(YOUR_PKG REQUIRED IMPORTED_TARGET grpc grpc++ protobuf)
#pkg-config --cflags --libs protobuf grpc++ grpc

add_subdirectory(lib)
add_subdirectory(tools)
