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
		repeat = false;

		function = nullptr;
		parent = nullptr;

		state = CoverKind::untest;

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

	void DFunction::setState(CoverKind kind) {
		if (state == CoverKind::cover && kind == CoverKind::uncover) {
			std::cerr << "error BasicBlock kind" << "\n";
		}
		state = kind;
	}

	void DFunction::update(CoverKind kind) {
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

	bool DFunction::isRepeat() const {
		return repeat;
	}

	static DFunction DFunction::MargeDFunction(DFunction *one, DFunction *two) {
		DFunction *f = new DFunction();
		return f;
	}

	void DFunction::setRepeat(bool repeat) {
		this->repeat = repeat;
	}

	void DFunction::setKind(FunctionKind kind) {
		switch (kind) {
			case dra::FunctionKind::IR: {
				setIR(true);
				break;
			}
			case dra::FunctionKind::O: {
				setObjudump(true);
				break;
			}
			case dra::FunctionKind::S: {
				setAsmSourceCode(true);
				break;
			}
			default: {

			}
		}
	}

	void DFunction::dump() {

		std::cout << "--------------------------------------------" << std::endl;
		std::cout << "Path :" << Path << std::endl;
		std::cout << "FunctionName :" << FunctionName << std::endl;

		std::cout << "Objudump :" << Objudump << std::endl;
		std::cout << "AsmSourceCode :" << AsmSourceCode << std::endl;
		std::cout << "IR :" << IR << std::endl;
		std::cout << "repeat :" << repeat << std::endl;
		std::cout << "CoverKind :" << state << std::endl;
		std::cout << "IRName :" << IRName << std::endl;
		std::cout << "Address :" << Address << std::endl;
		std::cout << "InstNum :" << InstNum << std::endl;
		std::cout << "CallInstNum :" << CallInstNum << std::endl;
		std::cout << "JumpInstNum :" << JumpInstNum << std::endl;
		std::cout << "BasicBlockNum :" << BasicBlockNum << std::endl;
		std::cout << "--------------------------------------------" << std::endl;

	}

} /* namespace dra */
