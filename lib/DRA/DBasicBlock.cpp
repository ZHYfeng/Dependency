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

	DBasicBlock::DBasicBlock() = default;

	DBasicBlock::~DBasicBlock() =default;

	void DBasicBlock::InitIRBasicBlock(llvm::BasicBlock *b) {
		DBasicBlock::basicBlock = b;
		for (auto &it : *basicBlock) {
			DLInstruction *i;
			i = new DLInstruction();
			InstIR.push_back(i);

			i->i = (&it);
		}
	}

	void DBasicBlock::setState(Kind kind) {
		if (state == Kind::cover && kind == Kind::uncover) {
			std::cerr << "error BasicBlock kind" << "\n";
		}
		state = kind;
	}

	void DBasicBlock::update(Kind kind) {
		setState(kind);
		for (std::vector<DLInstruction *>::iterator it = InstIR.begin(); it != InstIR.end();
				it++) {
			(*it)->setState(kind);
		}
		if (kind == Kind::cover) {
			parent->update(kind);
		}
	}

} /* namespace dra */
