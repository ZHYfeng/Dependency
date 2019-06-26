/*
 * DBasicBlock.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_BASICBLOCKALL_H_
#define LIB_DRA_BASICBLOCKALL_H_

#include <llvm/IR/BasicBlock.h>
#include <set>
#include <map>
#include <string>
#include <vector>

#include "DAInstruction.h"
#include "DInput.h"
#include "DLInstruction.h"

namespace dra {
    class DFunction;
} /* namespace dra */

namespace llvm {
    class BasicBlock;
} /* namespace llvm */

namespace dra {
    class DBasicBlock {
    public:
        DBasicBlock();

        virtual ~DBasicBlock();

        void InitIRBasicBlock(llvm::BasicBlock *b);

        void setState(CoverKind kind);

        void update(CoverKind kind, DInput *input);

        bool inferCoverBB(DInput *input, llvm::BasicBlock *b);

        void inferUncoverBB(llvm::BasicBlock *p, llvm::BasicBlock *b, int i);

        void inferSuccessors(llvm::BasicBlock *s, llvm::BasicBlock *b);

//        void inferPredecessors(llvm::BasicBlock *b);

//        void inferPredecessorsUncover(llvm::BasicBlock *b, llvm::BasicBlock *Pred);

        void infer();

        void addNewInput(DInput *i);

        bool isAsmSourceCode() const;

        void setAsmSourceCode(bool asmSourceCode);

        bool isIr() const;

        void setIr(bool ir);

        void dump();

        bool set_arrive(dra::DBasicBlock *db);

        void set_critical_condition();

        void add_critical_condition(dra::DBasicBlock *db, uint64_t condition);

    public:
        bool IR;
        bool AsmSourceCode;

        llvm::BasicBlock *basicBlock;
        DFunction *parent;
        CoverKind state;
        std::string name;
        uint64_t COVNum;
        uint64_t address;

        std::vector<DAInstruction *> InstASM;
        std::vector<DLInstruction *> InstIR;

        std::set<DInput *> input;
        DInput *lastInput;

        std::map<dra::DBasicBlock *, uint64_t > arrive;
        std::map<dra::DBasicBlock *, uint64_t > critical_condition;

        std::set<llvm::BasicBlock *> useLessPred;
    };

} /* namespace dra */

#endif /* LIB_DRA_BASICBLOCKALL_H_ */
