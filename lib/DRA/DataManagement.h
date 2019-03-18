/*
 * DataManagement.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_DATAMANAGEMENT_H_
#define LIB_DRA_DATAMANAGEMENT_H_

#include <string>
#include <unordered_map>

#include "DAInstruction.h"

namespace dra {
	class DModule;
} /* namespace dra */


namespace dra {


    class DataManagement {
    public:
        DataManagement();

        virtual ~DataManagement();

        void initializeModule(std::string objdump, std::string AssemblySourceCode, std::string InputFilename);

        void initMap();
    public:
        dra::DModule *Modules;
        std::unordered_map<std::string, DAInstruction *> Map;

        unsigned int FindNum;
        unsigned int UnFindNum;
        unsigned int SameNum;
        unsigned int DiffNum;
    };

} /* namespace dra */

#endif /* LIB_DRA_DATAMANAGEMENT_H_ */
