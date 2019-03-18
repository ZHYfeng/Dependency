/*
 * DataManagement.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DataManagement.h"

#include <algorithm>

#include "DModule.h"

#define PATH_SIZE 1000000

namespace dra {

	DataManagement::DataManagement() {
		FindNum = 0;
		UnFindNum = 0;
		SameNum = 0;
		DiffNum = 0;
		Modules = new dra::DModule();
		Address2BB.reserve(1000000);
	}

	DataManagement::~DataManagement() = default;

	void DataManagement::initializeModule(std::string objdump, std::string AssemblySourceCode, std::string InputFilename) {

		Modules->ReadBC(std::move(InputFilename));
		Modules->ReadObjdump(std::move(objdump));
		Modules->ReadAsmSourceCode(std::move(AssemblySourceCode));

	}

	void DataManagement::BuildAddress2BB(std::unordered_map<std::string, std::unordered_map<std::string, DFunction *>> Function) {
		for (auto file : Modules->Function) {
			for (auto function : file.second) {
				if (function.second->isRepeat()) {

				} else {
					for (auto inst : function.second->InstASM) {
						Address2BB[inst->Address] = inst;
					}
				}
			}
		}
	}

} /* namespace dra */
