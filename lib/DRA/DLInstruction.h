/*
 * DLInstruction.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_INSTRUCTIONLLVM_H_
#define LIB_DRA_INSTRUCTIONLLVM_H_

#include <string>
#include <llvm/IR/Instruction.h>

namespace dra {
    enum Kind {
        other, untest, cover, uncover,
    };

    class DBasicBlock;

    class DLInstruction {
    public:
        DLInstruction();

        virtual ~DLInstruction();

        void setState(Kind kind);

        void update(Kind kind);

    public:
        llvm::Instruction *i;
        Kind state;
        DBasicBlock *parent;

        std::string FileName;
        unsigned int Line;


    };

} /* namespace dra */

#endif /* LIB_DRA_INSTRUCTIONLLVM_H_ */
