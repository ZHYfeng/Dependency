/*
 * DBasicBlock.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DBasicBlock.h"

#include <llvm/IR/Instructions.h>
#include <iostream>

#include "DFunction.h"

namespace dra {

	DBasicBlock::DBasicBlock() {
		IR = false;
		AsmSourceCode = false;
		basicBlock = nullptr;
		parent = nullptr;
		state = CoverKind::untest;
		COVNum = 0;
	}

	DBasicBlock::~DBasicBlock() = default;

	void DBasicBlock::InitIRBasicBlock(llvm::BasicBlock *b) {
		DBasicBlock::basicBlock = b;
		for (auto &it : *basicBlock) {
			DLInstruction *i;

			i = new DLInstruction();
			InstIR.push_back(i);

			i->parent = this;
			i->i = (&it);
		}
	}

	void DBasicBlock::setState(CoverKind kind) {

        if (kind == CoverKind::cover) {
            llvm::Instruction *i;
            for (auto &it : *basicBlock) {
                i = &it;
            }
            switch (i->getOpcode()) {
                case llvm::Instruction::Br: {
                    auto *bi = llvm::cast<llvm::BranchInst>(i);
                    if (bi->isUnconditional()) {
                        infer(bi->getSuccessor(0), CoverKind::cover);
                    } else {
                        infer(bi->getSuccessor(0), CoverKind::uncover);
                        infer(bi->getSuccessor(1), CoverKind::uncover);
                    }
                    break;
                }
                case llvm::Instruction::Switch: {
                    auto *si = llvm::cast<llvm::SwitchInst>(i);
                    for (unsigned int i = 0, end = si->getNumCases(); i < end; i++) {
                        infer(si->getSuccessor(i), CoverKind::uncover);
                    }
                    break;
                }
                default: {
                    parent->dump();
                    std::cout << "basic block name : " << name << std::endl;
                    exit(0);
                }
            }
            state = kind;
        } else if (kind == CoverKind::uncover) {
            if (state == CoverKind::cover) {
                parent->dump();
                std::cout << "cover to uncover basic block name : " << name << std::endl;
                exit(0);
            }
        }
	}

	void DBasicBlock::update(CoverKind kind) {
		setState(kind);
		for (auto it : InstIR) {
			it->setState(kind);
		}
		for (auto it : InstASM) {
			it->setState(kind);
		}
		if (kind == CoverKind::cover) {
			parent->update(kind);
		}
	}

} /* namespace dra */

bool dra::DBasicBlock::isAsmSourceCode() const {
	return AsmSourceCode;
}

void dra::DBasicBlock::setAsmSourceCode(bool asmSourceCode) {
	AsmSourceCode = asmSourceCode;
}

bool dra::DBasicBlock::isIr() const {
	return IR;
}

void dra::DBasicBlock::infer(llvm::BasicBlock *b, CoverKind kind) {
    DBasicBlock *Db;
    std::string name = b->getName().str();
    if (parent->BasicBlock.find(name) != parent->BasicBlock.end()) {
        Db = parent->BasicBlock[b->getName().str()];
        if (Db->state == CoverKind::untest) {
            Db->setState(kind);
        } else if (Db->state == CoverKind::uncover) {
            if (kind == CoverKind::cover) {
                Db->setState(kind);
            }
        } else if (Db->state == CoverKind::cover) {
            if (kind == CoverKind::uncover) {
                parent->dump();
                std::cout << "cover to uncover basic block name : " << name << std::endl;
                exit(0);
            }
        }
    } else {
        parent->dump();
        std::cout << "not find basic block name : " << name << std::endl;
        exit(0);
    }
}

void dra::DBasicBlock::setIr(bool ir) {
	IR = ir;
}
