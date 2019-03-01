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

//	void InstructionLLVM::setLine(std::vector<InformationOfSourceCode*> &IS) {
//
//		if (llvm::DebugLoc id = i->getDebugLoc()) {
//			llvm::DILocation *di = id.get();
//			fileName = getDSPIPath(di);
//			line = di->getLine();
//			for (std::vector<InformationOfSourceCode*>::iterator it = IS.begin(); it != IS.end(); it++) {
//
//				if ((*it)->fileName == fileName) {
//
//					if (line != 0 && line < (*it)->lineOfCode) {
//						(*it)->allLine.at(line - 1)->addInstruction(this);
//					} else {
//						std::cerr << "----------------------------------------------------------------" << "\n";
//						i->dump();
//						std::cerr << "name : " << fileName << " line : " << line << "\n";
//						std::cerr << "it fileName : " << (*it)->fileName << "\n";
//						std::cerr << "it line : " << (*it)->lineOfCode << "\n";
//						std::cerr << "error Instruction sline" << "\n";
//					}
//				}
//			}
//		}
//	}

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
