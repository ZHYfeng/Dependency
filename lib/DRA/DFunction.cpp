/*
 * DFunction.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DFunction.h"

#include <llvm/IR/Function.h>
#include <iostream>

namespace dra {


    DFunction::DFunction() {

    }

    DFunction::~DFunction() = default;

    void DFunction::InitIRFunction(llvm::Function *f) {
        DFunction::function = f;
        for (auto &it : *function) {

            std::string Name = it.getName().str();
            DBasicBlock *b;
            b = new DBasicBlock();
            BasicBlock[Name] = b;

            b->InitIRBasicBlock(&it);
        }
    }

    void DFunction::setState(Kind kind) {
        if (state == Kind::cover && kind == Kind::uncover) {
            std::cerr << "error BasicBlock kind" << "\n";
        }
        state = kind;
    }

    void DFunction::update(Kind kind) {
        setState(kind);
    }

    bool DFunction::isObjudump() const {
        return Objudump;
    }

    void DFunction::setObjudump(bool Objudump) {
        DFunction::Objudump = Objudump;
    }

    bool DFunction::isAsmSourceCode() const {
        return AsmSourceCode;
    }

    void DFunction::setAsmSourceCode(bool AsmSourceCode) {
        DFunction::AsmSourceCode = AsmSourceCode;
    }

    bool DFunction::isIR() const {
        return IR;
    }

    void DFunction::setIR(bool IR) {
        DFunction::IR = IR;
    }

    bool DFunction::isMap() {
        return DFunction::Objudump && DFunction::AsmSourceCode && DFunction::IR;
    }

} /* namespace dra */
