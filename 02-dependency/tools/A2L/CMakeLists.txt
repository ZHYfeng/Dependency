add_executable(a2l a2l.cpp)
target_link_libraries(a2l PUBLIC
        RPC
        ${llvm_libs}
        PkgConfig::YOUR_PKG)