add_executable(a2i a2i.cpp)
target_link_libraries(a2i PUBLIC
        DCC DRA RPC JSON STA
        ${llvm_libs}
        PkgConfig::YOUR_PKG)