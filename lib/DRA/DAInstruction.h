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

        void setState(CoverKind kind);
        void update(CoverKind kind);

    public:
        CoverKind state;

        std::string SInst;
        std::string BasicBlockName;
        DBasicBlock *parent;

        std::string OInst;
        std::string Address;
    };

} /* namespace dra */

#endif /* LIB_DRA_INSTRUCTIONASM_H_ */
