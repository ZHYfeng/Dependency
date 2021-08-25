set(SOURCES
        DBasicBlock.cpp
        DataManagement.cpp
        DFunction.cpp
        DAInstruction.cpp
        DLInstruction.cpp
        DModule.cpp
        DInput.cpp
        )

add_library(DMM ${SOURCES})
set_target_properties(DMM PROPERTIES LINKER_LANGUAGE CXX)
target_link_libraries(DMM PUBLIC
        ${llvm_libs}
        )