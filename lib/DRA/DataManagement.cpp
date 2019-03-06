/*
 * DataManagement.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DataManagement.h"

#include <llvm/IR/BasicBlock.h>
#include <llvm/IR/CFG.h>
#include <fstream>
#include <iostream>
#include <sstream>

namespace dra {

	DataManagement::DataManagement() {
		FindNum = 0;
		UnFindNum = 0;
		SameNum = 0;
		DiffNum = 0;
		Modules = new dra::DModule();
	}

	DataManagement::~DataManagement() = default;

	void DataManagement::initializeModule(std::string objdump, std::string AssemblySourceCode,
			std::string InputFilename) {

		Modules->ReadBC(std::move(InputFilename));
        Modules->ReadObjdump(std::move(objdump));
        Modules->ReadAsmSourceCode(std::move(AssemblySourceCode));

	}

	void dra::DataManagement::Statistics() {
		std::cout << "Statistics :" << std::endl;
		auto Function = Modules->Function;
		for (auto it = Function.begin(), ie = Function.end(); it != ie; it++) {
			for (auto iit = (*it).second.begin(), iie = (*it).second.end(); iit != iie; iit++) {
				if (Function.find((*it).first) != Function.end()) {
					if (Function[(*it).first].find((*iit).first) != Function[(*it).first].end()) {
						FindNum++;
						if ((*iit).second->InstNum == Function[(*it).first][(*iit).first]->InstNum) {
							SameNum++;
						} else {
							DiffNum++;
							std::cout << "FunctionName of Num is different :" << (*it).first << std::endl;
							std::cout << "Num of s :" << Function[(*it).first][(*iit).first]->InstNum << std::endl;
							for (auto iiit = Function[(*it).first][(*iit).first]->InstASM.begin(), iiie =
									Function[(*it).first][(*iit).first]->InstASM.end(); iiit != iiie; iiit++) {
								std::cout << (*iiit)->Inst << std::endl;
							}
							std::cout << "Num of o :" << Function[(*it).first][(*iit).first]->InstNum << std::endl;
							for (auto iiit = Function[(*it).first][(*iit).first]->InstASM.begin(), iiie =
									Function[(*it).first][(*iit).first]->InstASM.end(); iiit != iiie; iiit++) {
								std::cout << (*iiit)->Inst << std::endl;
							}
						}

					} else {
						UnFindNum++;
						std::cout << "not find FunctionName :" << (*iit).first << std::endl;
					}
				} else {
					UnFindNum++;
					std::cout << "not find Path :" << (*it).first << std::endl;
				}
			}

		}
		std::cout << "FindNum :" << FindNum << std::endl;
		std::cout << "UnFindNum :" << UnFindNum << std::endl;
		std::cout << "SameNum :" << SameNum << std::endl;
		std::cout << "DiffNum :" << DiffNum << std::endl;
	}

} /* namespace dra */
