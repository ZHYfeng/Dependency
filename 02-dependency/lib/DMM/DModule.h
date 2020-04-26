/*
 * DModule.h
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_MODULEALL_H_
#define LIB_DRA_MODULEALL_H_

#include <memory>
#include <string>
#include <unordered_map>

#include "DFunction.h"

namespace dra {
    class address;
} /* namespace dra */

#define TEST 0
#define DEBUG_MAP 0
#define DEBUG_BC 0
#define DEBUG_OBJ_DUMP 0
#define DEBUG_ASM 0

namespace llvm {
    class Module;
} /* namespace llvm */

namespace dra {

    class DModule {
    public:
        DModule();

        virtual ~DModule();

        std::string exec(std::string cmd);

        void ReadObjdump(const std::string& objdump);

        void ReadAsmSourceCode(std::string AssemblySourceCode);

        void ReadBC(const std::string& InputFilename);

        void BuildLLVMFunction(llvm::Module *Module);

        void AddRepeatFunction(DFunction *function, FunctionKind kind);

        DFunction *CheckRepeatFunction(std::string Path, std::string FunctionName, FunctionKind kind);

        DFunction *CreatFunction(std::string Path, std::string FunctionName, FunctionKind kind);

        DFunction *get_DF_from_f(llvm::Function *b);

        DBasicBlock *get_DB_from_bb(llvm::BasicBlock *b);

        DBasicBlock *get_DB_from_i(llvm::Instruction *i);

        void add_number_basic_block_covered();

    public:
        std::unique_ptr<llvm::Module> module;
        std::unordered_map<std::string, std::unordered_map<std::string, DFunction *>> Function;
        std::unique_ptr<dra::address> addr2line;

        uint32_t NumberBasicBlock;
        uint32_t NumberBasicBlockReal;
        uint32_t NumberBasicBlockCovered;

        std::unordered_map<std::string, DFunction *> RepeatBCFunction;
        std::unordered_map<std::string, DFunction *> RepeatOFunction;
        std::unordered_map<std::string, std::unordered_map<std::string, DFunction *>> RepeatSFunction;

    };

} /* namespace dra */

#endif /* LIB_DRA_MODULEALL_H_ */
