set(SOURCES
		DBasicBlock.cpp
        DataManagement.cpp
		DFunction.cpp
		DAInstruction.cpp
		DLInstruction.cpp
		DModule.cpp
        DInput.cpp
)

add_llvm_loadable_module(DRA ${SOURCES})
add_llvm_library(LLVMDRA ${SOURCES})

if(DEFINED IN_SOURCE_BUILD)
    add_dependencies(DRA intrinsics_gen)
    add_dependencies(LLVMDRA intrinsics_gen)
endif()
