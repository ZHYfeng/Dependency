/*
 * DBasicBlock.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DBasicBlock.h"

#include <llvm/ADT/ilist_iterator.h>
#include <llvm/IR/BasicBlock.h>
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
		if (state == CoverKind::cover && kind == CoverKind::uncover) {
			std::cerr << "error BasicBlock kind" << "\n";
		}
		state = kind;
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

void dra::DBasicBlock::setIr(bool ir) {
	IR = ir;
}
