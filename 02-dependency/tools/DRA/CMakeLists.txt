add_executable(dra dra.cpp)
target_link_libraries(dra PUBLIC
        DCC DMM RPC JSON STA
        ${llvm_libs})

install(TARGETS dra RUNTIME DESTINATION bin)