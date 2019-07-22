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

namespace dra {

    static sigjmp_buf escapeCallJmpBuf;

    extern "C" {
    static void sigsegv_handler(int signal, siginfo_t *info, void *context) {
        siglongjmp(escapeCallJmpBuf, 1);
    }
    }

    void outputTime(std::string s);

    void handler(int nSignum, siginfo_t *si, void *vcontext);

    void deal_sig();

    llvm::BasicBlock *getRealBB(llvm::BasicBlock *b);

    llvm::BasicBlock *getFinalBB(llvm::BasicBlock *b);

    std::string getFileName(llvm::Function *f);

    std::string getFunctionName(llvm::Function *f);

    void dump_inst(llvm::Instruction *inst);
}

#endif //INC_2018_DEPENDENCY_GENERAL_H
