/*
 * DataManagement.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DataManagement.h"

#include <llvm/ADT/StringRef.h>
#include <llvm/IR/BasicBlock.h>
#include <llvm/IR/CFG.h>
#include <fstream>
#include <iostream>
#include <iterator>
#include <sstream>
#include <unordered_map>

#include "BasicBlockAll.h"
#include "DataASM.h"

namespace dra {


	DataManagement::DataManagement() {
		DataASM = new dra::DataASM();
		// TODO Auto-generated constructor stub

	}

	DataManagement::~DataManagement() {
		// TODO Auto-generated destructor stub
	}

	void DataManagement::initializeModule(std::string InputFilename) {
		Modules->initializeModule(InputFilename);
	}

	void DataManagement::MapBBfromStoBC() {
		std::stringstream ss;
		for (auto it = DataASM->AllFunctions.begin(), ie = DataASM->AllFunctions.end(); it != ie; it++) {
			std::string path = (*it).first;
			for (auto Functionib = (*it).second.begin(), Functionie = (*it).second.end(); Functionib != Functionie;
					Functionib++) {
				std::string name = (*Functionib).first;
#if DEBUG
				std::cout << "---------------------------------------" << std::endl;
				std::cout << "function name : " << name << std::endl;
				for (auto BBib = (Functionib)->second->BasicBlockVector.begin(), BBie = (Functionib)->second->BasicBlockVector.end(); BBib != BBie;
						BBib++) {
					std::cout << "s BBname : " << (*BBib)->name << std::endl;
					std::cout << "s BB covNum : " << (*BBib)->covNum << std::endl;
				}
#endif
					auto mm = Modules->AllFunctionbc;
					if ((mm.find(path) != mm.end()) && (mm[path].find(name) != mm[path].end())) {
						auto function = mm[path][name];
						for (auto BBiib = function->f->begin(), BBiie = function->f->end(); BBiib != BBiie; BBiib++) {
							llvm::BasicBlock *B = &(*BBiib);
							std::string BBname = B->getName().str();
							std::cout << "ll BBname : " << BBname << std::endl;
							for (auto predib = llvm::pred_begin(B), predie = llvm::pred_end(B); predib != predie;
									++predib) {
								llvm::BasicBlock* predecessor = *predib;
								std::cout << "predecessor name : " << predecessor->getName().str() << std::endl;
							}
						}

						for (auto BBiib = (Functionib)->second->BasicBlockVector.begin(), BBiie = (Functionib)->second->BasicBlockVector.end();
								BBiib != BBiie; BBiib++) {
							if ((*BBiib)->covNum > 0) {
								std::cout << "s BBname : " << (*BBiib)->name << std::endl;
								if ((*BBiib)->name.find(".i.") < (*BBiib)->name.size()) {
									//inline?

								} else if ((*BBiib)->name.find(".exit") < (*BBiib)->name.size()) {

								} else if ((*BBiib)->name.find("_crit_edge") < (*BBiib)->name.size()) {
									unsigned int end;
									if ((*BBiib)->name.find(".backedge") < (*BBiib)->name.size()) {
										end = (*BBiib)->name.find(".backedge");
									} else {
										end = (*BBiib)->name.find("_crit_edge");
									}
									bool Isfind = false;
									for (unsigned int start = 0; start < (*BBiib)->name.size() && !Isfind; start++) {
										if ((*BBiib)->name.at(start) == '.') {
											ss.str("");
											for (int i = start + 1; i < end; i++) {
												ss << (*BBiib)->name.at(i);
											}
											std::string realname = ss.str();
											std::cout << "realname : " << realname << std::endl;
											std::cout << "edge Isfind ? : ";
											for (auto BBiiib = function->f->begin(), BBiiie = function->f->end();
													BBiiib != BBiiie && !Isfind; BBiiib++) {
												llvm::BasicBlock *B = &(*BBiiib);
												std::string BBname = B->getName().str();
												if (BBname == realname) {
													std::cout << "dot find!" << std::endl;
													Isfind = true;
												}
											}
										}
									}
								} else {
									std::cout << "Isfind ? : ";
									for (auto BBiiib = function->f->begin(), BBiiie = function->f->end();
											BBiiib != BBiiie; BBiiib++) {

										llvm::BasicBlock *B = &(*BBiiib);
										std::string BBname = B->getName().str();
										if ((*BBiib)->name == BBname) {
											std::cout << "find!" << std::endl;
											break;
										}
									}
								}
							}
						}
				}
			}
		}
	}

} /* namespace dra */
