/*
 * BasicBlockAll.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "BasicBlockAll.h"

#include <llvm/ADT/ilist_iterator.h>
#include <llvm/IR/BasicBlock.h>
#include <iostream>
#include <iterator>

#include "FunctionAll.h"

namespace dra {

	BasicBlockAll::BasicBlockAll() :
			covNum(0), b(0), parent(0) {
		// TODO Auto-generated constructor stub
		state = Kind::other;

	}

	BasicBlockAll::~BasicBlockAll() {
		// TODO Auto-generated destructor stub
	}

	void BasicBlockAll::set(llvm::BasicBlock *b) {
		this->b = b;
		for (llvm::BasicBlock::iterator it = b->begin(); it != b->end(); it++) {
			InstructionLLVM *II = new InstructionLLVM();
			InstructionVector.push_back(II);

			II->i = (&*it);
		}
	}

	void BasicBlockAll::setState(Kind kind) {
		if (state == Kind::cover && kind == Kind::uncover) {
			std::cerr << "error BasicBlock kind" << "\n";
		}
		state = kind;
	}

	void BasicBlockAll::update(Kind kind) {
		setState(kind);
		for (std::vector<InstructionLLVM *>::iterator it = InstructionVector.begin(); it != InstructionVector.end();
				it++) {
			(*it)->setState(kind);
		}
		if (kind == Kind::cover) {
			parent->update(kind);
		}
	}

} /* namespace dra */
