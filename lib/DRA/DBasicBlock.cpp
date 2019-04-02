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
		infer();
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

void dra::DBasicBlock::infer() {

	if (this->state == CoverKind::cover) {
		inferSuccessors(this->basicBlock);
		inferPredecessors(this->basicBlock);
	}

}

void dra::DBasicBlock::setIr(bool ir) {
	IR = ir;
}

void dra::DBasicBlock::inferSuccessors(llvm::BasicBlock *b) {
	auto *inst = b->getTerminator();
	if (inst->getNumSuccessors() == 1) {
		setOtherBBState(inst->getSuccessor(0), CoverKind::cover);
		inferSuccessors(inst->getSuccessor(0));
	} else {
		for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
			std::string name = inst->getSuccessor(i)->getName().str();
			if (inst->getSuccessor(i)->hasName()) {
				setOtherBBState(inst->getSuccessor(i), CoverKind::uncover);
			} else {
				setOtherBBState(inst->getSuccessor(i), CoverKind::cover);
				inferSuccessors(inst->getSuccessor(i));
			}
		}
	}
}

void dra::DBasicBlock::setOtherBBState(llvm::BasicBlock *b, dra::CoverKind kind) {
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

void dra::DBasicBlock::inferPredecessors(llvm::BasicBlock *b) {

	if (b->getSinglePredecessor()) {
		setOtherBBState(b->getSinglePredecessor(), CoverKind::cover);
		inferPredecessors(b->getSinglePredecessor());
		inferPredecessorsUncover(b, b->getSinglePredecessor());
	} else if (useLessPred.size() > 0) {
		int num = 0;
		llvm::BasicBlock *pb;
		for (auto *Pred : llvm::predecessors(b)) {
			if (useLessPred.find(Pred) == useLessPred.end()) {
				pb = Pred;
				num++;
			}
		}
		if (num == 1) {
			setOtherBBState(pb, CoverKind::cover);
			inferPredecessors(pb);
			inferPredecessorsUncover(b, pb);
		}
	} else if(!b->hasName()){
		for (auto *Pred : llvm::predecessors(b)) {
			setOtherBBState(b->getSinglePredecessor(), CoverKind::cover);
			inferPredecessors(b->getSinglePredecessor());
			inferPredecessorsUncover(b, b->getSinglePredecessor());
		}
	} else{

	}

}

void dra::DBasicBlock::inferPredecessorsUncover(llvm::BasicBlock *b, llvm::BasicBlock *Pred) {
	auto *inst = Pred->getTerminator();
	if (inst->getNumSuccessors() == 1) {

	} else {
		for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
			if (inst->getSuccessor(i) != b && inst->getSuccessor(i)->hasName()) {
				setOtherBBState(inst->getSuccessor(i), CoverKind::uncover);
			}
		}
	}
}
