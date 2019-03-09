/*
 * DLInstruction.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DLInstruction.h"

#include <iostream>

#include "DBasicBlock.h"


namespace dra {

    DLInstruction::DLInstruction() :
            i(nullptr), parent(nullptr), Line(0) {
        state = Kind::other;

    }

    DLInstruction::~DLInstruction() = default;

    void DLInstruction::setState(Kind kind) {
        if (state == Kind::cover && kind == Kind::uncover) {
            std::cerr << "error InstIR kind" << "\n";
        }
        state = kind;
    }

    void DLInstruction::update(Kind kind) {
        setState(kind);
        parent->update(kind);
    }

} /* namespace dra */
