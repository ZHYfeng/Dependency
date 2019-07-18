/*
 * DBasicBlock.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DBasicBlock.h"

#include <llvm/IR/Instructions.h>
#include "llvm/IR/CFG.h"
#include <iostream>

#include "DFunction.h"
#include "DataManagement.h"

namespace dra {

    DBasicBlock::DBasicBlock() {
        IR = false;
        AsmSourceCode = false;
        basicBlock = nullptr;
        parent = nullptr;
        state = CoverKind::untest;
        tracr_num = 0;
        this->lastInput = nullptr;
    }

    DBasicBlock::~DBasicBlock() = default;

    void DBasicBlock::InitIRBasicBlock(llvm::BasicBlock *b) {
        for (auto &it : *b) {
            DLInstruction *i;

            i = new DLInstruction();
            InstIR.push_back(i);

            i->parent = this;
            i->i = (&it);
        }
    }

    void DBasicBlock::setState(CoverKind kind) {

        if (kind == CoverKind::cover) {
            state = kind;
        } else if (kind == CoverKind::uncover) {
            if (state == CoverKind::cover) {
                std::cout << "cover to uncover basic block name : " << name << std::endl;
                parent->dump();
                exit(0);
            }
        }
    }

    void DBasicBlock::update(CoverKind kind, DInput *input) {
        setState(kind);
        for (auto it : InstIR) {
            it->setState(kind);
        }
        for (auto it : InstASM) {
            it->setState(kind);
        }
        if (this->parent != nullptr) {
            parent->update(kind);
        } else {
//            std::cerr << "DBasicBlock update parent == nullptr" << "\n";
//            this->dump();
        }

        if (this->basicBlock != nullptr) {
            if (inferCoverBB(input, this->basicBlock)) {

            }
            infer();
        } else {
            std::cerr << "DBasicBlock update basicBlock == nullptr : " << this->trace_pc_address << "\n";
        }

    }

    bool DBasicBlock::isAsmSourceCode() const {
        return AsmSourceCode;
    }

    void DBasicBlock::setAsmSourceCode(bool asmSourceCode) {
        AsmSourceCode = asmSourceCode;
    }

    bool DBasicBlock::isIr() const {
        return IR;
    }

    void DBasicBlock::setIr(bool ir) {
        IR = ir;
    }

    bool DBasicBlock::inferCoverBB(DInput *input, llvm::BasicBlock *b) {
        DBasicBlock *Db;
        std::string name = b->getName().str();
        if (parent->BasicBlock.find(name) != parent->BasicBlock.end()) {
            Db = parent->BasicBlock[name];
            if (Db->state == CoverKind::untest || Db->state == CoverKind::uncover) {
                Db->setState(CoverKind::cover);
                Db->input.clear();
                this->addNewInput(input);
                return true;
            } else if (Db->state == CoverKind::cover) {
                this->addNewInput(input);
            }
        } else {
            std::cout << "inferCoverBB not find basic block name : " << name << std::endl;
            parent->dump();
            exit(0);
        }
        return false;
    }

    void DBasicBlock::inferUncoverBB(llvm::BasicBlock *p, llvm::BasicBlock *b, int i) {
        DBasicBlock *Dp;
        DBasicBlock *Db;
        std::string pname = p->getName().str();
        std::string name = b->getName().str();
        if (parent->BasicBlock.find(name) != parent->BasicBlock.end()) {
            Dp = parent->BasicBlock[pname];
            Db = parent->BasicBlock[name];
            if (parent->BasicBlock.find(pname) != parent->BasicBlock.end()) {
                DInput *dInput = parent->BasicBlock[pname]->lastInput;
                if (Db->state == CoverKind::untest) {
                    Db->setState(CoverKind::uncover);
                    Db->addNewInput(dInput);
                    dInput->addUncoveredAddress(Db->trace_pc_address, Dp->trace_pc_address, i);
                } else if (Db->state == CoverKind::uncover) {
                    Db->addNewInput(dInput);
                    dInput->addUncoveredAddress(Db->trace_pc_address, Dp->trace_pc_address, i);
                } else if (Db->state == CoverKind::cover) {

                }
#if DEBUGINPUT
                if (Db->state == CoverKind::uncover) {
                    std::cout << "-------uncover basic block-----------------" << std::endl;
                    Db->dump();
                }
#endif
            } else {
                std::cout << "inferUncoverBB not find basic block pname : " << name << std::endl;
            }

        } else {
            std::cout << "inferUncoverBB not find basic block name : " << name << std::endl;
            parent->dump();
            exit(0);
        }
    }

    void DBasicBlock::inferSuccessors(llvm::BasicBlock *s, llvm::BasicBlock *b) {
        auto *inst = b->getTerminator();
        if (inst->getNumSuccessors() == 1) {
//		setOtherBBState(addr2line, inst->getSuccessor(0), CoverKind::cover);
//		inferSuccessors(inst->getSuccessor(0));
        } else {
            for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
                if (inst->getSuccessor(i)->hasName()) {
                    inferUncoverBB(s, inst->getSuccessor(i), i);
                } else {
                    inferSuccessors(s, inst->getSuccessor(i));
                }
            }
        }
    }

//    void DBasicBlock::inferPredecessors(llvm::BasicBlock *b) {
//        std::string name = b->getName().str();
//        auto input = this->parent->BasicBlock[name]->lastInput;
//        if (b->getSinglePredecessor()) {
//            inferCoverBB(input, b->getSinglePredecessor());
//            inferPredecessors(b->getSinglePredecessor());
//            inferPredecessorsUncover(b, b->getSinglePredecessor());
//        } else if (useLessPred.size() > 0) {
//            int num = 0;
//            llvm::BasicBlock *pb;
//            for (auto *Pred : llvm::predecessors(b)) {
//                if (useLessPred.find(Pred) == useLessPred.end()) {
//                    pb = Pred;
//                    num++;
//                }
//            }
//            if (num == 1) {
//                inferCoverBB(input, pb);
//                inferPredecessors(pb);
//                inferPredecessorsUncover(b, pb);
//            }
//        } else if (!b->hasName()) {
//            for (auto *Pred : llvm::predecessors(b)) {
//                inferCoverBB(input, b->getSinglePredecessor());
//                inferPredecessors(b->getSinglePredecessor());
//                inferPredecessorsUncover(b, b->getSinglePredecessor());
//            }
//        } else {
//
//        }
//
//    }

//    void DBasicBlock::inferPredecessorsUncover(llvm::BasicBlock *b, llvm::BasicBlock *Pred) {
//        auto *inst = Pred->getTerminator();
//        if (inst->getNumSuccessors() == 1) {
//
//        } else {
//            for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
//                if (inst->getSuccessor(i) != b && inst->getSuccessor(i)->hasName()) {
//                    inferUncoverBB(b, inst->getSuccessor(i));
//                }
//            }
//        }
//    }

    void DBasicBlock::infer() {
        if (this->state == CoverKind::cover) {
            inferSuccessors(this->basicBlock, this->basicBlock);
//		    inferPredecessors(this->basicBlock);
        }

    }

    void DBasicBlock::addNewInput(DInput *i) {
        this->lastInput = i;
        this->input[i] = i->idx;
    }

    void DBasicBlock::dump() {

        std::cout << "--------------------------------------------" << std::endl;
        if (parent != nullptr) {
            std::cout << "Path :" << parent->Path << std::endl;
            std::cout << "FunctionName :" << parent->FunctionName << std::endl;
        }
        std::cout << "basicblock name :" << name << std::endl;
        std::cout << "AsmSourceCode :" << AsmSourceCode << std::endl;
        std::cout << "IR :" << IR << std::endl;
        std::cout << "CoverKind :" << state << std::endl;
        std::cout << "trace_pc_address :" << std::hex << trace_pc_address << std::endl;
//        basicBlock->dump();
        for (auto i : this->input) {
            std::cout << "input : " << i.second << " : " << i.first->sig << std::endl;
            std::cout << i.first->program;
        }
        std::cout << "--------------------------------------------" << std::endl;

    }

    // not work if there is a switch with more than 64 cases.
    bool DBasicBlock::set_arrive(dra::DBasicBlock *db) {
        bool res = false;
        uint64_t Num;
        auto *inst = dra::getFinalBB(this->basicBlock)->getTerminator();
        for (uint64_t i = 0, end = inst->getNumSuccessors(); i < end; i++) {
            if (inst->getSuccessor(i) == db->basicBlock) {
                Num = i;
                if (this->arrive.find(db) != this->arrive.end()) {
                    if ((this->arrive[db] & 1 << i) > 0) {

                    } else {
                        this->arrive[db] |= 1 << i;
                        res = true;
                    }
                } else {
                    this->arrive[db] = 1 << i;
                    res = true;
                }
            }
        }

        for (auto bb : db->arrive) {
            if (this->arrive.find(bb.first) != this->arrive.end()) {
                if ((this->arrive[bb.first] & 1 << Num) > 0) {

                } else {
                    this->arrive[bb.first] |= 1 << Num;
                    res = true;
                }
            } else {
                this->arrive[bb.first] = 1 << Num;
                res = true;
            }
        }
        return res;
    }

    void DBasicBlock::set_critical_condition() {

        if(this->basicBlock == nullptr) {

        } else {
            auto *fb = dra::getFinalBB(this->basicBlock);
            auto *inst = fb->getTerminator();
            auto successor_num = inst->getNumSuccessors();
            for (auto bb : this->arrive) {
                if (bb.second == ((1 << successor_num) - 1)) {

                } else {
                    bb.first->add_critical_condition(this, bb.second);
                }
            }
        }
    }

    void DBasicBlock::add_critical_condition(dra::DBasicBlock *db, uint64_t condition) {
        Condition *c = new Condition();
        c->set_condition_address(db->trace_pc_address);
        c->set_successor(condition);
        auto *inst = dra::getFinalBB(db->basicBlock)->getTerminator();
        for (uint64_t i = 0, end = inst->getNumSuccessors(); i < end; i++) {
            auto temp = this->get_DB_from_bb(inst->getSuccessor(i));
            if (condition && 1 << i) {
                c->add_right_branch_address(temp->trace_pc_address);
            } else {
                c->add_wrong_branch_address(temp->trace_pc_address);
            }
        }
        this->critical_condition[db] = c;
    }

    DBasicBlock *DBasicBlock::get_DB_from_bb(llvm::BasicBlock *b) {
        llvm::BasicBlock *bb = dra::getRealBB(b);
        std::string bbname = bb->getName().str();
        DBasicBlock *db = this->parent->BasicBlock[bbname];
        return db;
    }

} /* namespace dra */
