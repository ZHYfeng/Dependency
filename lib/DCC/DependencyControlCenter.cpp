/*
 * DependencyControlCenter.cpp
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#include "DependencyControlCenter.h"
#include <chrono>
#include <thread>
#include <utility>
#include <grpcpp/grpcpp.h>
#include <llvm/IR/DebugInfoMetadata.h>
#include "../DRA/DModule.h"
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

        this->client = new dra::DependencyRPCClient(grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials()));
        unsigned long long int vmOffsets = client->GetVmOffsets();
        DM.setVmOffsets(vmOffsets);
        this->current_time = std::time(NULL);
        std::cout << std::ctime(&this->current_time) << "*time : GetVmOffsets" << std::endl;

        std::thread *t1 = new std::thread(this->record);

    }

    void DependencyControlCenter::run() {
        while (true) {
            std::cout << "wait for get newInput" << std::endl;
            NewInput *newInput = client->GetNewInput();
            if (newInput != nullptr) {
                std::cout << "get newInput size : " << newInput->input_size() << std::endl;
                for (int j = 0; j < newInput->input_size(); j++) {
                    const Input &input = newInput->input(j);
                    std::cout << "new input : " << input.sig() << std::endl;
                    DInput *dInput = DM.getInput(input);
                    DependencyInput dependencyInput;
                    bool sendFlag = false;
                    dependencyInput.set_sig(dInput->sig);
                    std::cout << "dUncoveredAddress size : " << dInput->dUncoveredAddress.size() << std::endl;
                    for (auto u : dInput->dUncoveredAddress) {
                        this->uncovered_address_number++;
                        if (this->DM.isDriver(u->address)) {
                            this->uncovered_address_number_driver++;
                            if (DM.Address2BB.find(u->condition_address) != DM.Address2BB.end()) {
                                auto *b = DM.Address2BB[u->condition_address]->parent->basicBlock;
                                sta::MODS *allBasicblock = this->STA.GetAllGlobalWriteBBs(DM.getFinalBB(b));
                                if (allBasicblock == nullptr) {
                                    // no taint or out side
                                    std::cout << "allBasicblock == nullptr" << std::endl;

                                } else if (allBasicblock->size() == 0) {
                                    // unrelated to gv
                                    std::cout << "allBasicblock->size() == 0" << std::endl;

                                } else if (allBasicblock != nullptr && allBasicblock->size() != 0) {
                                    this->uncovered_address_number_gv_driver++;
                                    sendFlag = true;
                                    std::cout << "get useful static analysis result" << std::endl;

                                    unsigned long long int address = DM.getSyzkallerAddress(u->address);
                                    unsigned long long int condition_address = DM.getSyzkallerAddress(u->condition_address);
                                    UncoveredAddress *uncoveredAddress = dependencyInput.add_uncovered_address();
                                    uncoveredAddress->set_address(address);
                                    uncoveredAddress->set_idx(u->idx);
                                    uncoveredAddress->set_condition_address(condition_address);

                                    for (auto &x : *allBasicblock) {
                                        llvm::BasicBlock *bb = DM.getRealBB(x->B);
                                        //Hang: NOTE: now let's just use "ioctl" as the "related syscall"
                                        //Hang: Below "cmds" is the value set for "cmd" arg of ioctl to reach this write BB.
                                        std::set<uint64_t> *cmds = x->getIoctlCmdSet();
                                        std::string Path = dra::DModule::getFileName(bb->getParent());
                                        std::string FunctionName = dra::DModule::getFunctionName(bb->getParent());
                                        std::string bbname = bb->getName().str();
                                        auto db = DM.Modules->Function[Path][FunctionName]->BasicBlock[bbname];
                                        unsigned int writeAddress = DM.getSyzkallerAddress(db->address);
                                        auto function_name = "ioctl";
                                        auto related_address = uncoveredAddress->add_related_address();
                                        related_address->set_address(writeAddress);
                                        std::cout << "cmds size : " << cmds->size() << std::endl;
                                        for (auto c : *cmds) {
                                            auto related_syscall = related_address->add_related_syscall();
                                            related_syscall->set_name(function_name);
                                            related_syscall->set_number(c);
                                        }
                                        for (auto i : db->input) {
                                            auto related_input = related_address->add_related_input();
                                            related_input->set_sig(i->sig);
                                        }

                                        //TODO: need to free "allBasicblock" and "cmds" to avoid memory leak,
                                        // or we can also set up a cache to avoid repeated query to STA.
                                    }
                                }
                            } else {
                                std::cerr << "can not find condition_address : " << std::hex << u->condition_address << std::endl;
                            }
                        } else {
                        }
                    }

                    if (sendFlag) {
                        std::cerr << "SendDependencyInput : " << std::endl;

                        auto reply = client->SendDependencyInput(dependencyInput);
                        if (reply->name() == "") {

                        } else {
                            std::cerr << "SendDependencyInput error : " << reply->name() << std::endl;
                        }
//                    std::cerr << "SendDependencyInput size : " << reply->address() << std::endl;
//                    std::cerr << "test GetDependencyInput : " << std::endl;
//                    auto neww = client->GetDependencyInput();
//                    for (int ni = 0; ni < neww->dependencyinput_size(); ni++) {
//                        const DependencyInput &nn = neww->dependencyinput(ni);
//                        std::cerr << "GetDependencyInput sig : " << nn.sig() << std::endl;
//                        std::cerr << "GetDependencyInput prog : " << nn.prog() << std::endl;
//                    }
                    }
                }
                newInput->Clear();
                this->current_time = std::time(NULL);
                std::cout << std::ctime(&this->current_time) << "*time : sleep_for 10s." << std::endl;
                std::this_thread::sleep_for(std::chrono::seconds(10));
            } else {
                this->current_time = std::time(NULL);
                std::cout << std::ctime(&this->current_time) << "*time : sleep_for 60s." << std::endl;
                std::this_thread::sleep_for(std::chrono::seconds(60));
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

            sta::MODS *allBasicblock = this->STA.GetAllGlobalWriteBBs(b);
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

    void DependencyControlCenter::test_rpc() {

//        std::cout << "test_rpc : " << std::endl;
//        DependencyInput dependencyInput;
//        dependencyInput.set_sig("dependencyInput.set_sig");
//        UncoveredAddress *uncoveredAddress = dependencyInput.add_uncovered_address();
//        uncoveredAddress->set_address(0xfffffff1);
//        uncoveredAddress->set_idx(2);
//        uncoveredAddress->set_condition_address(0xfffffff2);
//
//        auto function_name = "ioctl";
//        RelatedSyscall *relatedSyscall = uncoveredAddress->add_related_syscall();
//        relatedSyscall->set_address(0xfffffff3);
//        relatedSyscall->set_name(function_name);
//        relatedSyscall->set_number(0xfffffff4);
//
//        RelatedInput *relatedInput = uncoveredAddress->add_related_input();
//        relatedInput->set_address(0xfffffff5);
//        relatedInput->set_sig("sig");
//
//        std::cout << "SendDependencyInput : " << std::endl;
//        auto r = client->SendDependencyInput(dependencyInput);
//        std::cout << "SendDependencyInput.r : " << r->name() << std::endl;
//        for (int j = 0; j < dependencyInput.uncovered_address_size(); j++) {
//            auto uu = dependencyInput.uncovered_address(j);
//            for (int k = 0; k < uu.related_input_size(); k++) {
//                auto ii = uu.related_input(k);
//                std::cout << "ii.sig : " << ii.sig() << std::endl;
//                std::cout << "ii.address : " << ii.address() << std::endl;
//            }
//
//            for (auto ss: uu.related_syscall()) {
//                std::cout << "ss.number : " << ss.number() << std::endl;
//                std::cout << "ss.name : " << ss.name() << std::endl;
//                std::cout << "ss.address : " << ss.address() << std::endl;
//            }
//        }
//
//        std::cout << "GetDependencyInput : " << std::endl;
//        auto newD = client->GetDependencyInput();
//        std::cout << "newD.dependencyinput_size : " << newD->dependencyinput_size() << std::endl;
//        for (int i = 0; i < newD->dependencyinput_size(); i++) {
//            auto dd = newD->dependencyinput(i);
//            std::cout << "dd.sig : " << dd.sig() << std::endl;
//            for (int j = 0; j < dd.uncovered_address_size(); j++) {
//                auto uu = dd.uncovered_address(j);
//                for (int k = 0; k < uu.related_input_size(); k++) {
//                    auto ii = uu.related_input(k);
//                    std::cout << "ii.sig : " << ii.sig() << std::endl;
//                    std::cout << "ii.address : " << std::hex << ii.address() << std::endl;
//                }
//
//                for (auto ss: uu.related_syscall()) {
//                    std::cout << "ss.number : " << ss.number() << std::endl;
//                    std::cout << "ss.name : " << ss.name() << std::endl;
//                    std::cout << "ss.address : " << ss.address() << std::endl;
//                }
//            }
//        }

        exit(0);
    }

} /* namespace dra */
