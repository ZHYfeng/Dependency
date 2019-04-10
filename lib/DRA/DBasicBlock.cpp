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
	this->lastInput = nullptr;
	this->realPred = nullptr;
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
			std::cout << "cover to uncover basic block name : " << name << std::endl;
			parent->dump();
			exit(0);
		}
	}
}

void DBasicBlock::update(CoverKind kind, DInput * input) {
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
	inferCoverBB(input, this->basicBlock);
	infer();
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

void DBasicBlock::inferCoverBB(DInput * input, llvm::BasicBlock *b) {
	DBasicBlock *Db;
	std::string name = b->getName().str();
	if (parent->BasicBlock.find(name) != parent->BasicBlock.end()) {
		Db = parent->BasicBlock[name];
		if (Db->state == CoverKind::untest || Db->state == CoverKind::uncover) {
			Db->setState(CoverKind::cover);
			Db->input.clear();
			this->addNewInput(input);
		} else if (Db->state == CoverKind::cover) {
			this->addNewInput(input);
		}
	} else {
		std::cout << "inferCoverBB not find basic block name : " << name << std::endl;
		parent->dump();
		exit(0);
	}
}

void DBasicBlock::inferUncoverBB(llvm::BasicBlock *p, llvm::BasicBlock *b) {
	DBasicBlock *Db;
	std::string name = b->getName().str();
	std::string pname = p->getName().str();
	if (parent->BasicBlock.find(name) != parent->BasicBlock.end()) {
		Db = parent->BasicBlock[name];
		if (Db->state == CoverKind::untest) {
			Db->setState(CoverKind::uncover);
			Db->addNewInput(parent->BasicBlock[pname]->lastInput);
			Db->realPred = parent->BasicBlock[pname];
		} else if (Db->state == CoverKind::uncover) {
			Db->addNewInput(parent->BasicBlock[pname]->lastInput);
			Db->realPred = parent->BasicBlock[pname];
		} else if (Db->state == CoverKind::cover) {

		}
#if DEBUGINPUT
		if (Db->state == CoverKind::uncover) {
			std::cout << "-------uncover basic block-----------------" << std::endl;
			Db->dump();
		}
#endif
	} else {
		std::cout << "inferUncoverBB not find basic block name : " << name << std::endl;
		parent->dump();
		exit(0);
	}
}

void DBasicBlock::inferSuccessors(llvm::BasicBlock *b) {
	auto *inst = b->getTerminator();
	std::string name = b->getName().str();
	auto input = this->parent->BasicBlock[name]->lastInput;
	if (inst->getNumSuccessors() == 1) {
//		setOtherBBState(b, inst->getSuccessor(0), CoverKind::cover);
//		inferSuccessors(inst->getSuccessor(0));
	} else {
		for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
			std::string name = inst->getSuccessor(i)->getName().str();
			if (inst->getSuccessor(i)->hasName()) {
				inferUncoverBB(b, inst->getSuccessor(i));
			} else {
				inferCoverBB(input, inst->getSuccessor(i));
				inferSuccessors(inst->getSuccessor(i));
			}
		}
	}
}

void DBasicBlock::inferPredecessors(llvm::BasicBlock *b) {
	std::string name = b->getName().str();
	auto input = this->parent->BasicBlock[name]->lastInput;
	if (b->getSinglePredecessor()) {
		inferCoverBB(input, b->getSinglePredecessor());
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
			inferCoverBB(input, pb);
			inferPredecessors(pb);
			inferPredecessorsUncover(b, pb);
		}
	} else if (!b->hasName()) {
		for (auto *Pred : llvm::predecessors(b)) {
			inferCoverBB(input, b->getSinglePredecessor());
			inferPredecessors(b->getSinglePredecessor());
			inferPredecessorsUncover(b, b->getSinglePredecessor());
		}
	} else {

	}

}

void DBasicBlock::inferPredecessorsUncover(llvm::BasicBlock *b, llvm::BasicBlock *Pred) {
	auto *inst = Pred->getTerminator();
	if (inst->getNumSuccessors() == 1) {

	} else {
		for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
			if (inst->getSuccessor(i) != b && inst->getSuccessor(i)->hasName()) {
				inferUncoverBB(b, inst->getSuccessor(i));
			}
		}
	}
}

void DBasicBlock::infer() {
	if (this->state == CoverKind::cover) {
		inferSuccessors(this->basicBlock);
//		inferPredecessors(this->basicBlock);
	}

}

void DBasicBlock::addNewInput(DInput* i) {
	this->lastInput = i;
	this->input.insert(i);
}

void DBasicBlock::dump() {

	std::cout << "--------------------------------------------" << std::endl;
	std::cout << "Path :" << parent->Path << std::endl;
	std::cout << "FunctionName :" << parent->FunctionName << std::endl;
	std::cout << "name :" << name << std::endl;

	std::cout << "AsmSourceCode :" << AsmSourceCode << std::endl;
	std::cout << "IR :" << IR << std::endl;
	std::cout << "CoverKind :" << state << std::endl;
	basicBlock->dump();
	if (realPred != nullptr) {
		std::cout << "realPred :" << realPred->name << std::endl;
		std::cout << "lastInput :" << lastInput->sig << std::endl;
	}
	std::cout << "--------------------------------------------" << std::endl;

}

} /* namespace dra */
