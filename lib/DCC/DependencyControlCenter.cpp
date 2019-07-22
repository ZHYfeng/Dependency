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
#include "general.h"

namespace dra {

    DependencyControlCenter::DependencyControlCenter() {
        dra::outputTime("start_time");
        dra::deal_sig();
    }

    DependencyControlCenter::~DependencyControlCenter() = default;

    void DependencyControlCenter::init(std::string objdump, std::string AssemblySourceCode, std::string InputFilename,
                                       const std::string &staticRes) {

        DM.initializeModule(std::move(objdump), std::move(AssemblySourceCode), std::move(InputFilename));
        dra::outputTime("initializeModule");



        //Deserialize the static analysis results.
        this->STA.initStaticRes(staticRes, &this->DM);

        this->setRPCConnection();
    }

    void DependencyControlCenter::run() {
        while (true) {
            std::cout << "wait for get newInput" << std::endl;
            Inputs *newInput = client->GetNewInput();
            if (newInput != nullptr) {
                for (auto &input : *newInput->mutable_input()) {
                    std::cout << "new input : " << input.sig() << std::endl;
                    std::cout << input.program() << std::endl;

                    DInput *dInput = DM.getInput(&input);
                    get_dependency_input(dInput);
                }
                newInput->Clear();

                outputTime("sleep_for 10s");
                std::this_thread::sleep_for(std::chrono::seconds(10));
            } else {
                outputTime("sleep_for 60s");
                std::this_thread::sleep_for(std::chrono::seconds(60));
                setRPCConnection();
            }

            this->DM.dump_cover();
            this->DM.dump_uncover();

            get_write_addresses();
        }
    }


    void DependencyControlCenter::setRPCConnection() {
        this->client = new dra::DependencyRPCClient(
                grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials()));
        unsigned long long int vmOffsets = client->GetVmOffsets();
        DM.setVmOffsets(vmOffsets);
        dra::outputTime("GetVmOffsets");
    }

    void DependencyControlCenter::get_dependency_input(DInput *dInput) {


        std::cout << "dUncoveredAddress size : " << std::dec << dInput->dUncoveredAddress.size()
                  << std::endl;
        uint64_t i = 0;
        for (auto u : dInput->dUncoveredAddress) {
            i++;
            outputTime("uncovered address count : " + std::to_string(i));

            if (this->DM.check_uncovered_address(u)) {

                if (this->DM.uncover.find(u->uncovered_address()) != this->DM.uncover.end()) {
                    this->DM.uncover[u->uncovered_address()]->belong_to_Driver = true;
                }

                sta::MODS *write_basicblock = get_write_basicblock(u);
                if (write_basicblock == nullptr) {

                } else {

                    Input *dependencyInput = new Input();
                    dependencyInput->set_sig(dInput->sig);
                    dependencyInput->set_program(dInput->program);

                    unsigned long long int syzkallerUncoveredAddress = DM.getSyzkallerAddress(u->uncovered_address());
                    unsigned long long int syzkallerConditionAddress = DM.getSyzkallerAddress(u->condition_address());

                    outputTime("");
                    std::cout << "condition trace_pc_address : " << std::hex << u->condition_address() << "\n";
                    std::cout << "uncovered trace_pc_address : " << std::hex << u->uncovered_address() << "\n";
                    std::cout << "condition getSyzkallerAddress : " << std::hex << syzkallerConditionAddress << "\n";
                    std::cout << "uncovered getSyzkallerAddress : " << std::hex << syzkallerUncoveredAddress << "\n";

                    UncoveredAddress *uncoveredAddress = dependencyInput->add_uncovered_address();
                    uncoveredAddress->set_uncovered_address(syzkallerUncoveredAddress);
                    uncoveredAddress->set_condition_address(syzkallerConditionAddress);

                    set_runtime_data(uncoveredAddress->mutable_run_time_date(), dependencyInput->program(), u->idx(),
                                     syzkallerConditionAddress, syzkallerUncoveredAddress);


                    if (this->DM.uncover.find(u->uncovered_address()) != this->DM.uncover.end()) {
                        this->DM.uncover[u->uncovered_address()]->related_to_gv = true;
                    }

                    for (auto &x : *write_basicblock) {
                        WriteAddress *writeAddress = uncoveredAddress->add_write_address();
                        get_write_address(x, u, writeAddress);

                        set_runtime_data(writeAddress->mutable_run_time_date(), dependencyInput->program(), u->idx(),
                                         syzkallerConditionAddress, syzkallerUncoveredAddress);

                        for (auto write_syscall : *writeAddress->mutable_write_syscall()) {
                            set_runtime_data(write_syscall.mutable_run_time_date(), dependencyInput->program(),
                                             u->idx(), syzkallerConditionAddress, writeAddress->write_address());
                        }

                    }

                    this->send_dependency_input(dependencyInput);
                }
            }
        }


    }

    void DependencyControlCenter::send_dependency_input(Input *dependencyInput) {
        if (dependencyInput != nullptr) {
            std::cout << "SendDependencyInput sig : " << dependencyInput->sig() << std::endl;
            std::cout << "dependencyInput size : " << dependencyInput->ByteSizeLong() << std::endl;
            if (dependencyInput->ByteSizeLong() < 0x7fffffff) {
                auto reply = client->SendDependencyInput(*dependencyInput);
            } else {
                std::cout << "dependencyInput is too big : " << dependencyInput->ByteSizeLong() << std::endl;
            }
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
        } else {

        }
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

    void DependencyControlCenter::get_write_addresses() {
        dra::Conditions *cs = client->GetCondition();
        for (auto condition : *cs->mutable_condition()) {
            sta::MODS *write_basicblock = get_write_basicblock(&condition);

            if (write_basicblock == nullptr) {

            } else {

                WriteAddresses *wa = new WriteAddresses();
                wa->set_allocated_condition(&condition);

                for (auto &x : *write_basicblock) {
                    WriteAddress *writeAddress = wa->add_writeaddress();
                    get_write_address(x, &condition, writeAddress);
                }

                send_write_address(wa);
            }
        }
    }

    void DependencyControlCenter::send_write_address(WriteAddresses *writeAddress) {
        if (writeAddress != nullptr) {
            std::cout << "send_write_address : " << writeAddress->condition().condition_address() << std::endl;
            auto reply = client->SendWriteAddress(*writeAddress);
#if DEBUG_RPC
#endif
        } else {

        }

    }

    sta::MODS *DependencyControlCenter::get_write_basicblock(Condition *u) {

        //TODO: need to free "allBasicblock" and "cmds" to avoid memory leak,
        // or we can also set up a cache to avoid repeated query to STA.

        sta::MODS *res = nullptr;

        DBasicBlock *p = DM.Address2BB[u->condition_address()]->parent;
        p->dump();

        llvm::BasicBlock *b = dra::getFinalBB(p->basicBlock);


        outputTime("GetAllGlobalWriteBBs : ");

        int64_t successor = u->successor();
        int64_t idx;
        if (successor == 1) {
            idx = 0;
        } else if (successor == 2) {
            idx = 1;
        } else {
            idx = 0;
            std::cerr << "switch case : " << std::hex << successor << std::endl;
        }

        if ((this->staticResult.find(b) != this->staticResult.end()) &&
            (this->staticResult[b].find(idx) != this->staticResult[b].end())) {
            res = this->staticResult[b][idx];
        } else {
            sta::MODS *write_basicblock = this->STA.GetAllGlobalWriteBBs(b, idx);
            if (write_basicblock == nullptr) {
                // no taint or out side
                std::cout << "allBasicblock == nullptr" << std::endl;
            } else if (write_basicblock->size() == 0) {
                // unrelated to gv
                std::cout << "allBasicblock->size() == 0" << std::endl;
            } else if (!write_basicblock->empty()) {
                std::cout << "get useful static analysis result : " << std::dec << write_basicblock->size()
                          << std::endl;
                res = write_basicblock;
            }

            this->staticResult[b].insert(std::pair<uint64_t, sta::MODS *>(idx, res));
        }


        return res;
    }

    void DependencyControlCenter::set_runtime_data(runTimeData *r, std::string program, uint32_t idx,
                                                   uint32_t condition, uint32_t address) {
        r->set_program(program);
        r->set_task_status(runTimeData_taskStatus_untested);
        r->set_rcursive_count(0);
        r->set_idx(idx);
        r->set_checkcondition(false);
        r->set_condition_address(condition);
        r->set_checkaddress(false);
        r->set_address(address);
        r->set_checkrightbranchaddress(false);
    }

    void DependencyControlCenter::get_write_address(sta::Mod *write_basicblock, Condition *condition,
                                                    WriteAddress *writeAddress) {

        dra::outputTime("write basicblock : ");

        dra::dump_inst(&write_basicblock->B->front());

        DBasicBlock *db = this->DM.get_DB_from_bb(write_basicblock->B);
        unsigned int write_address = DM.getSyzkallerAddress(db->trace_pc_address);


        std::cout << "write_address getSyzkallerAddress : " << std::hex << write_address << "\n";
        std::cout << "write_basicblock->repeat : " << std::hex << write_basicblock->repeat << "\n";
        std::cout << "write_basicblock->prio : " << std::hex << write_basicblock->prio << "\n";
        db->dump();

        dra::outputTime("get_cmd_ctx : start");
        std::vector<sta::cmd_ctx *> *cmd_ctx = write_basicblock->get_cmd_ctx();
        std::cout << "cmd size : " << std::dec << cmd_ctx->size() << "\n";
        dra::outputTime("get_cmd_ctx : finish");
        for (auto c: *cmd_ctx) {
            std::cout << "cmd dec: " << std::dec << c->cmd << "\n";
            std::cout << "cmd hex: " << std::hex << c->cmd << "\n";
            this->DM.dump_ctxs(&c->ctx);
        }

        writeAddress->set_write_address(write_address);
        writeAddress->set_condition_address(condition->syzkaller_condition_address());
        writeAddress->set_repeat(write_basicblock->repeat);
        writeAddress->set_prio(write_basicblock->prio);
        writeAddress->mutable_run_time_date();

        auto function_name = "ioctl";
        std::cout << "for (auto c : *cmd_ctx) {" << std::endl;
        for (auto c : *cmd_ctx) {
            std::cout << "for (auto c : *cmd_ctx) {" << std::endl;
            auto write_syscall = writeAddress->add_write_syscall();
            write_syscall->set_name(function_name);
            write_syscall->set_cmd(c->cmd);
            write_syscall->mutable_run_time_date();


            auto mm = write_syscall->mutable_critical_condition();
            bool parity = false;
            Condition *indirect_call = nullptr;
            for (auto i : c->ctx) {
                std::cout << "for (auto i : c->ctx) {" << std::endl;
                parity = !parity;
                if (parity) {
                    std::cout << "if (parity) {" << std::endl;
                    auto db = this->DM.get_DB_from_i(i);
                    if(db != nullptr){
                        std::cout << "if(db != nullptr){" << std::endl;
                        db->parent->compute_arrive();
                        if (indirect_call != nullptr) {
                            std::cout << "if (indirect_call != nullptr) {" << std::endl;
                            indirect_call->add_right_branch_address(db->trace_pc_address);
                            this->DM.set_condition(indirect_call);
                            auto ca = indirect_call->syzkaller_condition_address();
                            (*mm)[ca] = *indirect_call;
                        }
                    }
                } else {
                    std::cout << "if (parity) { else " << std::endl;
                    auto db = this->DM.get_DB_from_i(i);
                    if(db != nullptr) {
                        std::cout << "if(db != nullptr) {" << std::endl;
                        auto cc = db->critical_condition;
                        for (auto ccc : cc) {
                            std::cout << "for (auto ccc : cc) {" << std::endl;
                            this->DM.set_condition(ccc.second);
                            auto ca = ccc.second->syzkaller_condition_address();
                            (*mm)[ca] = *ccc.second;
                        }
                        indirect_call = new Condition();
                        indirect_call->set_condition_address(db->trace_pc_address);
                    }
                }
            }
        }
//        for (auto i : db->input) {
//            auto write_input = writeAddress->add_write_input();
//            write_input->set_sig(i.first->sig);
//        }
    }

} /* namespace dra */
