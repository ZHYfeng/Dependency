/*
 * InstructionASM.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_INSTRUCTIONASM_H_
#define LIB_DRA_INSTRUCTIONASM_H_

#include <string>

#include "InstructionLLVM.h"

namespace dra {

	class InstructionASM {
		public:
			InstructionASM();
			virtual ~InstructionASM();

		public:
			Kind state;
			std::string Inst;
			std::string BasicBlockName;
			BasicBlockAll *parent;
			std::string Address;
	};

} /* namespace dra */

#endif /* LIB_DRA_INSTRUCTIONASM_H_ */
