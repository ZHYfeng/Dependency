set(SOURCES
	BasicBlockAll.cpp
	DataASM.cpp
	DataManagement.cpp
	FunctionAll.cpp
	InformationOfLine.cpp
	InformationOfSourceCode.cpp
	InstructionASM.cpp
	InstructionLLVM.cpp
	ModuleAll.cpp
)

add_llvm_loadable_module(DRA ${SOURCES})
add_llvm_library(LLVMDRA ${SOURCES})

if(DEFINED IN_SOURCE_BUILD)
    add_dependencies(DRA intrinsics_gen)
    add_dependencies(LLVMDRA intrinsics_gen)
endif()
