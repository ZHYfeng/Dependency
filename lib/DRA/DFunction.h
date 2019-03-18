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
#include "DModule.h"

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
			void setState(CoverKind kind);
			void update(CoverKind kind);

			bool isObjudump() const;
			void setObjudump(bool Objudump);
			bool isAsmSourceCode() const;
			void setAsmSourceCode(bool AsmSourceCode);
			bool isIR() const;
			void setIR(bool IR);
			void setKind(FunctionKind kind);

			bool isMap();

			static DFunction MargeDFunction(DFunction *one, DFunction *two);

			bool isRepeat() const;
			void setRepeat(bool repeat);

			void dump();

		public:
			bool Objudump;
			bool AsmSourceCode;
			bool IR;

			bool repeat;

			llvm::Function *function;
			DModule *parent;
			CoverKind state;

			std::string FunctionName;
			std::string IRName;
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
