/*
 * DFunction.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_FUNCTION_H_
#define LIB_DRA_FUNCTION_H_

#include <string>
#include <unordered_map>
#include <vector>

#include "DAInstruction.h"
#include "DBasicBlock.h"
#include "DLInstruction.h"

namespace dra {
	class DModule;
} /* namespace dra */

namespace llvm {
	class Function;
} /* namespace llvm */

namespace dra {

	class DFunction {
		public:
			DFunction();

			virtual ~DFunction();

			void InitIRFunction(llvm::Function *f);

			void setState(Kind kind);

			void update(Kind kind);

			bool isObjudump() const;

			void setObjudump(bool Objudump);

			bool isAsmSourceCode() const;

			void setAsmSourceCode(bool AsmSourceCode);

			bool isIR() const;

			void setIR(bool IR);

			bool isMap();

		public:
			bool Objudump;
			bool AsmSourceCode;
			bool IR;

			llvm::Function *function;
			DModule *parent;
			Kind state;

			std::string Name;
			std::string Path;

			std::string Address;
			unsigned int InstNum;
			unsigned int CallInstNum;
			unsigned int JumpInstNum;
			std::vector<DAInstruction *> InstASM;

			unsigned int BasicBlockNum;
			std::unordered_map<std::string, DBasicBlock *> BasicBlock;

	};

} /* namespace dra */

#endif /* LIB_DRA_FUNCTION_H_ */
