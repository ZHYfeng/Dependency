cmake_minimum_required(VERSION 3.4.3)
project(2018_dependency CXX)

set(CMAKE_CXX_STANDARD 11)

set(CMAKE_CXX_FLAGS "-std=gnu++11 -O3 -pthread -I/usr/local/include")



if (DEFINED ENV{LLVM_DIR})
    find_package(LLVM REQUIRED CONFIG)
    list(APPEND CMAKE_MODULE_PATH "${LLVM_CMAKE_DIR}")
    include(AddLLVM)
    add_definitions(${LLVM_DEFINITIONS})
    include_directories(${LLVM_INCLUDE_DIRS})
else ()
    message(FATAL_ERROR "\
WARNING: The LLVM_DIR var was not set (required for an out-of-source build)!\n\
Please set this to environment variable to point to the LLVM build directory\
(e.g. on linux: export LLVM_DIR=/path/to/llvm/build/dir)")
endif ()

add_subdirectory(lib)
add_subdirectory(tools)
