//
// Created by yhao on 7/21/19.
//

#include <string>
#include <signal.h>
#include <iostream>
#include <llvm/IR/BasicBlock.h>
#include "llvm/IR/CFG.h"
#include <llvm/IR/DebugLoc.h>
#include <llvm/IR/DebugInfoMetadata.h>
#include "general.h"

namespace dra {

#ifndef _STACKTRACE_H_
#define _STACKTRACE_H_

#include <stdio.h>
#include <stdlib.h>
#include <execinfo.h>
#include <cxxabi.h>

/** Print a demangled stack backtrace of the caller function to FILE* out. */
    static inline void print_stacktrace(FILE *out = stderr, unsigned int max_frames = 63) {
        fprintf(out, "stack trace:\n");

        // storage array for stack trace address data
        void *addrlist[max_frames + 1];

        // retrieve current stack addresses
        int addrlen = backtrace(addrlist, sizeof(addrlist) / sizeof(void *));

        if (addrlen == 0) {
            fprintf(out, "  <empty, possibly corrupt>\n");
            return;
        }

        // resolve addresses into strings containing "filename(function+address)",
        // this array must be free()-ed
        char **symbollist = backtrace_symbols(addrlist, addrlen);

        // allocate string which will be filled with the demangled function name
        size_t funcnamesize = 256;
        char *funcname = (char *) malloc(funcnamesize);

        // iterate over the returned symbol lines. skip the first, it is the
        // address of this function.
        for (int i = 1; i < addrlen; i++) {
            char *begin_name = 0, *begin_offset = 0, *end_offset = 0;

            // find parentheses and +address offset surrounding the mangled name:
            // ./module(function+0x15c) [0x8048a6d]
            for (char *p = symbollist[i]; *p; ++p) {
                if (*p == '(')
                    begin_name = p;
                else if (*p == '+')
                    begin_offset = p;
                else if (*p == ')' && begin_offset) {
                    end_offset = p;
                    break;
                }
            }

            if (begin_name && begin_offset && end_offset && begin_name < begin_offset) {
                *begin_name++ = '\0';
                *begin_offset++ = '\0';
                *end_offset = '\0';

                // mangled name is now in [begin_name, begin_offset) and caller
                // offset in [begin_offset, end_offset). now apply
                // __cxa_demangle():

                int status;
                char *ret = abi::__cxa_demangle(begin_name,
                                                funcname, &funcnamesize, &status);
                if (status == 0) {
                    funcname = ret; // use possibly realloc()-ed string
                    fprintf(out, "  %s : %s+%s\n",
                            symbollist[i], funcname, begin_offset);
                } else {
                    // demangling failed. Output function name as a C function with
                    // no arguments.
                    fprintf(out, "  %s : %s()+%s\n",
                            symbollist[i], begin_name, begin_offset);
                }
            } else {
                // couldn't parse the line? print the whole line.
                fprintf(out, "  %s\n", symbollist[i]);
            }
        }

        free(funcname);
        free(symbollist);
    }

#endif // _STACKTRACE_H_

    void outputTime(std::string s) {
        time_t current_time;
        current_time = time(nullptr);
        std::cout << ctime(&current_time);
        std::cout << "#time : " << s << std::endl;
    }

    void handler(int nSignum, siginfo_t *si, void *vcontext) {
        std::cout << "Segmentation fault" << std::endl;

        ucontext_t *context = (ucontext_t *) vcontext;
        context->uc_mcontext.gregs[REG_RIP]++;

        print_stacktrace();

        exit(-1);
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

    void deal_sig() {

        struct sigaction action;
        memset(&action, 0, sizeof(struct sigaction));
        action.sa_flags = SA_SIGINFO;
        action.sa_sigaction = dra::handler;
        sigaction(SIGSEGV, &action, NULL);

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

    }
}