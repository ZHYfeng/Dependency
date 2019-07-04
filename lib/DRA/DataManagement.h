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
#include "DFunction.h"
#include "DInput.h"
#include "DModule.h"
#include "../RPC/DependencyRPCClient.h"
//#include "../RPC/Data.pb.h"

namespace dra {

    class coverage {
    public:
        std::time_t time;
        unsigned long long int address;
    };

    class uncover_info {
    public:
        uncover_info();

    public:
        std::time_t time;
        unsigned long long int address;
        unsigned long long int condition_address;
        bool belong_to_Driver;
        bool related_to_gv;
        bool covered;
        bool covered_by_dependency;
    };

    llvm::BasicBlock *getRealBB(llvm::BasicBlock *b);

    llvm::BasicBlock *getFinalBB(llvm::BasicBlock *b);

    void dump_inst(llvm::Instruction *inst);

    class DataManagement {
    public:
        DataManagement();

        virtual ~DataManagement();

        void initializeModule(std::string objdump, std::string AssemblySourceCode, std::string InputFilename);

        void BuildAddress2BB();

        void getVmOffsets(std::string vmOffsets);

        void setVmOffsets(unsigned long long int vmOffsets);

        void getInput(std::string coverfile);

        DInput *getInput(Input input);

        void setInput();

        unsigned long long int getRealAddress(unsigned long long int address);

        unsigned long long int getSyzkallerAddress(unsigned long long int address);

        bool isDriver(unsigned long long int address);

        void dump_address(unsigned long long int address);

        void dump_cover();

        void dump_uncover();

        void dump_ctxs(std::vector<llvm::Instruction *> *ctx);



    public:
        dra::DModule *Modules;
        std::unordered_map<unsigned long long int, DAInstruction *> Address2BB;
        std::unordered_map<std::string, DInput *> Inputs;
//        dra::all_data Add_Data;
        std::map<unsigned long long int, std::time_t> cover;
        std::map<unsigned long long int, uncover_info *> uncover;
        std::vector<coverage *> time;
        unsigned long long int vmOffsets;

    };

} /* namespace dra */

#endif /* LIB_DRA_DATAMANAGEMENT_H_ */
