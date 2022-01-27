add_executable(a2i a2i.cpp)
target_link_libraries(a2i PUBLIC
        DCC DMM RPC JSON STA
        ${llvm_libs})

install(TARGETS a2i RUNTIME DESTINATION bin)