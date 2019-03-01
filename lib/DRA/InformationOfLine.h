/*
 * InformationOfLine.h
 *
 *  Created on: Dec 10, 2018
 *      Author: yhao
 */

#ifndef LIB_DRA_INFORMATIONOFLINE_H_
#define LIB_DRA_INFORMATIONOFLINE_H_

#include <vector>

#include "InstructionLLVM.h"

namespace dra {
	class InstructionLLVM;
} /* namespace dra */

namespace dra {

class InformationOfLine {
public:
	InformationOfLine();
	virtual ~InformationOfLine();
	void dump();
	void addInstruction(InstructionLLVM *i);
	void setState(Kind kind);

public:
	Kind state;
	std::vector<InstructionLLVM *> allInstruction;
};

} /* namespace dra */

#endif /* LIB_DRA_INFORMATIONOFLINE_H_ */
