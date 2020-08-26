/*
 * DFunction.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DFunction.h"

#include <iostream>
#include <llvm/IR/CFG.h>
#include <llvm/IR/Function.h>
#include <llvm/IR/InstrTypes.h>
#include <llvm/IR/Value.h>
#include <llvm/IR/Instructions.h>
#include <llvm/Analysis/CFG.h>
#include "DataManagement.h"
#include "../DCC/general.h"
#include <set>

namespace dra {

    DFunction::DFunction() {
        Objudump = false;
        AsmSourceCode = false;
        IR = false;
        Repeat = false;

        function = nullptr;
        parent = nullptr;

        state = CoverKind::outside;

        InstNum = 0;
        CallInstNum = 0;
        JumpInstNum = 0;
        NumberBasicBlock = 0;
        NumberBasicBlockReal = 0;

        uncovered_basicblock = false;

        NumberBasicBlockCovered = 0;

        DT = nullptr;
    }

    DFunction::~DFunction() = default;

    void DFunction::InitIRFunction(llvm::Function *f) {
        DFunction::function = f;
        this->DT = new llvm::DominatorTree(*f);
        std::string Name;
        DBasicBlock *b;
        int64_t no = 0;
        for (auto &it : *function) {
            if (it.hasName()) {
                NumberBasicBlock++;
                NumberBasicBlockReal++;
                Name = it.getName().str();
                if (BasicBlock.find(Name) == BasicBlock.end()) {
                    b = new DBasicBlock();
                    BasicBlock[Name] = b;
                    b->name = Name;
                    b->basicBlock = &it;
                    BasicBlock[Name]->setIr(true);
                    BasicBlock[Name]->parent = this;
                } else {
                    std::cerr << "error same basic block name" << "\n";
                }
            } else {
                NumberBasicBlock++;
                Name = std::to_string(no);

                if (BasicBlock.find(Name) == BasicBlock.end()) {
                    b = new DBasicBlock();
                    BasicBlock[Name] = b;
                    b->name = Name;
                    b->basicBlock = &it;
                    BasicBlock[Name]->setIr(true);
                    BasicBlock[Name]->parent = this;
                } else {
                    std::cerr << "error same basic block name"
                              << "\n";
                }
            }
            no++;
            BasicBlock[Name]->InitIRBasicBlock(&it);
        }
    }

    void DFunction::setState(CoverKind kind) {
        if (kind < state) {

        } else {
            state = kind;
        }
    }

    void DFunction::update(CoverKind kind) {


        if (this->state < CoverKind::untest && kind >= CoverKind::untest) {
            setState(kind);
            updateUntestedState();
        } else {
            setState(kind);
        }
    }

    void DFunction::updateUntestedState() {
        std::set<llvm::Function *> res;
        this->get_function_call(res);
        for (auto f: res) {
            if(this->parent->get_DF_from_f(f))
                this->parent->get_DF_from_f(f)->update(CoverKind::untest);
        }
    }

    bool DFunction::isObjudump() const { return Objudump; }

    void DFunction::setObjudump(bool objudump) { DFunction::Objudump = objudump; }

    bool DFunction::isAsmSourceCode() const { return AsmSourceCode; }

    void DFunction::setAsmSourceCode(bool asmSourceCode) {
        DFunction::AsmSourceCode = asmSourceCode;
    }

    bool DFunction::isIR() const { return IR; }

    void DFunction::setIR(bool ir) { DFunction::IR = ir; }

    bool DFunction::isMap() const {
        return DFunction::Objudump && DFunction::AsmSourceCode && DFunction::IR;
    }

    bool DFunction::isRepeat() const { return Repeat; }

    void DFunction::setRepeat(bool repeat) { this->Repeat = repeat; }

    void DFunction::setKind(FunctionKind kind) {
        switch (kind) {
            case dra::FunctionKind::IR: {
                setIR(true);
                break;
            }
            case dra::FunctionKind::O: {
                setObjudump(true);
                break;
            }
            case dra::FunctionKind::S: {
                setAsmSourceCode(true);
                break;
            }
            default: {
            }
        }
    }

    void DFunction::dump() const {

        std::cout << "--------------------------------------------" << std::endl;
        std::cout << "Path :" << Path << std::endl;
        std::cout << "FunctionName :" << FunctionName << std::endl;

        std::cout << "Objudump :" << Objudump << std::endl;
        std::cout << "AsmSourceCode :" << AsmSourceCode << std::endl;
        std::cout << "IR :" << IR << std::endl;
        std::cout << "repeat :" << Repeat << std::endl;
        std::cout << "CoverKind :" << state << std::endl;
        std::cout << "IRName :" << IRName << std::endl;
        std::cout << "Address :" << Address << std::endl;
        std::cout << "InstNum :" << InstNum << std::endl;
        std::cout << "CallInstNum :" << CallInstNum << std::endl;
        std::cout << "JumpInstNum :" << JumpInstNum << std::endl;
        std::cout << "NumberBasicBlock :" << NumberBasicBlock << std::endl;
        if (this->function != nullptr) {
            //            function->dump();
        }
        std::cout << "--------------------------------------------" << std::endl;
    }

    uint32_t DFunction::get_number_uncovered_instructions(std::map<std::string, dra::DBasicBlock *> &res) {
        uint64_t uncovered_basicblock_number = 0;
        for (const auto &b : this->BasicBlock) {
            if (b.second->state != CoverKind::cover && b.second->basicBlock != nullptr &&
                b.second->basicBlock->hasName()) {
                uncovered_basicblock_number += b.second->get_number_uncovered_instructions();
                res[this->FunctionName + b.first] = b.second;
            }
        }
        return 0;
    }

    void DFunction::get_function_call(std::set<llvm::Function *> &res) {
        for (auto b : this->BasicBlock) {
            b.second->get_function_call(res);
        }
    }


    uint32_t DFunction::get_dominator_uncovered_instructions(llvm::BasicBlock *b,
                                                             std::map<std::string, dra::DBasicBlock *> &res) {
        uint32_t count = 0;
        if (b->hasName()) {
            std::string Name = b->getName().str();
            if (BasicBlock.find(Name) != BasicBlock.end()) {
                if (this->BasicBlock[Name]->state != CoverKind::cover) {
                    count = count + this->BasicBlock[Name]->get_number_uncovered_instructions();
                    res[this->FunctionName + Name] = this->BasicBlock[Name];
                }
            }
        }

        for (auto c : DT->getNode(dra::getFinalBB(b))->getChildren()) {
            if (c->getBlock()->hasName()) {
                std::string Name = c->getBlock()->getName().str();
                if (BasicBlock.find(Name) != BasicBlock.end()) {
                    if (this->BasicBlock[Name]->state != CoverKind::cover) {
                        count = count + this->get_dominator_uncovered_instructions(c->getBlock(), res);
                    }
                }
            }
        }
        return count;
    }

    void DFunction::add_number_basic_block_covered() {
        this->NumberBasicBlockCovered++;
        this->parent->add_number_basic_block_covered();
    }

    uint32_t
    DFunction::get_uncovered_instructions(llvm::BasicBlock *b, std::map<std::string, dra::DBasicBlock *> &res) {
        uint32_t count = 0;
        if (b->hasName()) {
            std::string Name = b->getName().str();
            if (BasicBlock.find(Name) != BasicBlock.end()) {
                if (this->BasicBlock[Name]->state != CoverKind::cover) {
                    count = count + this->BasicBlock[Name]->get_number_uncovered_instructions();
                    res[this->FunctionName + Name] = this->BasicBlock[Name];
                }
            }
        }

        std::set<llvm::Function *> uncovered_function;
        std::set<llvm::Function *> new_uncovered_functions;
        uncovered_function.insert(this->function);

        for (const auto &bb: this->BasicBlock) {
            if (bb.second->basicBlock != nullptr && bb.second->basicBlock != nullptr &&
                bb.second->basicBlock->hasName()) {
                if (llvm::isPotentiallyReachable(b, bb.second->basicBlock, this->DT)) {
                    if (bb.second->state != CoverKind::cover) {
                        count = count + bb.second->get_number_uncovered_instructions();
                        res[this->FunctionName + bb.first] = bb.second;
                        bb.second->get_function_call(new_uncovered_functions);
                    }
                }
            }
        }

        while (!new_uncovered_functions.empty()) {
            std::set<llvm::Function *> temp;
            for (auto f : new_uncovered_functions) {
                uncovered_function.insert(f);
                DFunction *df = this->parent->get_DF_from_f(f);
                if (df != nullptr) {
                    if (df->state != CoverKind::cover) {
                        count += df->get_number_uncovered_instructions(res);
                        df->get_function_call(temp);
                    }
                }
            }
            new_uncovered_functions.clear();
            for (auto f : temp) {
                if (uncovered_function.find(f) == uncovered_function.end()) {
                    new_uncovered_functions.insert(f);
                    uncovered_function.insert(f);
                }
            }
        }

        return count;
    }

} /* namespace dra */
