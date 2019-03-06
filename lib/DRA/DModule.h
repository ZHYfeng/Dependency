/*
 * DModule.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_MODULEALL_H_
#define LIB_DRA_MODULEALL_H_

#include <string>
#include <unordered_map>

#include "DFunction.h"

#define DEBUG 1
#define DEBUGOBJDUMP 0
#define DEBUGASM 0

namespace llvm {
    class Module;
} /* namespace llvm */

namespace dra {

    class DModule {
    public:
        DModule();

        virtual ~DModule();



        std::string exec(std::string cmd);

        void ReadObjdump(std::string objdump);

        void ReadAsmSourceCode(std::string AssemblySourceCode);

        void ReadBC(std::string InputFilename);

        void BuildLLVMFunction(llvm::Module *Module);

    public:

        std::unordered_map<std::string, std::unordered_map<std::string, DFunction *>> Function;

    };

} /* namespace dra */

#endif /* LIB_DRA_MODULEALL_H_ */
