add_executable(a2l a2l.cpp)
target_link_libraries(a2l PUBLIC
        RPC
        ${llvm_libs})

install(TARGETS a2l RUNTIME DESTINATION bin)