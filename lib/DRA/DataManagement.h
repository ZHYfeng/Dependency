/*
 * DataManagement.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_DATAMANAGEMENT_H_
#define LIB_DRA_DATAMANAGEMENT_H_

#include <string>
#include <unordered_map>

#include "DAInstruction.h"
#include "DFunction.h"
#include "DInput.h"

namespace dra {
class DModule;
} /* namespace dra */

namespace dra {

class DataManagement {
public:
	DataManagement();
	virtual ~DataManagement();

	void initializeModule(std::string objdump, std::string AssemblySourceCode, std::string InputFilename);
	void BuildAddress2BB(std::unordered_map<std::string, std::unordered_map<std::string, DFunction *>> Function);

	void getVmOffsets(std::string vmOffsets);
	void getInput(std::string coverfile);
	void setInput();

public:
	dra::DModule *Modules;
	std::unordered_map<unsigned long long int, DAInstruction *> Address2BB;
	std::unordered_map<std::string, DInput *> Inputs;
	std::set<unsigned long long int> cover;
	unsigned long long int vmOffsets;

};

} /* namespace dra */

#endif /* LIB_DRA_DATAMANAGEMENT_H_ */
