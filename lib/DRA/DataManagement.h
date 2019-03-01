/*
 * DataManagement.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_DATAMANAGEMENT_H_
#define LIB_DRA_DATAMANAGEMENT_H_

#include <memory>
#include <string>
#include <vector>

#include "FunctionAll.h"
#include "InformationOfSourceCode.h"
#include "ModuleAll.h"

namespace dra {
	class DataASM;
} /* namespace dra */
namespace llvm {
	class Module;
} /* namespace llvm */

namespace dra {


	class DataManagement {
		public:
			DataManagement();
			virtual ~DataManagement();

			int GetInformationFromCoverFile(std::string CoverFileName);
			void initializeModule(std::string InputFilename);
			void setState();
			void MapBBfromStoBC();

		public:
			dra::ModuleAll *Modules;

			std::vector<InformationOfSourceCode*> allSourceCode;
			dra::DataASM *DataASM;
	};

} /* namespace dra */

#endif /* LIB_DRA_DATAMANAGEMENT_H_ */
