/*
 * ModuleAll.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "ModuleAll.h"

#include <llvm/ADT/ilist.h>
#include <llvm/ADT/ilist_iterator.h>
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
#include <cstdlib>
#include <initializer_list>
#include <iostream>
#include <memory>
#include <utility>


namespace dra {

	ModuleAll::ModuleAll() {
		// TODO Auto-generated constructor stub

	}

	ModuleAll::~ModuleAll() {
		// TODO Auto-generated destructor stub
	}

	void ModuleAll::initializeModule(std::string InputFilename) {
		llvm::LLVMContext *cxts;
		llvm::SMDiagnostic Err;
		cxts = new llvm::LLVMContext[1];
		modules = llvm::parseIRFile(InputFilename, Err, cxts[0]);
		if (!modules) {
			std::cerr << "load module: " << InputFilename << " failed\n";
			exit(0);
		} else {
			llvm::Module *mm = modules.get();

			std::cerr << "size : " << mm->getNamedMDList().size() << "\n";
			for (llvm::ilist<llvm::NamedMDNode>::iterator i = mm->getNamedMDList().begin();
					i != mm->getNamedMDList().end(); i++) {
				(*i).dump();
			}

		}

	}


	void ModuleAll::set(llvm::Module* m) {
		this->m = m;
		for (llvm::Module::iterator it = m->begin(); it != m->end(); it++) {
			llvm::SmallVector<std::pair<unsigned, llvm::MDNode *>, 4> MDs;
			(*it).getAllMetadata(MDs);
			for (auto &MD : MDs) {
				if (llvm::MDNode *N = MD.second) {
					if (auto *SP = llvm::dyn_cast<llvm::DISubprogram>(N)) {
						llvm::StringRef file = SP->getFilename();
						std::string name = (*it).getName().str();
						std::string realname;
						if (name.find('.') < name.size()) {
							realname = name.substr(0, name.find('.'));
						} else {
							realname = name;
						}

						std::cout << "file : " << file.str() << std::endl;
						std::cout << "name : " << name << std::endl;
						std::cout << "realname : " << realname << std::endl;

						if (AllFunctionbc.find(file.str()) != AllFunctionbc.end()
								&& AllFunctionbc[file.str()].find(file.str()) != AllFunctionbc[file.str()].end()) {
							std::cout << "repeat" << std::endl;
						}
						AllFunctionbc[file.str()].insert(
								std::pair<std::string, FunctionAll*>(name, new FunctionAll()));
						AllFunctionbc[file.str()][name]->set(&*it);
					}
				}
			}

		}
	}

	void ModuleAll::setLine(std::vector<InformationOfSourceCode*> &IS) {
		for (auto it = AllFunctionbc.begin(), ie = AllFunctionbc.end(); it != ie; it++) {
			for (auto iit = (*it).second.begin(), iie = (*it).second.end(); iit != iie; iit++) {
				(*iit).second->setLine(IS);
			}
		}
	}

} /* namespace dra */
