/*
 * DFunction.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DFunction.h"

#include <llvm/IR/Function.h>
#include <iostream>

namespace dra {

	DFunction::DFunction() {
		Objudump = false;
		AsmSourceCode = false;
		IR = false;

		function = nullptr;
		parent = nullptr;

		state = Kind::other;

		InstNum = 0;
		CallInstNum = 0;
		JumpInstNum = 0;
		BasicBlockNum = 0;
	}

	DFunction::~DFunction() = default;

	void DFunction::InitIRFunction(llvm::Function *f) {
		DFunction::function = f;
		for (auto &it : *function) {
			BasicBlockNum++;

			std::string Name = it.getName().str();
			DBasicBlock *b;
			if (BasicBlock.find(Name) == BasicBlock.end()) {

				b = new DBasicBlock();
				BasicBlock[Name] = b;
			}

			BasicBlock[Name]->setIr(true);
			BasicBlock[Name]->parent = this;
			BasicBlock[Name]->InitIRBasicBlock(&it);
		}
	}

	void DFunction::setState(Kind kind) {
		if (state == Kind::cover && kind == Kind::uncover) {
			std::cerr << "error BasicBlock kind" << "\n";
		}
		state = kind;
	}

	void DFunction::update(Kind kind) {
		setState(kind);
	}

	bool DFunction::isObjudump() const {
		return Objudump;
	}

	void DFunction::setObjudump(bool Objudump) {
		DFunction::Objudump = Objudump;
	}

	bool DFunction::isAsmSourceCode() const {
		return AsmSourceCode;
	}

	void DFunction::setAsmSourceCode(bool AsmSourceCode) {
		DFunction::AsmSourceCode = AsmSourceCode;
	}

	bool DFunction::isIR() const {
		return IR;
	}

	void DFunction::setIR(bool IR) {
		DFunction::IR = IR;
	}

	bool DFunction::isMap() {
		return DFunction::Objudump && DFunction::AsmSourceCode && DFunction::IR;
	}

} /* namespace dra */
