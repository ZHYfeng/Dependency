/*
 * BasicBlockAll.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_BASICBLOCKALL_H_
#define LIB_DRA_BASICBLOCKALL_H_

#include <string>
#include <vector>

#include "InstructionASM.h"
#include "InstructionLLVM.h"

namespace dra {
	class FunctionAll;
} /* namespace dra */
namespace llvm {
	class BasicBlock;
} /* namespace llvm */

namespace dra {

	class BasicBlockAll {
		public:
			BasicBlockAll();
			virtual ~BasicBlockAll();

			void set(llvm::BasicBlock *b);
			void setState(Kind kind);
			void update(Kind kind);

		public:

			llvm::BasicBlock *b;
			FunctionAll *parent;
			Kind state;
			std::string name;
			unsigned int covNum;

			std::vector<InstructionASM *> InstASM;
			std::vector<InstructionLLVM *> InstructionVector;
	};

} /* namespace dra */

#endif /* LIB_DRA_BASICBLOCKALL_H_ */
