/*
 * DAInstruction.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DAInstruction.h"

#include <iostream>

#include "DBasicBlock.h"

namespace dra {

    DAInstruction::DAInstruction() {
        state = CoverKind::untest;
        parent = nullptr;
        address = 0;

    }

    DAInstruction::~DAInstruction() = default;

    void DAInstruction::setState(CoverKind kind) {
        if (state == CoverKind::cover && kind == CoverKind::uncover) {
            std::cerr << "error InstIR kind" << "\n";
        }
        state = kind;
    }

    void DAInstruction::update(CoverKind kind) {
        setState(kind);
        parent->update(kind);
    }

    void DAInstruction::setAddr(std::string addr) {
        this->Address = addr;
        this->address = std::stoul(addr, nullptr, 16);

    }

} /* namespace dra */
