/*
 * DataASM.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_DATAASM_H_
#define LIB_DRA_DATAASM_H_

#include <string>
#include <unordered_map>

#include "FunctionAll.h"

#define DEBUG 1

namespace dra {

	class DataASM {
		public:
			DataASM();
			virtual ~DataASM();
			std::string exec(std::string cmd);

			void ReadFromObjdump(std::string objdump);
			void ReadFromAsmSourceCode(std::string AssemblySourceCode);
			void Statistics();
			void Analysis(std::string AssemblySourceCode, std::string objdump);

		public:
			std::unordered_map<std::string, std::unordered_map<std::string, FunctionAll*>> AllFunctiono;
			std::unordered_map<std::string, std::unordered_map<std::string, FunctionAll*>> AllFunctions;
			unsigned int FindNum, UnFindNum;
			unsigned int SameNum, DiffNum;
	};

} /* namespace dra */

#endif /* LIB_DRA_DATAASM_H_ */
