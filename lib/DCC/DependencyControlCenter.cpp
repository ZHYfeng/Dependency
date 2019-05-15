/*
 * DependencyControlCenter.cpp
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#include "DependencyControlCenter.h"

#include <utility>
#include <grpcpp/grpcpp.h>
#include <llvm/IR/DebugInfoMetadata.h>

#include "../DRA/DFunction.h"

namespace dra {

    DependencyControlCenter::DependencyControlCenter() :
            client(grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials())) {

    }

    DependencyControlCenter::~DependencyControlCenter() = default;

    void DependencyControlCenter::init(std::string objdump, std::string AssemblySourceCode, std::string InputFilename, const std::string &staticRes) {
        DM.initializeModule(std::move(objdump), std::move(AssemblySourceCode), std::move(InputFilename));
        unsigned long long int vmOffsets = client.GetVmOffsets();
        DM.setVmOffsets(vmOffsets);
        //Deserialize the static analysis results.
        this->STA.initStaticRes(staticRes, (DM.Modules->module).get());
    }

    void DependencyControlCenter::run() {
        while (true) {
            NewInput *newInput = client.GetNewInput();
            if (newInput != nullptr) {
                for (int j = 0; j < newInput->input_size(); j++) {
                    const Input &input = newInput->input(j);
                    DInput *dInput = DM.getInput(input);
                    DependencyInput dependencyInput;
                    for (auto u : dInput->dUncoveredAddress) {
                        unsigned long long int address = DM.getSyzkallerAddress(u->address);
                        unsigned long long int condition_address = DM.getSyzkallerAddress(u->condition_address);

                        UncoveredAddress *uncoveredAddress = dependencyInput.add_uncovered_address();
                        uncoveredAddress->set_address(address);
                        uncoveredAddress->set_idx(u->idx);
                        uncoveredAddress->set_condition_address(condition_address);

                        if (DM.Address2BB.find(condition_address) != DM.Address2BB.end()) {
                            llvm::BasicBlock *b = DM.Address2BB[condition_address]->parent->basicBlock;

                            MOD_BBS *allBasicblock = this->STA.GetAllGlobalWriteBBs(b);
                            if (allBasicblock != nullptr && allBasicblock->size() != 0) {
                                for (auto &x : *allBasicblock) {
                                    llvm::BasicBlock *bb = x.first;
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
                                            }
                                        }
                                    }
                                    //TODO: need to free "allBasicblock" and "cmds" to avoid memory leak, or we can also set up a cache to avoid repeated query to STA.
                                }
                            } else if (allBasicblock == nullptr) {
                                // no taint or out side

                            } else if (allBasicblock->size() == 0) {

                            }
                        }
                        client.SendDependencyInput(dependencyInput);
                    }
                }
            }
        }
    }

} /* namespace dra */
