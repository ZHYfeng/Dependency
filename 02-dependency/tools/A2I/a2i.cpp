/*
 * dra.cpp
 *
 *  Created on: Nov 28, 2018
 *      Author: yhao
 */
#include <llvm/Support/Signals.h>
#include <llvm/ADT/StringRef.h>
#include <llvm/Support/CommandLine.h>
#include <iostream>

#include "../../lib/DCC/DependencyControlCenter.h"

llvm::cl::opt<std::string> obj_dump("objdump", llvm::cl::desc("The obj dump file."), llvm::cl::init("./vmlinux.objdump"));
llvm::cl::opt<std::string> assembly("asm", llvm::cl::desc("The assembly source code."), llvm::cl::init("./build-in.s"));
llvm::cl::opt<std::string> bit_code("bc", llvm::cl::desc("The bit code."), llvm::cl::init("./built-in.bc"));

llvm::cl::opt<std::string> config(llvm::cl::Positional, llvm::cl::desc("The dra config file"), llvm::cl::init("dra.json"));
llvm::cl::opt<std::string> file("file", llvm::cl::desc("The file of uncovered address."),
                                llvm::cl::init("./not_covered.txt"));
int main(int argc, char **argv)
{
    llvm::sys::PrintStackTraceOnErrorSignal(argv[0]);
    llvm::cl::ParseCommandLineOptions(argc, argv, "a2i\n");
#if DEBUG
    std::cout << "AssemblySourceCode : " << AssemblySourceCode << std::endl;
    std::cout << "objdump : " << objdump << std::endl;
    std::cout << "InputFilename : " << InputFilename << std::endl;
    std::cout << "staticRes : " << staticRes << std::endl;
#endif

    auto *dcc = new dra::DependencyControlCenter();

    dcc->init(obj_dump, assembly, bit_code, config);
    dcc->check_uncovered_addresses_depednency(file);
    return 0;
}
