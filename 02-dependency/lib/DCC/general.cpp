//
// Created by yhao on 7/21/19.
//

#include <string>
#include <iostream>
#include <llvm/IR/BasicBlock.h>
#include "llvm/IR/CFG.h"
#include <llvm/IR/DebugLoc.h>
#include <llvm/IR/DebugInfoMetadata.h>
#include "general.h"

namespace dra {

    void outputTime(std::string s) {
#if DEBUG
        time_t current_time;
        current_time = time(nullptr);
        std::cout << ctime(&current_time);
        std::cout << "#time : " << s << std::endl;
#endif
    }

    llvm::BasicBlock *getRealBB(llvm::BasicBlock *b) {
        llvm::BasicBlock *rb;
        if (b->hasName()) {
            rb = b;
        } else {
            for (auto *Pred : llvm::predecessors(b)) {
                rb = getRealBB(Pred);
                break;
            }
        }
        return rb;
    }

    llvm::BasicBlock *getFinalBB(llvm::BasicBlock *b) {
        auto *inst = b->getTerminator();
        for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
            std::string name = inst->getSuccessor(i)->getName().str();
            if (inst->getSuccessor(i)->hasName()) {
            } else {
                return getFinalBB(inst->getSuccessor(i));
            }
        }
        return b;
    }

    std::string getFileName(llvm::Function *f) {
        llvm::SmallVector<std::pair<unsigned, llvm::MDNode *>, 4> MDs;
        f->getAllMetadata(MDs);
        for (auto &MD : MDs) {
            if (llvm::MDNode *N = MD.second) {
                if (auto *SP = llvm::dyn_cast<llvm::DISubprogram>(N)) {
                    std::string Path = SP->getFilename().str();
                    return Path;
                }
            }
        }
        return "";
    }

    std::string getFunctionName(llvm::Function *f) {
        std::string name = f->getName().str();
        std::string FunctionName;
        if (name.find('.') < name.size()) {
            FunctionName = name.substr(0, name.find('.'));
        } else {
            FunctionName = name;
        }
        return FunctionName;
    }

    void dump_inst(llvm::Instruction *inst) {
        if (inst != nullptr) {

        } else {
            return;
        }
        auto b = inst->getParent();
        auto f = b->getParent();

        std::string Path = dra::getFileName(f);
        std::string FunctionName = dra::getFunctionName(f);
        std::cout << Path << " : ";
        std::cout << FunctionName << " : ";

        const llvm::DebugLoc &debugInfo = inst->getDebugLoc();
        int line = debugInfo->getLine();
        int column = debugInfo->getColumn();
        std::cout << std::dec << line << " : ";
        std::cout << column << " : ";


        std::string BasicBlockName = getRealBB(b)->getName();
        std::cout << BasicBlockName << " : ";

//        std::string directory = debugInfo->getDirectory().str();
//        std::string filePath = debugInfo->getFilename().str();

        std::cout << std::endl;

        //https://elixir.bootlin.com/linux/v4.16/source/drivers/tty/tty_io.c#L2358
        std::cout << "https://elixir.bootlin.com/linux/v4.16/source/" << Path << "#L" << line << std::endl;

    }
}