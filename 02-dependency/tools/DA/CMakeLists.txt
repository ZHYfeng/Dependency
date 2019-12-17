add_executable(da da.cpp)
target_link_libraries(da PUBLIC
        DCC DRA RPC JSON STA
        ${llvm_libs}
        PkgConfig::YOUR_PKG)