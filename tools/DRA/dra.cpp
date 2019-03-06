/*
 * dra.cpp
 *
 *  Created on: Nov 28, 2018
 *      Author: yhao
 */

#include <llvm/Support/CommandLine.h>
#include <iostream>

#include "../../lib/DRA/DataManagement.h"

llvm::cl::opt<std::string> coverfile("coverfile", llvm::cl::desc("The path of cover file."),
                                     llvm::cl::init("./0-cover.html"));

llvm::cl::opt<std::string> objdump("objdump", llvm::cl::init("./vmlinux.objdump"),
                                   llvm::cl::desc("The path of objdump."));

llvm::cl::opt<std::string> AssemblySourceCode("asm", llvm::cl::desc("The path of assembly source code."),
                                              llvm::cl::init("./all.s"));

llvm::cl::opt<std::string> InputFilename(llvm::cl::Positional, llvm::cl::desc("<input bitcode>"), llvm::cl::init("-"));

int main(int argc, char **argv) {

    llvm::cl::ParseCommandLineOptions(argc, argv, "dra\n");
#if DEBUG
    std::cout << "AssemblySourceCode : " << AssemblySourceCode << std::endl;
    std::cout << "objdump : " << objdump << std::endl;
    std::cout << "InputFilename : " << InputFilename << std::endl;
#endif

    dra::DataManagement *MI = new dra::DataManagement();
//	MI->GetInformationFromCoverFile(coverfile);
    MI->initializeModule(objdump, AssemblySourceCode, InputFilename);
//	MI->a->Analysis(AssemblySourceCode, objdump);
//	MI->MapBBfromStoBC();

    return 0;
}
