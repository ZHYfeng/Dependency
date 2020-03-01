/*
 * DBasicBlock.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_DBASICBLOCK_CPP_
#define LIB_DRA_DBASICBLOCK_CPP_

#include "DBasicBlock.h"
#include <llvm/IR/IntrinsicInst.h>
#include <llvm/IR/Instructions.h>
#include <llvm/IR/CallSite.h>
#include <llvm/Analysis/PostDominators.h>
#include <iostream>

#include "DFunction.h"
#include "DataManagement.h"
#include "../DCC/general.h"

namespace dra {

    DBasicBlock::DBasicBlock() {
        IR = false;
        AsmSourceCode = false;
        basicBlock = nullptr;
        parent = nullptr;
        state = CoverKind::untest;
        tracr_num = 0;
        this->lastInput = nullptr;
        this->number_instructions = 1;
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

        for (auto &it : *b) {
            if (llvm::isa<llvm::PHINode>(it) || llvm::isa<llvm::DbgInfoIntrinsic>(it))
                continue;
            if (auto *II = llvm::dyn_cast<llvm::IntrinsicInst>(&it)) {

                if (II->getIntrinsicID() == llvm::Intrinsic::lifetime_start ||
                    II->getIntrinsicID() == llvm::Intrinsic::lifetime_end) {
                    continue;
                }
            }
            this->number_instructions++;
        }

        for (auto temp = b->getNextNode(); temp != nullptr && !temp->hasName(); temp = temp->getNextNode()) {
            for (auto &it : *temp) {
                if (llvm::isa<llvm::PHINode>(it) || llvm::isa<llvm::DbgInfoIntrinsic>(it))
                    continue;
                if (auto *II = llvm::dyn_cast<llvm::IntrinsicInst>(&it)) {

                    if (II->getIntrinsicID() == llvm::Intrinsic::lifetime_start ||
                        II->getIntrinsicID() == llvm::Intrinsic::lifetime_end) {
                        continue;
                    }
                }
                this->number_instructions++;
            }
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

    void DBasicBlock::update(CoverKind kind, DInput *dInput) {
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
            if (inferCoverBB(dInput, this->basicBlock)) {

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

    bool DBasicBlock::inferCoverBB(DInput *dInput, llvm::BasicBlock *b) {
        DBasicBlock *Db;

        std::string bname = b->getName().str();
        if (parent->BasicBlock.find(bname) != parent->BasicBlock.end()) {
            Db = parent->BasicBlock[bname];
            if (Db->state == CoverKind::untest || Db->state == CoverKind::uncover) {
                Db->setState(CoverKind::cover);
                Db->input.clear();
                this->addNewInput(dInput);
                return true;
            } else if (Db->state == CoverKind::cover) {
                this->addNewInput(dInput);
            }
        } else {
            std::cout << "inferCoverBB not find basic block name : " << bname << std::endl;
            parent->dump();
            exit(0);
        }
        return false;
    }

    void DBasicBlock::inferUncoverBB(llvm::BasicBlock *p, llvm::TerminatorInst *end, u_int i) {
        DBasicBlock *Dp;
        DBasicBlock *Db;
        std::string pname = p->getName().str();
        std::string bname = end->getSuccessor(i)->getName().str();
        if (parent->BasicBlock.find(bname) != parent->BasicBlock.end()) {
            Db = parent->BasicBlock[bname];
            if (parent->BasicBlock.find(pname) != parent->BasicBlock.end()) {
                Dp = parent->BasicBlock[pname];
                DInput *dInput = Dp->lastInput;

                std::vector<uint64_t> branch;
                for(uint64_t j = 0, e = end->getNumSuccessors(); j < e; j++) {
                    if(end->getSuccessor(j)->hasName()) {
                        auto n = end->getSuccessor(j)->getName();
                        if(this->parent->BasicBlock.find(n) != this->parent->BasicBlock.end()) {
                            branch.push_back(this->parent->BasicBlock[n]->trace_pc_address);
                        }
                    }
                }

                if (Db->state == CoverKind::untest) {
                    Db->setState(CoverKind::uncover);
                    Db->addNewInput(dInput);
                    dInput->addUncoveredAddress(Dp->trace_pc_address, Db->trace_pc_address, branch, i);
                } else if (Db->state == CoverKind::uncover) {
                    Db->addNewInput(dInput);
                    dInput->addUncoveredAddress(Dp->trace_pc_address, Db->trace_pc_address, branch, i);
                } else if (Db->state == CoverKind::cover) {

                }
#if DEBUG_INPUT
                if (Db->state == CoverKind::uncover) {
                    std::cout << "-------uncover basic block-----------------" << std::endl;
                    Db->dump();
                }
#endif
            } else {
                std::cout << "inferUncoverBB not find basic block pname : " << bname << std::endl;
            }

        } else {
            std::cout << "inferUncoverBB not find basic block bname : " << bname << std::endl;
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
                    inferUncoverBB(s, inst, i);
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
        if (this->input.find(i) == this->input.end()) {
            this->input[i] = 1U << i->idx;
        } else {
            this->input[i] = this->input[i] | 1U << i->idx;
        }
    }

    void DBasicBlock::dump() {
        std::cout << "********************************************" << std::endl;
        if (parent != nullptr) {
            std::cout << "Path : " << parent->Path << std::endl;
            std::cout << "FunctionName : " << parent->FunctionName << std::endl;
        }
        std::cout << "basicblock name : " << name << std::endl;
        std::cout << "AsmSourceCode : " << AsmSourceCode << std::endl;
        std::cout << "IR : " << IR << std::endl;
        std::cout << "CoverKind : " << state << std::endl;
        std::cout << "trace_pc_address : 0x" << std::hex << trace_pc_address << std::endl;
        if (this->basicBlock != nullptr) {
            dump_inst(this->basicBlock->getTerminator());
            std::string ld;
            llvm::raw_string_ostream rso(ld);
            this->basicBlock->print(rso);
            std::cout << ld;
        }

        for (auto i : this->input) {
            std::cout << "input : " << i.second << " : " << i.first->sig << std::endl;
            std::cout << i.first->program;
        }
        std::cout << "--------------------------------------------" << std::endl;
    }

    void DBasicBlock::real_dump(int kind) {
        std::cout << "********************************************" << std::endl;
        if (parent != nullptr) {
            std::cout << "Path              : " << parent->Path << std::endl;
            std::cout << "FunctionName      : " << parent->FunctionName << std::endl;
        }
        std::cout << "basicblock name   : " << name << std::endl;
        std::cout << "AsmSourceCode     : " << AsmSourceCode << std::endl;
        std::cout << "IR                : " << IR << std::endl;
        std::cout << "CoverKind         : " << state << std::endl;
        std::cout << "trace_pc_address  : 0x" << std::hex << trace_pc_address << std::endl;

        if (this->basicBlock != nullptr) {
            std::cout << "all dominator uncovered instructions : " <<
                      std::dec << this->get_number_all_dominator_uncovered_instructions() << std::endl;
            std::cout << "all arrive uncovered instructions : " << this->get_number_arrive_uncovered_instructions()
                      << std::endl;

            std::string ld;
            llvm::raw_string_ostream rso(ld);
            this->basicBlock->print(rso);
            auto bb = this->basicBlock;
            auto temp = bb->getNextNode();
            for (; temp != nullptr && !temp->hasName(); temp = bb->getNextNode()) {
                bb = temp;
                bb->print(rso);
            }
            // 0 is condition(br), 1 is uncovered branch, 2 is write statement(store)
            if (kind == 0) {
                auto inst = bb->getTerminator();
                dump_inst(inst);
            } else if (kind == 1) {
                auto inst = this->basicBlock->getFirstNonPHIOrDbgOrLifetime();
                dump_inst(inst);
            } else if (kind == 2) {
                for (temp = this->basicBlock;;) {
                    for (auto &inst : *temp) {
                        if (inst.getOpcode() == llvm::Instruction::Store) {
                            dump_inst(&inst);
                        }
                    }
                    temp = temp->getNextNode();
                    if (temp == nullptr || temp->hasName()) {
                        break;
                    }

                }
            }
            std::cout << ld;
        }

        for (auto i : this->input) {
            std::cout << "input : " << i.second << " : " << i.first->sig << std::endl;
            std::cout << i.first->program;
        }
        std::cout << "--------------------------------------------" << std::endl;
    }

    // not work if there is a switch with more than 64 cases.
    bool DBasicBlock::set_arrive(dra::DBasicBlock *db) {
        bool res = false;
        uint64_t Num = 0;
        auto *inst = dra::getFinalBB(this->basicBlock)->getTerminator();
        for (uint64_t i = 0, end = inst->getNumSuccessors(); i < end; i++) {
            if (inst->getSuccessor(i) == db->basicBlock) {
                Num = i;
                if (this->arrive.find(db) != this->arrive.end()) {
                    if ((this->arrive[db] & 1U << i) > 0) {

                    } else {
                        this->arrive[db] |= 1U << i;
                        res = true;
                    }
                } else {
                    this->arrive[db] = 1U << i;
                    res = true;
                }
            }
        }

        for (auto bb : db->arrive) {
            if (this->arrive.find(bb.first) != this->arrive.end()) {
                if ((this->arrive[bb.first] & 1U << Num) > 0) {

                } else {
                    this->arrive[bb.first] |= 1U << Num;
                    res = true;
                }
            } else {
                this->arrive[bb.first] = 1U << Num;
                res = true;
            }
        }
        return res;
    }

    void DBasicBlock::set_critical_condition() {

        if (this->basicBlock == nullptr) {

        } else {
            auto *fb = dra::getFinalBB(this->basicBlock);
            auto *inst = fb->getTerminator();
            auto successor_num = inst->getNumSuccessors();
            for (auto bb : this->arrive) {
                if (bb.second == ((1U << successor_num) - 1)) {

                } else {
                    bb.first->add_critical_condition(this, bb.second);
                }
            }
        }
    }

    void DBasicBlock::add_critical_condition(dra::DBasicBlock *db, uint64_t condition) {
        auto *c = new Condition();
        c->set_condition_address(db->trace_pc_address);
        c->set_successor(condition);
        auto *inst = dra::getFinalBB(db->basicBlock)->getTerminator();
        for (uint64_t i = 0, end = inst->getNumSuccessors(); i < end; i++) {
            auto temp = this->get_DB_from_bb(inst->getSuccessor(i));
            if (condition && 1U << i) {
                c->add_right_branch_address(temp->trace_pc_address);
//                (*c->mutable_right_branch_address())[temp->trace_pc_address] = 0;
            } else {
//                c->add_wrong_branch_address(temp->trace_pc_address);
//                (*c->mutable_wrong_branch_address())[temp->trace_pc_address] = 0;
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

    uint32_t DBasicBlock::get_number_uncovered_instructions() {
        if (this->state == CoverKind::cover) {
            return 0;
        } else {
            return this->number_instructions;
        }
    }

    void DBasicBlock::get_function_call(std::set<llvm::Function *> &res) {
        for (auto i : this->InstIR) {
            if (i->i->getOpcode() == llvm::Instruction::Call) {
                llvm::CallSite cs(i->i);
                llvm::Function *f = cs.getCalledFunction();
                if (f != nullptr && f != this->parent->function) {
                    res.insert(f);
                }
            }
        }
    }

    uint32_t DBasicBlock::get_number_arrive_uncovered_instructions() {
        std::set<llvm::Function *> uncovered_function;
        std::set<llvm::Function *> new_uncovered_functions;
        uncovered_function.insert(this->parent->function);
        uint32_t number_uncovered_instructions = this->get_number_uncovered_instructions();
        for (auto b : this->arrive) {
            if (b.first->state != CoverKind::cover) {
                number_uncovered_instructions =
                        number_uncovered_instructions + b.first->get_number_uncovered_instructions();
            }
            b.first->get_function_call(new_uncovered_functions);
        }
        while (!new_uncovered_functions.empty()) {
            std::set<llvm::Function *> temp;
            for (auto f : new_uncovered_functions) {
                uncovered_function.insert(f);
                DFunction *df = this->parent->parent->get_DF_from_f(f);
                number_uncovered_instructions += df->get_number_uncovered_instructions();
                df->get_function_call(temp);
            }
            new_uncovered_functions.clear();
            for (auto f : temp) {
                if (uncovered_function.find(f) == uncovered_function.end()) {
                    new_uncovered_functions.insert(f);
                    uncovered_function.insert(f);
                }
            }
        }
        return number_uncovered_instructions;
    }

    uint32_t DBasicBlock::get_number_all_dominator_uncovered_instructions() {
        if (this->basicBlock != nullptr) {
            return this->parent->get_number_dominator_uncovered_instructions(this->basicBlock);
        } else {
            return 0;
        }
    }

} /* namespace dra */

#endif /* LIB_DRA_DBASICBLOCK_CPP_ */