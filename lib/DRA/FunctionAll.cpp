/*
 * FunctionAll.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "FunctionAll.h"

#include <llvm/ADT/ilist_iterator.h>
#include <llvm/IR/Function.h>
#include <iostream>
#include <iterator>
#include <vector>

#include "InstructionLLVM.h"

namespace dra {


	FunctionAll::FunctionAll() :
			f(0), parent(0) {
		// TODO Auto-generated constructor stub
		state = Kind::other;
		Name = "";
				InstNum = 0;
				CallInstNum = 0;
				JumpInstNum = 0;
				BasicBlockNum = 0;
	}

	FunctionAll::~FunctionAll() {
		// TODO Auto-generated destructor stub
	}

	void FunctionAll::set(llvm::Function *f) {
		this->f = f;
		for (llvm::Function::iterator it = f->begin(); it != f->end(); it++) {
			BasicBlockAll *IB = new BasicBlockAll();
			BasicBlockVector.push_back(IB);

			IB->set(&*it);
		}
	}

	void FunctionAll::setState(Kind kind) {
		if(state == Kind::cover && kind == Kind::uncover){
			std::cerr << "error BasicBlock kind" << "\n";
		}
		state = kind;
	}

	void FunctionAll::update(Kind kind) {
		setState(kind);
	}

} /* namespace dra */
