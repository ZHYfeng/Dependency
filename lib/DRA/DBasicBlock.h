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
#include <string>
#include <vector>

#include "DAInstruction.h"
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

        void update(CoverKind kind);

        void setOtherBBState(llvm::BasicBlock *b, CoverKind kind);

        void inferSuccessors(llvm::BasicBlock *b);

        void inferPredecessors(llvm::BasicBlock *b);

        void inferPredecessorsUncover(llvm::BasicBlock *b, llvm::BasicBlock *Pred);

        void infer();

        bool isAsmSourceCode() const;

        void setAsmSourceCode(bool asmSourceCode);

        bool isIr() const;

        void setIr(bool ir);

    public:
        bool IR;
        bool AsmSourceCode;

        llvm::BasicBlock *basicBlock;
        std::set<llvm::BasicBlock *> useLessPred;
        DFunction *parent;
        CoverKind state;
        std::string name;
        unsigned int COVNum;

        std::vector<DAInstruction *> InstASM;
        std::vector<DLInstruction *> InstIR;
    };

} /* namespace dra */

#endif /* LIB_DRA_BASICBLOCKALL_H_ */
