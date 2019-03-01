/*
 * InformationOfSourceCode.cpp
 *
 *  Created on: Nov 28, 2018
 *      Author: yhao
 */

#include "InformationOfSourceCode.h"

#include <llvm/ADT/StringRef.h>
#include <llvm/IR/DebugInfoMetadata.h>
#include <iostream>
#include <iterator>

namespace dra {

InformationOfSourceCode::InformationOfSourceCode(std::string OptionValue,
		std::string FileName) :
		optionValue(OptionValue), fileName(FileName), lineOfCode(0), start(0), end(
				0) {

}

InformationOfSourceCode::~InformationOfSourceCode() {
	// TODO Auto-generated destructor stub
}

void InformationOfSourceCode::setStart(unsigned int number) {
	start = number;
}

void InformationOfSourceCode::setEnd(unsigned int number) {
	end = number;
	lineOfCode = end - start;
}

void InformationOfSourceCode::addToLine(Kind kind) {
	InformationOfLine *l = new InformationOfLine();
	allLine.push_back(l);
	l->state = kind;
}

void InformationOfSourceCode::addCover(unsigned int number) {
	cover.push_back(number);
	addToLine(Kind::cover);

}
void InformationOfSourceCode::addUncover(unsigned int number) {
	uncover.push_back(number);
	addToLine(Kind::uncover);

}
void InformationOfSourceCode::addOther(unsigned int number) {
	addToLine(Kind::other);
}

void InformationOfSourceCode::dump() {
	std::cerr
			<< "--------------------------------------------------------------"
			<< "\n";
	std::cerr << "optionValue : " << optionValue << "\n";
	std::cerr << "fileName : " << fileName << "\n";
	std::cerr << "start : " << start << "\n";
	std::cerr << "end : " << end << "\n";
	std::cerr << "lineOfCode : " << lineOfCode << "\n";

	std::cerr << "stateOfLine : " << "\n";
	for (std::vector<InformationOfLine *>::iterator i = allLine.begin();
			i != allLine.end(); i++) {
		(*i)->dump();

	}
	std::cerr << "\n";
	std::cerr << "cover : " << "\n";
	for (std::vector<unsigned int>::iterator i = cover.begin(); i < cover.end();
			i++) {
		std::cerr << *i << " ";

	}
	std::cerr << "\n";
	std::cerr << "uncover : " << "\n";
	for (std::vector<unsigned int>::iterator i = uncover.begin();
			i < uncover.end(); i++) {
		std::cerr << *i << " ";

	}
	std::cerr << "\n";

}

void InformationOfSourceCode::setCover() {
	for (std::vector<unsigned int>::iterator it = cover.begin();it != cover.end(); it++) {
		allLine.at(*it-start)->setState(Kind::cover);
	}
}

void InformationOfSourceCode::setUncover() {
	for (std::vector<unsigned int>::iterator it = uncover.begin();it != uncover.end(); it++) {
		allLine.at(*it-start)->setState(Kind::uncover);
	}
}

void InformationOfSourceCode::setState() {
	if (lineOfCode == allLine.size()) {
		setCover();
		setUncover();
	} else {
		std::cerr << "error lineOfCode or allLine.size" << "\n";
	}



}

std::string InformationOfSourceCode::getDSPIPath(llvm::DILocation *Loc) {
	std::string dir = Loc->getDirectory();
	std::string file = Loc->getFilename();
	if (file[0] == '.') {
		return file;
	} else {
		return "./" + file;
	}
}

} /* namespace dra */
