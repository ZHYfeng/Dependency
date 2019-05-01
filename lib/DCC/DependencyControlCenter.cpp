#include <utility>

#include <utility>

#include <utility>

/*
 * DependencyControlCenter.cpp
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#include "DependencyControlCenter.h"

#include <grpcpp/grpcpp.h>

namespace dra {

    DependencyControlCenter::DependencyControlCenter() :
            client(grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials())) {

    }

    DependencyControlCenter::~DependencyControlCenter() = default;

    void DependencyControlCenter::init(std::string objdump, std::string AssemblySourceCode, std::string InputFilename) {
        DM.initializeModule(std::move(objdump), std::move(AssemblySourceCode), std::move(InputFilename));
        unsigned long long int vmOffsets = client.GetVmOffsets();
        DM.setVmOffsets(vmOffsets);
    }

    void DependencyControlCenter::run() {
        while (true) {
            NewInput *newInput = client.GetNewInput();
            if (newInput != nullptr) {
                for (int j = 0; j < newInput->input_size(); j++) {
                    const Input &input = newInput->input(j);
                    DInput *dInput = new DInput();
                    // TODO(Yu): set input coverage and get uncover address
                    DependencyInput dependencyInput;
                    for (auto u : dInput->dUncoveredAddress) {
                        unsigned long long int address = u.address;
                        unsigned long long int condition_address = u.condition_address;

                        UncoveredAddress *uncoveredAddress = dependencyInput.add_uncovered_address();
                        uncoveredAddress->set_address(address);
                        uncoveredAddress->set_condition_address(condition_address);

                        llvm::BasicBlock *b = DM.Address2BB[condition_address]->parent->basicBlock;
                        // TODO(hang): GetGlobalWriteBB
                        auto allbb = GetGlobalWriteBB(b);
                        for (auto bb : allbb) {
                            unsigned long long int writeAddress = DM.Modules->Function[bb.path][bb.name].address;
                            // TODO(hang): GetGlobalWriteBB
                            auto relatedsyscall = GetRelatedSyscall(bb);
                            RelatedSyscall *relatedSyscall = uncoveredAddress->add_related_syscall();
                            relatedSyscall->set_address(writeAddress);
                            relatedSyscall->set_name(relatedsyscall);
                        }
                    }
                    // TODO(Yu): set dependencyInput
                    client.SendDependencyInput(dependencyInput);
                }
            }

        }
    }

} /* namespace dra */
