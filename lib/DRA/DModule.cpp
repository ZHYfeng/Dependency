#include "DModule.h"

#include <bits/types/FILE.h>
#include <llvm/ADT/simple_ilist.h>
#include <llvm/ADT/SmallVector.h>
#include <llvm/ADT/StringRef.h>
#include <llvm/IR/DebugInfoMetadata.h>
#include <llvm/IR/Function.h>
#include <llvm/IR/LLVMContext.h>
#include <llvm/IR/Metadata.h>
#include <llvm/IR/Module.h>
#include <llvm/IRReader/IRReader.h>
#include <llvm/Support/Casting.h>
#include <llvm/Support/SourceMgr.h>
#include <cstdio>
#include <cstdlib>
#include <fstream>
#include <iostream>
#include <memory>
#include <sstream>
#include <utility>
#include <vector>

#include "DAInstruction.h"
#include "DBasicBlock.h"

#define PATH_SIZE 10000

namespace dra {

	DModule::DModule() {
		// TODO Auto-generated constructor stub
//		Function.reserve(PATH_SIZE);

	}

	DModule::~DModule() = default;

	void DModule::ReadBC(std::string InputFilename) {
		std::unique_ptr<llvm::Module> module;
		llvm::LLVMContext *cxts;
		llvm::SMDiagnostic Err;
		cxts = new llvm::LLVMContext[1];
		module = llvm::parseIRFile(InputFilename, Err, cxts[0]);

		if (!module) {
			std::cerr << "load module: " << InputFilename << " failed\n";
			exit(0);
		} else {

			std::cerr << "size : " << module->getNamedMDList().size() << "\n";
			for (auto &i : module->getNamedMDList()) {
				i.dump();
			}
			BuildLLVMFunction(module.get());
		}
	}

	void DModule::BuildLLVMFunction(llvm::Module *Module) {

		for (auto &it : *Module) {
			llvm::SmallVector<std::pair<unsigned, llvm::MDNode *>, 4> MDs;
			it.getAllMetadata(MDs);
			for (auto &MD : MDs) {
				MD.second->dump();
				if (llvm::MDNode *N = MD.second) {
					if (auto *SP = llvm::dyn_cast<llvm::DISubprogram>(N)) {
						std::string Path = SP->getFilename().str();
						std::string name = it.getName().str();
						std::string FunctionName;
						if (name.find('.') < name.size()) {
							FunctionName = name.substr(0, name.find('.'));
						} else {
							FunctionName = name;
						}

						if ((Function.find(Path) != Function.end())
								&& (Function[Path].find(FunctionName) != Function[Path].end())) {
							if (Function[Path][FunctionName]->isObjudump()) {

							} else if (Function[Path][FunctionName]->isIR()) {
								std::cout << "--------------------------------------------" << std::endl;
								std::cout << "ir repeat function : " << std::endl;
								std::cout << "Path : " << Path << std::endl;
								std::cout << "name : " << name << std::endl;
								std::cout << "FunctionName : " << FunctionName << std::endl;
							} else if (Function[Path][FunctionName]->isAsmSourceCode()) {

							}

						} else {
							Function[Path].insert(std::pair<std::string, DFunction *>(FunctionName, new DFunction()));
						}

						Function[Path][FunctionName]->setIR(true);
						Function[Path][FunctionName]->InitIRFunction(&it);
					}

				}
			}

		}
	}

	std::string dra::DModule::exec(std::string cmd) {
		std::string data;
		FILE *stream;
		const int max_buffer = 256;
		char buffer[max_buffer];
		cmd.append(" 2>&1");

		stream = popen(cmd.c_str(), "r");
		if (stream) {
			while (!feof(stream))
				if (fgets(buffer, max_buffer, stream) != nullptr)
					data.append(buffer);
			pclose(stream);
		}
		return data;
	}

	void dra::DModule::ReadObjdump(std::string objdump) {
		std::string Line;
		std::string Addr;
		std::string FunctionName;
		std::string Path;
		std::string Cmd;
		std::string Result;
		std::stringstream ss;
		unsigned int LineNum;
		unsigned int InsNum;
		unsigned long end;
		unsigned long start;

#if DEBUGOBJDUMP
		std::cout << "objdump :" << objdump << std::endl;
#endif

		std::ifstream objdumpFile(objdump);
		InsNum = 0;
		LineNum = 0;
		FunctionName = "";
		if (objdumpFile.is_open()) {
			while (getline(objdumpFile, Line)) {
				LineNum++;
				if (!Line.empty()) {
#if DEBUGOBJDUMP
					std::cout << "Line :" << Line << std::endl;
#endif
					if (Line.find(">:") < Line.size()) {
#if DEBUGOBJDUMP
						std::cout << ">: :" << std::endl;
#endif

						// get address
						ss.str("");
						for (unsigned long i = 0; i < 16; i++) {
							ss << Line.at(i);
						}
						Addr = ss.str();

						// get function name
						ss.str("");
						start = Line.find('<');
						end = Line.find('>');
						for (unsigned long i = start + 1; i < end; i++) {
							ss << Line.at(i);
						}
						FunctionName = ss.str();
#if DEBUGOBJDUMP
						std::cout << "o FunctionName :" << FunctionName << std::endl;
#endif
						// get path
						Cmd = "addr2line -a -i -function -e vmlinux " + Addr;
#if DEBUGOBJDUMP
						std::cout << "o Cmd :" << Cmd << std::endl;
#endif
						Result = exec(Cmd);
						ss.str("");
						start = Result.find("_bc/");
						end = Result.find(':');
						for (unsigned long i = start + 4; i < end; i++) {
							ss << Result.at(i);
						}
						Path = ss.str();
#if DEBUGOBJDUMP
						std::cout << "o Path :" << Path << std::endl;
#endif
						if ((Function.find(Path) != Function.end())
								&& (Function[Path].find(FunctionName) != Function[Path].end())) {

							if (Function[Path][FunctionName]->isObjudump()) {
								std::cout << "--------------------------------------------" << std::endl;
								std::cout << "o repeat Path :" << Path << std::endl;
								std::cout << "o repeat FunctionName :" << FunctionName << std::endl;
							} else if (Function[Path][FunctionName]->isIR()) {

							} else if (Function[Path][FunctionName]->isAsmSourceCode()) {

							}

						} else {
							Function[Path].insert(std::pair<std::string, DFunction *>(FunctionName, new DFunction()));

						}

						Function[Path][FunctionName]->Name = FunctionName;
						Function[Path][FunctionName]->Path = Path;

					} else {
						//asm instruction
						if (Line.at(0) == '.') {
#if DEBUGOBJDUMP
							std::cout << "dot :" << std::endl;
#endif
						} else if (Line.at(0) == 'D') {
#if DEBUGOBJDUMP
							std::cout << "D :" << std::endl;
#endif
						} else if (Line.size() <= 40) {
							// deal with no asm
						} else if (Line.find("nop") < Line.size()) {
							// deal with nop
#if DEBUGOBJDUMP
							std::cout << "nop :" << std::endl;
#endif
						} else if (Line.find("xchg") < Line.size()) {
							// deal with xchg
#if DEBUGOBJDUMP
							std::cout << "xchg :" << std::endl;
#endif
						} else {
#if DEBUGOBJDUMP
							std::cout << "else :" << std::endl;
#endif
							InsNum++;
							auto *inst = new DAInstruction();
							ss.str("");
							for (char i : Line) {
								ss << i;
							}
							inst->Inst = ss.str();
							Function[Path][FunctionName]->InstASM.push_back(inst);
							if (ss.str().find("call") <= ss.str().size()) {
								Function[Path][FunctionName]->CallInstNum++;
							}
							if (ss.str().find('j') <= ss.str().size()) {
								Function[Path][FunctionName]->JumpInstNum++;
							}
						}
					}

				} else if (InsNum > 0) {
					// need add a space line at the end of objdump file.
					if (!FunctionName.empty()) {
#if DEBUGOBJDUMP
						std::cout << "Line :" << std::endl;
						std::cout << "FunctionName :" << FunctionName << std::endl;
						std::cout << "InsNum :" << InsNum << std::endl;
#endif
						Function[Path][FunctionName]->InstNum = InsNum;
						InsNum = 0;
						FunctionName = "";
					}
				}
			}
			objdumpFile.close();
		} else {
			std::cerr << "Unable to open objdump file " << objdump << "\n";
		}
	}

	void dra::DModule::ReadAsmSourceCode(std::string AssemblySourceCode) {
		std::string line;
		std::string Path;
		std::string FunctionName;
		std::string BasicBlockName;
		std::stringstream ss;
		unsigned int LineNum;
		unsigned int InsNum = 0;
		unsigned int COVNum;

#if DEBUGASM
		std::cout << "AssemblySourceCode :" << AssemblySourceCode << std::endl;
#endif

		std::ifstream AssemblySourceCodeFile(AssemblySourceCode);
		LineNum = 0;
		COVNum = 0;
		if (AssemblySourceCodeFile.is_open()) {
			while (getline(AssemblySourceCodeFile, line)) {
				LineNum++;
				if (!line.empty()) {
#if DEBUGASM
					std::cout << "line :" << line << std::endl;
#endif
					switch (line.at(0)) {
						case '.': {
							//label
#if DEBUGASM
							std::cout << "dot :" << std::endl;
#endif
							if (line.find(".Lfunc_end") < line.size()) {
								Function[Path][FunctionName]->InstNum = InsNum;
#if DEBUGASM
								std::cout << "FunctionName :" << FunctionName << std::endl;
								std::cout << "InsNum :" << InsNum << std::endl;
#endif
								InsNum = 0;
							} else if (line.find("# %") < line.size()) {
								if (Function[Path][FunctionName]->BasicBlock.find(BasicBlockName)
										!= Function[Path][FunctionName]->BasicBlock.end()) {
									(Function[Path][FunctionName]->BasicBlock[BasicBlockName])->COVNum = COVNum;
								}

								ss.str("");
								for (unsigned long i = line.find('%') + 1; i < line.size(); i++) {
									ss << line.at(i);
								}
								BasicBlockName = ss.str();

								Function[Path][FunctionName]->BasicBlock[BasicBlockName] = new DBasicBlock();
								(Function[Path][FunctionName]->BasicBlock[BasicBlockName])->name = BasicBlockName;
								COVNum = 0;

#if DEBUGASM
								std::cout << ". bb name :" << ss.str() << std::endl;
#endif
							}
							break;
						}
						case '#': {
							// bb
#if DEBUGASM
							std::cout << "sharp :" << std::endl;
#endif
							if (line.find("# %") < line.size()) {

								if (Function[Path][FunctionName]->BasicBlock.find(BasicBlockName)
										!= Function[Path][FunctionName]->BasicBlock.end()) {
									(Function[Path][FunctionName]->BasicBlock[BasicBlockName])->COVNum = COVNum;
								}

								ss.str("");
								for (unsigned long i = line.find('%') + 1; i < line.size(); i++) {
									if (line.at(i) == '%') {
										for (i++; i < line.size(); i++) {
											ss << line.at(i);
										}
										BasicBlockName = ss.str();

										Function[Path][FunctionName]->BasicBlock[BasicBlockName] = new DBasicBlock();
										(Function[Path][FunctionName]->BasicBlock[BasicBlockName])->name =
												BasicBlockName;
										COVNum = 0;

										break;
									}
								}

#if DEBUGASM
								std::cout << "# bb name :" << ss.str() << std::endl;
#endif
							}
							break;
						}
						case '	': {
#if DEBUGASM
							std::cout << "tab :" << std::endl;
							std::cout << "line.size() :" << line.size() << std::endl;
#endif
							if (line.size() == 1) {

							} else if (line.at(1) == '.') {
								if (Path.empty() && !FunctionName.empty() && line.find('#') < line.size()) {

									ss.str("");
									for (unsigned long i = line.find('#') + 2; i < line.find(':'); i++) {
										ss << line.at(i);
									}
									Path = ss.str();
#if DEBUGASM
									std::cout << "s Path :" << Path << std::endl;
#endif
									if ((Function.find(Path) != Function.end())
											&& (Function[Path].find(FunctionName) != Function[Path].end())) {

										if (Function[Path][FunctionName]->isObjudump()) {

										} else if (Function[Path][FunctionName]->isIR()) {

										} else if (Function[Path][FunctionName]->isAsmSourceCode()) {
											std::cout << "--------------------------------------------" << std::endl;
											std::cout << "s repeat Path :" << Path << std::endl;
											std::cout << "s repeat FunctionName :" << FunctionName << std::endl;

										}

									} else {
										Function[Path].insert(
												std::pair<std::string, DFunction *>(FunctionName, new DFunction()));
										Function[Path][FunctionName]->Name = FunctionName;
										Function[Path][FunctionName]->Path = Path;
										Function[Path][FunctionName]->setAsmSourceCode(true);

									}
								}
							} else if (line.at(1) == '#') {

							} else if (line.at(1) >= 'a' && line.at(1) <= 'z') {
								//asm instruction
								auto *inst = new DAInstruction();
								InsNum++;
								ss.str("");
								for (unsigned long i = 1; i < line.size(); i++) {
									ss << line.at(i);
								}
								inst->Inst = ss.str();
								inst->BasicBlockName = BasicBlockName;
								Function[Path][FunctionName]->InstASM.push_back(inst);
#if DEBUGASM
								std::cout << "inst :" << inst->Inst << std::endl;
#endif
								if (ss.str().find("call") <= ss.str().size()) {
									Function[Path][FunctionName]->CallInstNum++;
									if (ss.str().find("__sanitizer_cov_trace_pc") <= ss.str().size()) {
										(Function[Path][FunctionName]->BasicBlock[BasicBlockName])->COVNum++;
										COVNum++;
									}
								} else if (ss.str().find('j') <= ss.str().size()) {
									Function[Path][FunctionName]->JumpInstNum++;
								}
							}

							break;
						}
						case ' ': {
							//comment
#if DEBUGASM
							std::cout << "space :" << std::endl;
#endif
							break;
						}
						default: {
							if (line.find(':') < line.size()) {
								if (line.find("# @") < line.size()) {
									ss.str("");
									for (unsigned long i = 0; line.at(i) != ':'; i++) {
										ss << line.at(i);
									}
									FunctionName = ss.str();
									Path = "";
#if DEBUGASM
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

} /* namespace dra */
