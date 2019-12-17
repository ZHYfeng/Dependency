/*
 * DAInstruction.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_INSTRUCTIONASM_H_
#define LIB_DRA_INSTRUCTIONASM_H_

#include <string>

#include "DLInstruction.h"

namespace dra {
    class DInput;
} /* namespace dra */

namespace dra {
    class DBasicBlock;

    class DAInstruction {
    public:
        DAInstruction();

        virtual ~DAInstruction();

        void setState(CoverKind kind);

        void update(CoverKind kind, DInput *input);

        void setAddr(std::string addr);

    public:
        CoverKind state;

        std::string SInst;
        std::string BasicBlockName;
        DBasicBlock *parent;

        std::string OInst;
        std::string Address;
        unsigned long long int address;
    };

} /* namespace dra */

#endif /* LIB_DRA_INSTRUCTIONASM_H_ */
