/*
 * dra.cpp
 *
 *  Created on: Nov 28, 2018
 *      Author: yhao
 */

#include <llvm/Support/CommandLine.h>
#include <iostream>
#include <grpcpp/grpcpp.h>

#include "../../lib/DCC/DependencyControlCenter.h"
#include "../../lib/DRA/DataManagement.h"

llvm::cl::opt<std::string> objdump("objdump", llvm::cl::desc("The path of objdump."), llvm::cl::init("./vmlinux.objdump"));
llvm::cl::opt<std::string> AssemblySourceCode("asm", llvm::cl::desc("The path of assembly source code."), llvm::cl::init("./build-in.s"));
llvm::cl::opt<std::string> InputFilename(llvm::cl::Positional, llvm::cl::desc("<input bitcode>"), llvm::cl::init("./built-in.bc"));
llvm::cl::opt<std::string> coverfile("coverfile", llvm::cl::desc("The path of cover file."), llvm::cl::init("./cover.txt"));
llvm::cl::opt<std::string> vmOffsets("vmOffsets", llvm::cl::desc("The path of vmOffsets.txt."), llvm::cl::init("./vmOffsets.txt"));
//The file holding the serialized static analysis results.
llvm::cl::opt<std::string> staticRes("staticRes", llvm::cl::desc("The path of serialized static analysis results."), llvm::cl::init("./taint_info_serialize"));

int main(int argc, char **argv) {

    llvm::cl::ParseCommandLineOptions(argc, argv, "dra\n");
#if DEBUG
    std::cout << "AssemblySourceCode : " << AssemblySourceCode << std::endl;
    std::cout << "objdump : " << objdump << std::endl;
    std::cout << "InputFilename : " << InputFilename << std::endl;
    std::cout << "staticRes : " << staticRes << std::endl;
#endif

    dra::DependencyControlCenter dcc;

    dcc.init(objdump, AssemblySourceCode, InputFilename, staticRes);
    dcc.run();
    return 0;
}
