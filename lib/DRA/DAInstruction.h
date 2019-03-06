/*
 * DAInstruction.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_INSTRUCTIONASM_H_
#define LIB_DRA_INSTRUCTIONASM_H_

class DBasicBlock;

#include <string>

#include "DLInstruction.h"

namespace dra {

    class DAInstruction {
    public:
        DAInstruction();

        virtual ~DAInstruction();

    public:
        Kind state;
        std::string Inst;
        std::string BasicBlockName;
        DBasicBlock *parent;
        std::string Address;
    };

} /* namespace dra */

#endif /* LIB_DRA_INSTRUCTIONASM_H_ */
