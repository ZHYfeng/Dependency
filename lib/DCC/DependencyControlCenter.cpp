/*
 * DependencyControlCenter.cpp
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#include "DependencyControlCenter.h"

#include <utility>
#include <grpcpp/grpcpp.h>
#include <fstream>

namespace dra {

    DependencyControlCenter::DependencyControlCenter() :
            client(grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials())) {

    }

    DependencyControlCenter::~DependencyControlCenter() = default;

    void DependencyControlCenter::init(std::string objdump, std::string AssemblySourceCode, std::string InputFilename, std::string staticRes) {
        DM.initializeModule(std::move(objdump), std::move(AssemblySourceCode), std::move(InputFilename));
        unsigned long long int vmOffsets = client.GetVmOffsets();
        DM.setVmOffsets(vmOffsets);
        //Deserialize the static analysis results.
        this->initStaticRes(staticRes);
    }

    int DependencyControlCenter::initStaticRes(std::string staticRes) {
        try{
            std::ifstream infile;
            infile.open(staticRes);
            infile >> this->j_taintedBrs >> this->j_analysisCtxMap >> this->j_tagMap >> this->j_modInstCtxMap;
            infile.close();
            return 0;
        }catch(){
            std::cout << "Fail to deserialize the static analysis results!\n";
        }
        return 1;
    }

    LOC_INF *DependencyControlCenter::getLocInf(llvm::Instruction* I) {
        if(!I){
            return nullptr;
        }
        std::string inst,bb,func,mod;
        std::string str;
        llvm::raw_string_ostream ss(str);
        ss << *I;
        inst = ss.str();
        if(I->getParent()){
            bb = I->getParent()->getName().str();
        }
        if(I->getFunction()){
            func = I->getFunction()->getName().str();
        }
        if(I->getModule()){
            mod = I->getModule()->getName().str();
        }
        LOC_INF *loc_inf = new LOC_INF;
        loc_inf->push_back(inst);
        loc_inf->push_back(bb);
        loc_inf->push_back(func);
        loc_inf->push_back(mod);
        return loc_inf;
    }

    LOC_INF *DependencyControlCenter::getLocInf(llvm::BasicBlock* B) {
        if(!B){
            return nullptr;
        }
        return this->getLocInf(B->begin());
    }

    void DependencyControlCenter::run() {
        while (true) {
            NewInput *newInput = client.GetNewInput();
            if (newInput != nullptr) {
                for (int j = 0; j < newInput->input_size(); j++) {
                    const Input &input = newInput->input(j);
                    DInput *dInput = DM.getInput(input);
                    // TODO(Yu): set input coverage and get uncover address
                    DependencyInput dependencyInput;
                    for (auto u : dInput->dUncoveredAddress) {
                        unsigned long long int address = DM.getSyzkallerAddress(u->address);
                        unsigned long long int condition_address = DM.getSyzkallerAddress(u->condition_address);

                        UncoveredAddress *uncoveredAddress = dependencyInput.add_uncovered_address();
                        uncoveredAddress->set_address(address);
                        uncoveredAddress->set_idx(u->idx);
                        uncoveredAddress->set_condition_address(condition_address);

                        llvm::BasicBlock *b = DM.Address2BB[condition_address]->parent->basicBlock;
                        // TODO(hang): GetGlobalWriteBB
                        auto allbb = GetGlobalWriteBB(b);
                        for (auto bb : allbb) {
                            auto db = DM.Modules->Function[bb.path][bb.name];
                            unsigned long long int writeAddress = db.address;

                            // TODO(hang): GetGlobalWriteBB
                            auto relatedsyscall = GetRelatedSyscall(bb);
                            RelatedSyscall *relatedSyscall = uncoveredAddress->add_related_syscall();
                            relatedSyscall->set_address(writeAddress);
                            relatedSyscall->set_name(relatedsyscall);

                            RelatedInput *relatedInput = uncoveredAddress->add_related_input();
                            relatedInput->set_address(writeAddress);
                            for(auto i : db->input){
                                relatedInput->set_sig(i->sig);
                            }
                        }
                    }
                    client.SendDependencyInput(dependencyInput);
                }
            }

        }
    }

} /* namespace dra */
