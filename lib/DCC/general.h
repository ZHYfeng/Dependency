//
// Created by yhao on 7/18/19.
//

#ifndef INC_2018_DEPENDENCY_GENERAL_H
#define INC_2018_DEPENDENCY_GENERAL_H

#include <string>
#include <csetjmp>
#include <csignal>
#include <iostream>
#include <llvm/IR/BasicBlock.h>

#define DEBUG 0
#define DEBUG_ERR 0

namespace dra {

    void outputTime(std::string s);

    llvm::BasicBlock *getRealBB(llvm::BasicBlock *b);

    llvm::BasicBlock *getFinalBB(llvm::BasicBlock *b);

    std::string getFileName(llvm::Function *f);

    std::string getFunctionName(llvm::Function *f);

    void dump_inst(llvm::Instruction *inst);
}

#endif //INC_2018_DEPENDENCY_GENERAL_H
