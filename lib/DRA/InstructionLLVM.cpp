/*
 * InstructionLLVM.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "InstructionLLVM.h"

#include <iostream>

#include "BasicBlockAll.h"


namespace dra {

	InstructionLLVM::InstructionLLVM() :
			i(0), parent(0), Line(0) {
		// TODO Auto-generated constructor stub
		state = Kind::other;

	}

	InstructionLLVM::~InstructionLLVM() {
		// TODO Auto-generated destructor stub
	}

	void InstructionLLVM::setState(Kind kind) {
		if (state == Kind::cover && kind == Kind::uncover) {
			std::cerr << "error Instruction kind" << "\n";
		}
		state = kind;
	}

	void InstructionLLVM::update(Kind kind) {
		setState(kind);
		parent->update(kind);
	}

} /* namespace dra */
