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

#include "DModule.h"


namespace dra {


    class DataManagement {
    public:
        DataManagement();

        virtual ~DataManagement();

        void initializeModule(std::string objdump, std::string AssemblySourceCode, std::string InputFilename);

        void Statistics();

    public:
        dra::DModule *Modules;

        unsigned int FindNum;
        unsigned int UnFindNum;
        unsigned int SameNum;
        unsigned int DiffNum;
    };

} /* namespace dra */

#endif /* LIB_DRA_DATAMANAGEMENT_H_ */
