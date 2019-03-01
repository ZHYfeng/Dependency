/*
 * InformationOfLine.cpp
 *
 *  Created on: Dec 10, 2018
 *      Author: yhao
 */

#include "InformationOfLine.h"

namespace dra {

InformationOfLine::InformationOfLine() {
	// TODO Auto-generated constructor stub
	state = Kind::other;
}

InformationOfLine::~InformationOfLine() {
	// TODO Auto-generated destructor stub
}

void InformationOfLine::dump() {
	// TODO Auto-generated destructor stub
}

void InformationOfLine::addInstruction(InstructionLLVM *i) {
	allInstruction.push_back(i);
}

void InformationOfLine::setState(Kind kind) {
	state = kind;
	for (std::vector<InstructionLLVM *>::iterator it =
			allInstruction.begin(); it != allInstruction.end(); it++) {
		(*it)->update(kind);
	}
}

} /* namespace dra */
