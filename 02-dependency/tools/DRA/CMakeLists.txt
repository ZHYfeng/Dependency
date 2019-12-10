add_executable(dra dra.cpp)
target_link_libraries(dra PUBLIC
        DCC DRA RPC JSON STA
        ${llvm_libs}
        PkgConfig::YOUR_PKG)