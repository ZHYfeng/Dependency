/*
 * InformationOfSourceCode.h
 *
 *  Created on: Nov 28, 2018
 *      Author: yhao
 */

#ifndef LIB_DRA_INFORMATIONOFSOURCECODE_H_
#define LIB_DRA_INFORMATIONOFSOURCECODE_H_

#include <string>
#include <vector>

#include "InformationOfLine.h"
#include "InstructionLLVM.h"

namespace llvm {
	class DILocation;
} /* namespace llvm */

namespace dra {

class InformationOfSourceCode {
public:
	InformationOfSourceCode(std::string OptionValue, std::string FileName);
	virtual ~InformationOfSourceCode();
	void setStart(unsigned int number);
	void setEnd(unsigned int number);
	void addToLine(Kind kind);
	void addCover(unsigned int number);
	void addUncover(unsigned int number);
	void addOther(unsigned int number);
	void dump();
	void setCover();
	void setUncover();
	void setState();
	static std::string getDSPIPath(llvm::DILocation *Loc);

public:
	std::string optionValue;
	std::string fileName;
	unsigned int start;
	unsigned int end;
	unsigned int lineOfCode;

	std::vector<InformationOfLine *> allLine;
	std::vector<unsigned int> cover;
	std::vector<unsigned int> uncover;
};

} /* namespace dra */

#endif /* LIB_DRA_INFORMATIONOFSOURCECODE_H_ */
