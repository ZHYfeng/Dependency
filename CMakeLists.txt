cmake_minimum_required(VERSION 3.4.3)
project(2018_dependency CXX)

set(CMAKE_CXX_STANDARD 11)

set(CMAKE_CXX_FLAGS "-std=gnu++11 -O3 -pthread -I/usr/local/include")

find_package(LLVM REQUIRED CONFIG)
list(APPEND CMAKE_MODULE_PATH "${LLVM_CMAKE_DIR}")
include(AddLLVM)

add_subdirectory(lib)
add_subdirectory(tools)
