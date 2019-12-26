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
#include "llvm/IR/Instructions.h"
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

        state = CoverKind::untest;

        InstNum = 0;
        CallInstNum = 0;
        JumpInstNum = 0;
        this->RealBasicBlockNum = 0;
        BasicBlockNum = 0;

        critical_condition = false;
        this->uncovered_basicblock = false;

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
                RealBasicBlockNum++;
                BasicBlockNum++;
                Name = it.getName().str();
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
            } else {
                BasicBlockNum++;
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
//        compute_arrive();
        //	inferUseLessPred();
        //	inferUseLessPred(&f->getEntryBlock());
    }

    void DFunction::setState(CoverKind kind) {
        if (state == CoverKind::cover && kind == CoverKind::uncover) {
            std::cerr << "error DFunction kind"
                      << "\n";
        }
        state = kind;
    }

    void DFunction::update(CoverKind kind) { setState(kind); }

    bool DFunction::isObjudump() const { return Objudump; }

    void DFunction::setObjudump(bool objudump) { DFunction::Objudump = objudump; }

    bool DFunction::isAsmSourceCode() const { return AsmSourceCode; }

    void DFunction::setAsmSourceCode(bool asmSourceCode) {
        DFunction::AsmSourceCode = asmSourceCode;
    }

    bool DFunction::isIR() const { return IR; }

    void DFunction::setIR(bool ir) { DFunction::IR = ir; }

    bool DFunction::isMap() {
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

    void DFunction::dump() {

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
        std::cout << "BasicBlockNumber :" << BasicBlockNum << std::endl;
        if (this->function != nullptr) {
            //            function->dump();
        }
        std::cout << "--------------------------------------------" << std::endl;
    }

    void DFunction::inferUseLessPred(llvm::BasicBlock *b) {
        for (auto i : path) {
            std::cout << " " << i->getName().str();
            if (i == b) {
                auto bb = path.back();
                auto name = b->getName().str();
                BasicBlock[name]->useLessPred.insert(bb);
                return;
            }
        }
        path.push_back(b);
        auto *inst = b->getTerminator();
        if (inst->getNumSuccessors() == 0) {

        } else {
            for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
                inferUseLessPred(inst->getSuccessor(i));
            }
        }
        path.pop_back();
    }

    void DFunction::inferUseLessPred() {
        for (auto &it : *function) {
            order.insert(&it);
            if (it.getSinglePredecessor()) {

            } else {
                for (auto *pred : llvm::predecessors(&it)) {
                    auto name = it.getName().str();
                    if (order.find(pred) == order.end()) {
                        BasicBlock[name]->useLessPred.insert(pred);
                        std::cout << "function : " << FunctionName << std::endl;
                        std::cout << "name : " << name << std::endl;
                        std::cout << "use less : " << pred->getName().str() << std::endl;
                    }
                }
            }
        }
    }

    void DFunction::compute_arrive() {
        if (this->isIR()) {
            if (this->critical_condition) {
                return;
            } else {
                this->critical_condition = true;
                std::vector<dra::DBasicBlock *> terminator_bb;
                get_terminator(terminator_bb);

                for (auto db : terminator_bb) {
                    set_pred_successor(db);
                }
                set_critical_condition();
                return;
            }
        } else {
            std::cerr << "compute_arrive is not ir" << std::endl;
        }

    }

    void DFunction::get_terminator(std::vector<dra::DBasicBlock *> &terminator_bb) {
        for (const auto &db : this->BasicBlock) {
            for (auto di : db.second->InstIR) {
                if (auto *RI = llvm::dyn_cast<llvm::ReturnInst>(di->i)) {
                    terminator_bb.push_back(db.second);
                }
            }
        }
    }

    void DFunction::set_pred_successor(DBasicBlock *db) {
        for (auto *pred : llvm::predecessors(db->basicBlock)) {
            std::string basicblock_name = dra::getRealBB(pred)->getName().str();
            if (this->BasicBlock.find(basicblock_name) != this->BasicBlock.end()) {
                auto pred_db = this->BasicBlock[basicblock_name];
                bool new_basicblock = pred_db->set_arrive(db);
                if (new_basicblock) {
                    set_pred_successor(pred_db);
                }
            }
        }
    }

    void DFunction::set_critical_condition() {
        for (auto db : this->BasicBlock) {
            db.second->set_critical_condition();
        }
    }

    uint32_t DFunction::get_number_uncovered_instructions() {
        uint64_t uncovered_basicblock_number = 0;
        for (auto b : this->BasicBlock) {
            if (b.second->state != CoverKind::cover) {
                uncovered_basicblock_number += b.second->get_number_uncovered_instructions();
            }
        }
        return 0;
    }

    void DFunction::get_function_call(std::set<llvm::Function *> &res) {
        for (auto b : this->BasicBlock) {
            b.second->get_function_call(res);
        }
    }


    uint32_t DFunction::get_number_dominator_uncovered_instructions(llvm::BasicBlock *b) {
        uint32_t count = 0;
        count = count + this->parent->get_DB_from_bb(b)->get_number_uncovered_instructions();
        if (DT->getNode(b)->getNumChildren() != 0) {
            for (auto c : DT->getNode(b)->getChildren()) {
                auto df = this->parent->get_DB_from_bb(c->getBlock());
                count = count + df->get_number_uncovered_instructions();
                count = count + this->get_number_dominator_uncovered_instructions(df->basicBlock);
            }
        }
    }

} /* namespace dra */
