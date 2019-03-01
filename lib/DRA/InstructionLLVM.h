/*
 * InstructionLLVM.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_INSTRUCTIONLLVM_H_
#define LIB_DRA_INSTRUCTIONLLVM_H_

#include <string>

namespace dra {
	class BasicBlockAll;
} /* namespace dra */


namespace llvm {
	class Instruction;
} /* namespace llvm */

namespace dra {

	enum Kind {
			other, untest, cover, uncover,
		};


	class InstructionLLVM {
		public:
			InstructionLLVM();
			virtual ~InstructionLLVM();
			void setState(Kind kind);
			void update(Kind kind);

		public:
			llvm::Instruction *i;
			Kind state;

			BasicBlockAll *parent;

			std::string FileName;
			unsigned int Line;


	};

} /* namespace dra */

#endif /* LIB_DRA_INSTRUCTIONLLVM_H_ */
