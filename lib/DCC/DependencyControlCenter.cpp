/*
 * DependencyControlCenter.cpp
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#include "DependencyControlCenter.h"

#include <thread>
#include <utility>
#include <grpcpp/grpcpp.h>
#include <llvm/IR/DebugInfoMetadata.h>

#include "../DRA/DFunction.h"

namespace dra {

    DependencyControlCenter::DependencyControlCenter() {
        this->start_time = std::time(nullptr);
        std::cout << std::ctime(&this->start_time) << "*time : start_time" << std::endl;

        this->uncovered_address_number = 0;
        this->uncovered_address_number_driver = 0;
        this->uncovered_address_number_gv_driver = 0;
    }

    DependencyControlCenter::~DependencyControlCenter() = default;

    void DependencyControlCenter::init(std::string objdump, std::string AssemblySourceCode, std::string InputFilename, const std::string &staticRes) {


        DM.initializeModule(std::move(objdump), std::move(AssemblySourceCode), std::move(InputFilename));
        this->current_time = std::time(NULL);
        std::cout << std::ctime(&this->current_time) << "*time : initializeModule" << std::endl;

        //Deserialize the static analysis results.
        this->STA.initStaticRes(staticRes, &this->DM);
        this->current_time = std::time(NULL);
//        this->current_time
        std::cout << std::ctime(&this->current_time) << "*time : initStaticRes" << std::endl;

        this->test_sta();

        this->client = new dra::DependencyRPCClient(grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials()));
        unsigned long long int vmOffsets = client->GetVmOffsets();
        DM.setVmOffsets(vmOffsets);
        this->current_time = std::time(NULL);
        std::cout << std::ctime(&this->current_time) << "*time : GetVmOffsets" << std::endl;

        std::thread *t1 = new std::thread(this->record);

    }

    void DependencyControlCenter::run() {
        while (true) {
            NewInput *newInput = client->GetNewInput();
            if (newInput != nullptr) {
                for (int j = 0; j < newInput->input_size(); j++) {
                    const Input &input = newInput->input(j);
                    DInput *dInput = DM.getInput(input);
                    DependencyInput dependencyInput;
                    for (auto u : dInput->dUncoveredAddress) {

                        this->uncovered_address_number++;
                        if (this->DM.isDriver(u->address)) {
                            this->uncovered_address_number_driver++;

                            unsigned long long int address = DM.getSyzkallerAddress(u->address);
                            unsigned long long int condition_address = DM.getSyzkallerAddress(u->condition_address);

                            UncoveredAddress *uncoveredAddress = dependencyInput.add_uncovered_address();
                            uncoveredAddress->set_address(address);
                            uncoveredAddress->set_idx(u->idx);
                            uncoveredAddress->set_condition_address(condition_address);

                            if (DM.Address2BB.find(condition_address) != DM.Address2BB.end()) {
                                auto *b = DM.Address2BB[condition_address]->parent->basicBlock;
                                b.
                                MOD_BBS *allBasicblock = this->STA.GetAllGlobalWriteBBs(DM.getFinalBB(b));
                                if (allBasicblock == nullptr) {
                                    // no taint or out side

                                } else if (allBasicblock->size() == 0) {
                                    // unrelated to gv

                                } else if (allBasicblock != nullptr && allBasicblock->size() != 0) {
                                    this->uncovered_address_number_gv_driver++;

                                    for (auto &x : *allBasicblock) {
                                        llvm::BasicBlock *bb = DM.getRealBB(x.first);
                                        MOD_INF &mod_inf = x.second;
                                        //Hang: NOTE: now let's just use "ioctl" as the "related syscall"
                                        //Hang: Below "cmds" is the value set for "cmd" arg of ioctl to reach this write BB.
                                        std::set<uint64_t> *cmds = this->STA.getIoctlCmdSet(&mod_inf);

                                        llvm::SmallVector<std::pair<unsigned, llvm::MDNode *>, 4> MDs;
                                        bb->getParent()->getAllMetadata(MDs);
                                        for (auto &MD : MDs) {
                                            if (llvm::MDNode *N = MD.second) {
                                                if (auto *SP = llvm::dyn_cast<llvm::DISubprogram>(N)) {
                                                    std::string Path = SP->getFilename().str();
                                                    std::string name = bb->getParent()->getName().str();
                                                    std::string FunctionName;
                                                    if (name.find('.') < name.size()) {
                                                        FunctionName = name.substr(0, name.find('.'));
                                                    } else {
                                                        FunctionName = name;
                                                    }
                                                    std::string bbname = bb->getName().str();
                                                    auto db = DM.Modules->Function[Path][FunctionName]->BasicBlock[bbname];
                                                    unsigned long long int writeAddress = db->address;

                                                    for (auto c : *cmds) {
                                                        auto function_name = "ioctl";
                                                        RelatedSyscall *relatedSyscall = uncoveredAddress->add_related_syscall();
                                                        relatedSyscall->set_address(writeAddress);
                                                        relatedSyscall->set_name(function_name);
                                                        relatedSyscall->set_number(c);
                                                    }

                                                    for (auto i : db->input) {
                                                        RelatedInput *relatedInput = uncoveredAddress->add_related_input();
                                                        relatedInput->set_address(writeAddress);
                                                        relatedInput->set_sig(i->sig);
                                                    }

                                                    client->SendDependencyInput(dependencyInput);
                                                }
                                            }
                                        }
                                        //TODO: need to free "allBasicblock" and "cmds" to avoid memory leak, or we can also set up a cache to avoid repeated query to STA.
                                    }
                                }
                            }
                        } else {

                        }
                    }
                }
            }
        }
    }

    void DependencyControlCenter::record() {
        std::chrono::seconds timespan(1800);
        while (true) {
            std::this_thread::sleep_for(timespan);

        }
    }

    void DependencyControlCenter::test_sta() {
        auto f = this->DM.Modules->Function["block/blk-core.c"]["blk_flush_plug_list"];
        for (auto B : f->BasicBlock) {
            auto b = B.second->basicBlock;
            std::cout << "b name : " << B.second->name << std::endl;

            MOD_BBS *allBasicblock = this->STA.GetAllGlobalWriteBBs(b);
            if (allBasicblock == nullptr) {
                // no taint or out side
                std::cout << "allBasicblock == nullptr" << std::endl;
            } else if (allBasicblock->size() == 0) {
                // unrelated to gv
                std::cout << "allBasicblock->size() == 0" << std::endl;
            } else if (allBasicblock != nullptr && allBasicblock->size() != 0) {
                std::cout << "allBasicblock != nullptr && allBasicblock->size() != 0" << std::endl;

            }
        }

        exit(0);
    }

} /* namespace dra */
