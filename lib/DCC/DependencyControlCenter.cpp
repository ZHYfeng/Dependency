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
#include "../RPC/DependencyRPC.pb.h"

namespace dra {

    DependencyControlCenter::DependencyControlCenter() {
        this->start_time = std::time(nullptr);
        std::cout << std::ctime(&this->start_time) << "*time : start_time" << std::endl;
    }

    DependencyControlCenter::~DependencyControlCenter() = default;

    void DependencyControlCenter::init(std::string objdump, std::string AssemblySourceCode, std::string InputFilename,
                                       const std::string &staticRes) {

        DM.initializeModule(std::move(objdump), std::move(AssemblySourceCode), std::move(InputFilename));
        this->current_time = std::time(nullptr);
        std::cout << std::ctime(&this->current_time) << "*time : initializeModule" << std::endl;

        //Deserialize the static analysis results.
        this->STA.initStaticRes(staticRes, &this->DM);
        this->current_time = std::time(nullptr);
        std::cout << std::ctime(&this->current_time) << "*time : initStaticRes" << std::endl;

        this->setRPCConnection();
    }

    void DependencyControlCenter::run() {
        while (true) {
            std::cout << "wait for get newInput" << std::endl;
            Input *newInput = client->GetNewInput();
            if (newInput != nullptr) {
                std::cout << "new input : " << newInput->sig() << std::endl;
                std::cout << newInput->prog() << std::endl;

                DInput *dInput = DM.getInput(newInput);
                std::cout << "dUncoveredAddress size : " << std::dec << dInput->dUncoveredAddress.size()
                          << std::endl;
                get_dependency_input(dInput);

                newInput->Clear();
                this->current_time = std::time(nullptr);
                std::cout << std::ctime(&this->current_time) << "*time : sleep_for 10s." << std::endl;
                std::this_thread::sleep_for(std::chrono::seconds(10));
            } else {
                this->current_time = std::time(nullptr);
                std::cout << std::ctime(&this->current_time) << "*time : sleep_for 60s." << std::endl;
                std::this_thread::sleep_for(std::chrono::seconds(60));
                setRPCConnection();
            }

            this->DM.dump_cover();
            this->DM.dump_uncover();

        }
    }


    void DependencyControlCenter::setRPCConnection() {
        this->client = new dra::DependencyRPCClient(
                grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials()));
        unsigned long long int vmOffsets = client->GetVmOffsets();
        DM.setVmOffsets(vmOffsets);
        this->current_time = std::time(nullptr);
        std::cout << std::ctime(&this->current_time) << "*time : GetVmOffsets" << std::endl;
    }

    void DependencyControlCenter::get_dependency_input(DInput *dInput) {

        Input *dependencyInput = new Input();
        dependencyInput->set_prog(dInput->prog);
        bool send_flag = false;

        for (auto u : dInput->dUncoveredAddress) {
            if (this->DM.check_uncovered_address(u)) {

                if (this->DM.uncover.find(u->address) != this->DM.uncover.end()) {
                    this->DM.uncover[u->address]->belong_to_Driver = true;
                }


                unsigned long long int address = DM.getSyzkallerAddress(u->address);
                unsigned long long int condition_address = DM.getSyzkallerAddress(u->condition_address);

                this->current_time = std::time(nullptr);
                std::cout << std::ctime(&current_time);
                std::cout << "condition trace_pc_address : " << std::hex << u->condition_address << "\n";
                std::cout << "uncovered trace_pc_address : " << std::hex << u->address << "\n";
                std::cout << "condition getSyzkallerAddress : " << std::hex << condition_address << "\n";
                std::cout << "uncovered getSyzkallerAddress : " << std::hex << address << "\n";

                UncoveredAddress *uncoveredAddress = dependencyInput->add_uncovered_address();
                uncoveredAddress->set_uncovered_address(address);
                uncoveredAddress->set_idx(u->idx);
                uncoveredAddress->set_condition_address(condition_address);

                sta::MODS *write_basicblock = get_write_basicblock(u);

                if (write_basicblock == nullptr) {

                } else {

                    send_flag = true;

                    if (this->DM.uncover.find(u->address) != this->DM.uncover.end()) {
                        this->DM.uncover[u->address]->related_to_gv = true;
                    }

                    for (auto &x : *write_basicblock) {

                        this->current_time = std::time(nullptr);
                        std::cout << std::ctime(&current_time);
                        std::cout << "write basicblock : " << std::endl;

                        dra::dump_inst(&x->B->front());

                        DBasicBlock *db = this->DM.get_DB_from_bb(x->B);
                        unsigned int write_address = DM.getSyzkallerAddress(db->trace_pc_address);


                        std::cout << "write_address getSyzkallerAddress : " << std::hex << write_address << "\n";
                        std::cout << "x->repeat : " << std::hex << x->repeat << "\n";
                        std::cout << "x->prio : " << std::hex << x->prio << "\n";
                        db->dump();

                        this->current_time = std::time(nullptr);
                        std::cout << std::ctime(&current_time);
                        std::vector<sta::cmd_ctx *> *cmd_ctx = x->get_cmd_ctx();
                        std::cout << "cmd size : " << std::dec << cmd_ctx->size() << "\n";
                        this->current_time = std::time(nullptr);
                        std::cout << std::ctime(&current_time);
                        for (auto c: *cmd_ctx) {
                            std::cout << "cmd dec: " << std::dec << c->cmd << "\n";
                            std::cout << "cmd hex: " << std::hex << c->cmd << "\n";
                            this->DM.dump_ctxs(&c->ctx);
                        }
                        this->current_time = std::time(nullptr);
                        std::cout << std::ctime(&current_time);

                        WriteAddress *writeAddress = uncoveredAddress->add_write_address();
                        writeAddress->set_write_address(write_address);
                        writeAddress->set_condition_address(condition_address);
                        writeAddress->set_repeat(x->repeat);
                        writeAddress->set_prio(x->prio);

                        auto function_name = "ioctl";
                        for (auto c : *cmd_ctx) {
                            auto write_syscall = writeAddress->add_write_syscall();
                            write_syscall->set_name(function_name);
                            write_syscall->set_number(c->cmd);
                        }
                        for (auto i : db->input) {
                            auto write_input = writeAddress->add_write_input();
                            write_input->set_sig(i.first->sig);
                        }

                        this->send_dependency_input(dependencyInput);
                    }
                }
            }
        }

        if (send_flag) {
            this->send_dependency_input(dependencyInput);
        }

    }

    void DependencyControlCenter::send_dependency_input(Input *dependencyInput) {
        std::cout << "SendDependencyInput sig : " << dependencyInput->sig() << std::endl;
        auto reply = client->SendDependencyInput(*dependencyInput);
#if DEBUG_RPC
        for (auto ua : dependencyInput->uncovered_address()) {
            std::cout << "uncover address : " << ua.address() << std::endl;
            std::cout << "uncover idx : " << ua.idx() << std::endl;
            std::cout << "uncover condition address : " << ua.condition_address() << std::endl;
            for (auto ra : ua.write_address()) {
                std::cout << "ra.address() : " << ra.address() << std::endl;
                std::cout << "ra.repeat() : " << ra.repeat() << std::endl;
                std::cout << "ra.prio() : " << ra.prio() << std::endl;
            }
        }
#endif
    }

    void DependencyControlCenter::test_sta() {
        auto f = this->DM.Modules->Function["block/blk-core.c"]["blk_flush_plug_list"];
        for (auto B : f->BasicBlock) {
            auto b = B.second->basicBlock;
            std::cout << "b name : " << B.second->name << std::endl;

            sta::MODS *allBasicblock = this->STA.GetAllGlobalWriteBBs(b, true);
            if (allBasicblock == nullptr) {
                // no taint or out side
                std::cout << "allBasicblock == nullptr" << std::endl;
            } else if (allBasicblock->empty()) {
                // unrelated to gv
                std::cout << "allBasicblock->size() == 0" << std::endl;
            } else if (!allBasicblock->empty()) {
                std::cout << "allBasicblock != nullptr && allBasicblock->size() != 0" << std::endl;

            }
        }

        exit(0);
    }

    void DependencyControlCenter::test_rpc() {

        exit(0);
    }

    void DependencyControlCenter::get_write_address() {

    }

    void DependencyControlCenter::send_write_address(WriteAddress *writeAddress) {

    }

    sta::MODS *DependencyControlCenter::get_write_basicblock(DUncoveredAddress *u) {

        //TODO: need to free "allBasicblock" and "cmds" to avoid memory leak,
        // or we can also set up a cache to avoid repeated query to STA.

        sta::MODS *res = nullptr;

        DBasicBlock *p = DM.Address2BB[u->condition_address]->parent;
        llvm::BasicBlock *b = dra::getFinalBB(p->basicBlock);

        this->current_time = std::time(nullptr);
        std::cout << std::ctime(&current_time);
        std::cout << "GetAllGlobalWriteBBs : " << std::endl;
        sta::MODS *write_basicblock = this->STA.GetAllGlobalWriteBBs(b, u->successor_idx);
        this->current_time = std::time(nullptr);
        std::cout << std::ctime(&current_time);

        p->dump();
        if (write_basicblock == nullptr) {
            // no taint or out side
            std::cout << "allBasicblock == nullptr" << std::endl;
        } else if (write_basicblock->size() == 0) {
            // unrelated to gv
            std::cout << "allBasicblock->size() == 0" << std::endl;
        } else if (!write_basicblock->empty()) {
            std::cout << "get useful static analysis result : " << std::dec << write_basicblock->size() << std::endl;
            res = write_basicblock;
        }

        return res;
    }

} /* namespace dra */
