/*
 * ModuleAll.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_MODULEALL_H_
#define LIB_DRA_MODULEALL_H_

#include <memory>
#include <string>
#include <unordered_map>
#include <vector>

#include "FunctionAll.h"
#include "InformationOfSourceCode.h"

namespace llvm {
	class Module;
} /* namespace llvm */

#define PATH_SIZE 10000

namespace dra {

	class ModuleAll {
		public:
			ModuleAll();
			virtual ~ModuleAll();
			void initializeModule(std::string InputFilename);
			void set(llvm::Module* m);
			void setLine(std::vector<InformationOfSourceCode*> &IS);

		public:
			std::unique_ptr<llvm::Module> modules;
			llvm::Module* m;
			std::unordered_map<std::string, std::unordered_map<std::string, FunctionAll *>> AllFunctionbc;
	};

} /* namespace dra */

#endif /* LIB_DRA_MODULEALL_H_ */
