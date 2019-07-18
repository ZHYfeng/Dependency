//
// Created by yhao on 7/18/19.
//

#ifndef INC_2018_DEPENDENCY_GENERAL_H
#define INC_2018_DEPENDENCY_GENERAL_H

#include <string>
#include <llvm/IR/BasicBlock.h>

namespace dra {



    static void outputTime(std::string s) {
        std::time_t current_time;
        current_time = std::time(nullptr);
        std::cout << std::ctime(&current_time);
        std::cout << "#time : " << s << std::endl;
    }

    llvm::BasicBlock *getRealBB(llvm::BasicBlock *b);

    llvm::BasicBlock *getFinalBB(llvm::BasicBlock *b);

    void dump_inst(llvm::Instruction *inst);
}

#endif //INC_2018_DEPENDENCY_GENERAL_H
