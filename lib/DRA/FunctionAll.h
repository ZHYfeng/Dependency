/*
 * FunctionAll.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_FUNCTIONALL_H_
#define LIB_DRA_FUNCTIONALL_H_

#include <string>
#include <vector>

#include "BasicBlockAll.h"
#include "InformationOfSourceCode.h"
#include "InstructionASM.h"
#include "InstructionLLVM.h"

namespace dra {
	class ModuleAll;
} /* namespace dra */

namespace llvm {
	class Function;
} /* namespace llvm */

namespace dra {

	class FunctionAll {
		public:
			FunctionAll();
			virtual ~FunctionAll();
			void set(llvm::Function *f);
			void setState(Kind kind);
			void update(Kind kind);

		public:
			llvm::Function *f;
			ModuleAll *parent;
			Kind state;



			std::string Name;
			std::string Path;

			std::string Address;
			unsigned int InstNum;
			unsigned int CallInstNum;
			unsigned int JumpInstNum;
			std::vector<InstructionASM *> InstASM;

			unsigned int BasicBlockNum;
			std::vector<BasicBlockAll *> BasicBlockVector;
	};

} /* namespace dra */

#endif /* LIB_DRA_FUNCTIONALL_H_ */
