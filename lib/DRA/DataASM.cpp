/*
 * DataASM.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DataASM.h"

#include <bits/types/FILE.h>
#include <cstdio>
#include <fstream>
#include <initializer_list>
#include <iostream>
#include <sstream>
#include <utility>
#include <vector>

#include "BasicBlockAll.h"
#include "InstructionASM.h"

namespace dra {

	DataASM::DataASM() {
		// TODO Auto-generated constructor stub
		FindNum = 0;
		UnFindNum = 0;
		SameNum = 0;
		DiffNum = 0;

	}

	DataASM::~DataASM() {
		// TODO Auto-generated destructor stub
	}

	std::string DataASM::exec(std::string cmd) {
		std::string data;
		FILE * stream;
		const int max_buffer = 256;
		char buffer[max_buffer];
		cmd.append(" 2>&1");

		stream = popen(cmd.c_str(), "r");
		if (stream) {
			while (!feof(stream))
				if (fgets(buffer, max_buffer, stream) != NULL)
					data.append(buffer);
			pclose(stream);
		}
		return data;
	}

	void DataASM::ReadFromObjdump(std::string objdump) {
		std::string Line;
		std::string Addr;
		std::string FunctionName;
		std::string Path;
		std::string Cmd;
		std::string Result;
		std::stringstream ss;
		unsigned int LineNum;
		unsigned int inum;
		unsigned int start, end;

#if DEBUG
		std::cout << "objdump :" << objdump << std::endl;
#endif

		std::ifstream objdumpFile(objdump);
		inum = 0;
		LineNum = 0;
		FunctionName = "";
		if (objdumpFile.is_open()) {
			while (getline(objdumpFile, Line)) {
				LineNum++;
				if (Line.size() > 0) {
#if DEBUG
					std::cout << "Line :" << Line << std::endl;
#endif
					if (Line.find(">:") < Line.size()) {
#if DEBUG
						std::cout << ">: :" << std::endl;
#endif

						// get address
						ss.str("");
						for (int i = 0; i < 16; i++) {
							ss << Line.at(i);
						}
						Addr = ss.str();

						// get function name
						ss.str("");
						start = Line.find("<");
						end = Line.find(">");
						for (int i = start + 1; i < end; i++) {
							ss << Line.at(i);
						}
						FunctionName = ss.str();
#if DEBUG
						std::cout << "o FunctionName :" << FunctionName << std::endl;
#endif
						// get path
						Cmd = "addr2line -a -i -f -e vmlinux " + Addr;
#if DEBUG
						std::cout << "o Cmd :" << Cmd << std::endl;
#endif
						Result = exec(Cmd);
						ss.str("");
						start = Result.find("_bc/");
						end = Result.find(':');
						for (int i = start + 4; i < end; i++) {
							ss << Result.at(i);
						}
						Path = ss.str();
#if DEBUG
						std::cout << "o Path :" << Path << std::endl;
#endif
						if ((AllFunctiono.find(Path) != AllFunctiono.end())
								&& (AllFunctiono[Path].find(FunctionName) != AllFunctiono[Path].end())) {
							std::cout << "o repeat Path :" << Path << std::endl;
							std::cout << "o repeat FunctionName :" << FunctionName << std::endl;
						} else {
							AllFunctiono[Path].insert(
									std::pair<std::string, FunctionAll*>(FunctionName, new FunctionAll()));
							AllFunctiono[Path][FunctionName]->Name = FunctionName;
							AllFunctiono[Path][FunctionName]->Path = Path;
						}
					} else {
						//asm instruction
						if (Line.at(0) == '.') {
#if DEBUG
							std::cout << "dot :" << std::endl;
#endif
						} else if (Line.at(0) == 'D') {
#if DEBUG
							std::cout << "D :" << std::endl;
#endif
						} else if (Line.size() <= 40) {
							// deal with no asm
						} else if (Line.find("nop") < Line.size()) {
							// deal with nop
#if DEBUG
							std::cout << "nop :" << std::endl;
#endif
						} else if (Line.find("xchg") < Line.size()) {
							// deal with xchg
#if DEBUG
							std::cout << "xchg :" << std::endl;
#endif
						} else {
#if DEBUG
							std::cout << "else :" << std::endl;
#endif
							inum++;
							InstructionASM *inst = new InstructionASM();
							ss.str("");
							for (int i = 0; i < Line.size(); i++) {
								ss << Line.at(i);
							}
							inst->Inst = ss.str();
							AllFunctiono[Path][FunctionName]->InstASM.push_back(inst);
							if (ss.str().find("call") <= ss.str().size()) {
								AllFunctiono[Path][FunctionName]->CallInstNum++;
							}
							if (ss.str().find("j") <= ss.str().size()) {
								AllFunctiono[Path][FunctionName]->JumpInstNum++;
							}
						}
					}

				} else if (inum > 0) {
					// need add a space line at the end of objdump file.
					if (FunctionName != "") {
#if DEBUG
						std::cout << "Line :" << std::endl;
						std::cout << "FunctionName :" << FunctionName << std::endl;
						std::cout << "inum :" << inum << std::endl;
#endif
						AllFunctiono[Path][FunctionName]->InstNum = inum;
						inum = 0;
						FunctionName = "";
					}

				}
			}
			objdumpFile.close();
		} else {
			std::cerr << "Unable to open objdumpFile " << objdump << "\n";
		}
	}

	void DataASM::ReadFromAsmSourceCode(std::string AssemblySourceCode) {
		std::string line;
		std::string Path;
		std::string FunctionName;
		std::string BasicBlockName;
		std::stringstream ss;
		unsigned int LineNum;
		unsigned int FunctionLineNum;
		unsigned int inum = 0;
		unsigned int covNum;

#if DEBUG
		std::cout << "AssemblySourceCode :" << AssemblySourceCode << std::endl;
#endif

		std::ifstream AssemblySourceCodeFile(AssemblySourceCode);
		LineNum = 0;
		covNum = 0;
		if (AssemblySourceCodeFile.is_open()) {
			while (getline(AssemblySourceCodeFile, line)) {
				LineNum++;
				if (line.size() > 0) {
#if DEBUG
					std::cout << "line :" << line << std::endl;
#endif
					switch (line.at(0)) {
						case '.': {
							//label
#if DEBUG
							std::cout << "dot :" << std::endl;
#endif
							if (line.find(".Lfunc_end") < line.size()) {
								AllFunctions[Path][FunctionName]->InstNum = inum;
#if DEBUG
								std::cout << "FunctionName :" << FunctionName << std::endl;
								std::cout << "inum :" << inum << std::endl;
#endif
								inum = 0;
							} else if (line.find("# %") < line.size()) {
								ss.str("");
								for (int i = line.find('%') + 1; i < line.size(); i++) {
									ss << line.at(i);
								}
								BasicBlockName = ss.str();
								AllFunctions[Path][FunctionName]->BasicBlockVector.push_back(new BasicBlockAll());
								(AllFunctions[Path][FunctionName]->BasicBlockVector.back())->name = BasicBlockName;
								covNum = 0;
#if DEBUG
								std::cout << ". bb name :" << ss.str() << std::endl;
#endif
							}
							break;
						}
						case '#': {
							// bb
#if DEBUG
							std::cout << "sharp :" << std::endl;
#endif
							if (line.find("# %") < line.size()) {
								ss.str("");
								for (int i = line.find('%') + 1; i < line.size(); i++) {
									if (line.at(i) == '%') {
										for (i++; i < line.size(); i++) {
											ss << line.at(i);
										}
										AllFunctions[Path][FunctionName]->BasicBlockVector.push_back(new BasicBlockAll());
										(AllFunctions[Path][FunctionName]->BasicBlockVector.back())->name = ss.str();
										covNum = 0;

										break;
									}
								}

#if DEBUG
								std::cout << "# bb name :" << ss.str() << std::endl;
#endif
							}
							break;
						}
						case '	': {
#if DEBUG
							std::cout << "tab :" << std::endl;
							std::cout << "line.size() :" << line.size() << std::endl;
#endif
							if (line.size() == 1) {

							} else if (line.at(1) == '.') {
								if (Path == "" && FunctionName != "" && line.find('#') < line.size()) {
									if ((AllFunctions.find(Path) != AllFunctions.end())
											&& (AllFunctions[Path].find(FunctionName) != AllFunctions[Path].end())) {
										std::cout << "s repeat Path :" << Path << std::endl;
										std::cout << "s repeat FunctionName :" << FunctionName << std::endl;
									} else {

										ss.str("");
										for (int i = line.find('#') + 2; i < line.find(':'); i++) {
											ss << line.at(i);
										}
										Path = ss.str();
#if DEBUG
										std::cout << "o Path :" << Path << std::endl;
#endif
										AllFunctions[Path].insert(
												std::pair<std::string, FunctionAll*>(FunctionName, new FunctionAll()));
										AllFunctions[Path][FunctionName]->Name = FunctionName;
										AllFunctions[Path][FunctionName]->Path = Path;

									}
								}
							} else if (line.at(1) == '#') {

							} else if (line.at(1) >= 'a' && line.at(1) <= 'z') {
								//asm instruction
								InstructionASM *inst = new InstructionASM();
								inum++;
								ss.str("");
								for (int i = 1; i < line.size(); i++) {
									ss << line.at(i);
								}
								inst->Inst = ss.str();
								inst->BasicBlockName = BasicBlockName;
								AllFunctions[Path][FunctionName]->InstASM.push_back(inst);
#if DEBUG
								std::cout << "inst :" << inst->Inst << std::endl;
#endif
								if (ss.str().find("call") <= ss.str().size()) {
									AllFunctions[Path][FunctionName]->CallInstNum++;
									if (ss.str().find("__sanitizer_cov_trace_pc") <= ss.str().size()) {
										(AllFunctions[Path][FunctionName]->BasicBlockVector.back())->covNum++;
									}
								} else if (ss.str().find("j") <= ss.str().size()) {
									AllFunctions[Path][FunctionName]->JumpInstNum++;
								}
							}

							break;
						}
						case ' ': {
							//comment
#if DEBUG
							std::cout << "space :" << std::endl;
#endif
							break;
						}
						default: {
							if (line.find(':') < line.size()) {
								if (line.find("# @") < line.size()) {
									ss.str("");
									for (int i = 0; line.at(i) != ':'; i++) {
										ss << line.at(i);
									}
									FunctionName = ss.str();
									Path = "";
#if DEBUG
									std::cout << "FunctionName :" << FunctionName << std::endl;
#endif
								}
							}
						}
					}
				}
			}
			AssemblySourceCodeFile.close();
		} else {
			std::cerr << "Unable to open AssemblySourceCodeFile " << AssemblySourceCode << ">\n";
		}
	}

	void DataASM::Statistics() {
		std::cout << "Statistics :" << std::endl;
		for (auto it = AllFunctiono.begin(), ie = AllFunctiono.end(); it != ie; it++) {
			for (auto iit = (*it).second.begin(), iie = (*it).second.end(); iit != iie; iit++) {
				if (AllFunctions.find((*it).first) != AllFunctions.end()) {
					if (AllFunctions[(*it).first].find((*iit).first) != AllFunctions[(*it).first].end()) {
						FindNum++;
						if ((*iit).second->InstNum == AllFunctions[(*it).first][(*iit).first]->InstNum) {
							SameNum++;
						} else {
							DiffNum++;
							std::cout << "FunctionName of Num is different :" << (*it).first << std::endl;
							std::cout << "Num of s :" << AllFunctions[(*it).first][(*iit).first]->InstNum << std::endl;
							for (std::vector<InstructionASM *>::iterator iiit =
									AllFunctions[(*it).first][(*iit).first]->InstASM.begin(), iiie =
									AllFunctions[(*it).first][(*iit).first]->InstASM.end(); iiit != iiie; iiit++) {
								std::cout << (*iiit)->Inst << std::endl;
							}
							std::cout << "Num of o :" << AllFunctiono[(*it).first][(*iit).first]->InstNum << std::endl;
							for (std::vector<InstructionASM *>::iterator iiit =
									AllFunctiono[(*it).first][(*iit).first]->InstASM.begin(), iiie =
									AllFunctiono[(*it).first][(*iit).first]->InstASM.end(); iiit != iiie; iiit++) {
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

	void DataASM::Analysis(std::string AssemblySourceCode, std::string objdump) {

		ReadFromAsmSourceCode(AssemblySourceCode);
		ReadFromObjdump(objdump);
		Statistics();
	}


} /* namespace dra */
