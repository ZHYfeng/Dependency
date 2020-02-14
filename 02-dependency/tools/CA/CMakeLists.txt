add_executable(ca ca.cpp)
target_link_libraries(ca PUBLIC
        DCC DMM RPC JSON STA
        ${llvm_libs}
        PkgConfig::YOUR_PKG)