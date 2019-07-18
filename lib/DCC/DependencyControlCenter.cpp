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
        this->outputTime("start_time");
    }

    DependencyControlCenter::~DependencyControlCenter() = default;

    void DependencyControlCenter::init(std::string objdump, std::string AssemblySourceCode, std::string InputFilename,
                                       const std::string &staticRes) {

        DM.initializeModule(std::move(objdump), std::move(AssemblySourceCode), std::move(InputFilename));
        this->outputTime("initializeModule");



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

                this->outputTime("sleep_for 10s");
                std::this_thread::sleep_for(std::chrono::seconds(10));
            } else {
                this->outputTime("sleep_for 60s");
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
        this->outputTime("GetVmOffsets");
    }

    void DependencyControlCenter::get_dependency_input(DInput *dInput) {

        Input *dependencyInput = new Input();
        dependencyInput->set_sig(dInput->sig);
        dependencyInput->set_program(dInput->program);

        bool send_flag = false;
        std::cout << "dUncoveredAddress size : " << std::dec << dInput->dUncoveredAddress.size()
                  << std::endl;
        uint64_t i = 0;
        for (auto u : dInput->dUncoveredAddress) {
            i++;
            this->outputTime("uncovered address count : " + std::to_string(i));
            if (this->DM.check_uncovered_address(u)) {

                if (this->DM.uncover.find(u->uncovered_address()) != this->DM.uncover.end()) {
                    this->DM.uncover[u->uncovered_address()]->belong_to_Driver = true;
                }

                sta::MODS *write_basicblock = get_write_basicblock(u);

                if (write_basicblock == nullptr) {

                } else {

                    unsigned long long int syzkallerUncoveredAddress = DM.getSyzkallerAddress(u->uncovered_address());
                    unsigned long long int syzkallerConditionAddress = DM.getSyzkallerAddress(u->condition_address());

                    this->current_time = std::time(nullptr);
                    std::cout << std::ctime(&current_time);
                    std::cout << "condition trace_pc_address : " << std::hex << u->condition_address() << "\n";
                    std::cout << "uncovered trace_pc_address : " << std::hex << u->uncovered_address() << "\n";
                    std::cout << "condition getSyzkallerAddress : " << std::hex << syzkallerConditionAddress << "\n";
                    std::cout << "uncovered getSyzkallerAddress : " << std::hex << syzkallerUncoveredAddress << "\n";

                    UncoveredAddress *uncoveredAddress = dependencyInput->add_uncovered_address();
                    uncoveredAddress->set_uncovered_address(syzkallerUncoveredAddress);
                    uncoveredAddress->set_condition_address(syzkallerConditionAddress);

                    set_runtime_data(uncoveredAddress->mutable_run_time_date(), dependencyInput->program(), u->idx(),
                                     syzkallerConditionAddress, syzkallerUncoveredAddress);

                    send_flag = true;

                    if (this->DM.uncover.find(u->uncovered_address()) != this->DM.uncover.end()) {
                        this->DM.uncover[u->uncovered_address()]->related_to_gv = true;
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
                        writeAddress->set_repeat(x->repeat);
                        writeAddress->set_prio(x->prio);
                        writeAddress->set_write_address(write_address);
                        writeAddress->set_condition_address(syzkallerConditionAddress);
                        set_runtime_data(writeAddress->mutable_run_time_date(), dependencyInput->program(), u->idx(),
                                         syzkallerConditionAddress, syzkallerUncoveredAddress);

                        auto function_name = "ioctl";
                        for (auto c : *cmd_ctx) {
                            auto write_syscall = writeAddress->add_write_syscall();
                            write_syscall->set_name(function_name);
                            write_syscall->set_cmd(c->cmd);

                            bool parity = false;
                            auto mm = write_syscall->mutable_critical_condition();
                            for (auto i : c->ctx) {
                                parity = !parity;
                                if (parity) {

                                    auto db = this->DM.get_DB_from_bb(i->getParent());
                                    db->parent->compute_arrive();
                                } else {
                                    auto cc = this->DM.get_DB_from_bb(i->getParent())->critical_condition;
                                    for (auto ccc : cc) {
                                        auto ca = ccc.second->syzkaller_condition_address();
                                        (*mm)[ca] = *ccc.second;
                                    }
                                }
                            }

                            set_runtime_data(write_syscall->mutable_run_time_date(), dependencyInput->program(),
                                             u->idx(),
                                             syzkallerConditionAddress, syzkallerUncoveredAddress);
                        }

                        // need something
                        for (auto i : db->input) {
                            auto write_input = writeAddress->add_write_input();
                            write_input->set_sig(i.first->sig);
                            write_input->set_program(i.first->program);
                        }
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

    void DependencyControlCenter::get_write_addresses() {
        dra::Conditions *cs = client->GetCondition();
        for (auto condition : *cs->mutable_condition()) {
            sta::MODS *write_basicblock = get_write_basicblock(&condition);
            if (write_basicblock == nullptr) {

            } else {
                WriteAddresses *wa = new WriteAddresses();
                wa->set_allocated_condition(&condition);

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

                    WriteAddress *writeAddress = wa->add_writeaddress();
                    writeAddress->set_write_address(write_address);
                    writeAddress->set_condition_address(condition.syzkaller_condition_address());
                    writeAddress->set_repeat(x->repeat);
                    writeAddress->set_prio(x->prio);

                    auto function_name = "ioctl";
                    for (auto c : *cmd_ctx) {
                        auto write_syscall = writeAddress->add_write_syscall();
                        write_syscall->set_name(function_name);
                        write_syscall->set_cmd(c->cmd);


                        bool parity = false;
                        auto mm = write_syscall->mutable_critical_condition();
                        for (auto i : c->ctx) {
                            parity = !parity;
                            if (parity) {
                                this->DM.get_DB_from_bb(i->getParent())->parent->compute_arrive();
                            } else {
                                auto cc = this->DM.get_DB_from_bb(i->getParent())->critical_condition;
                                for (auto ccc : cc) {
                                    auto ca = ccc.second->syzkaller_condition_address();
                                    (*mm)[ca] = *ccc.second;
                                }
                            }
                        }
                    }
                    for (auto i : db->input) {
                        auto write_input = writeAddress->add_write_input();
                        write_input->set_sig(i.first->sig);
                    }
                }

                this->send_write_address(wa);
            }
        }
    }

    void DependencyControlCenter::send_write_address(WriteAddresses *writeAddress) {
        std::cout << "send_write_address : " << writeAddress->condition().condition_address() << std::endl;
        auto reply = client->SendWriteAddress(*writeAddress);
#if DEBUG_RPC
#endif
    }

    sta::MODS *DependencyControlCenter::get_write_basicblock(Condition *u) {

        //TODO: need to free "allBasicblock" and "cmds" to avoid memory leak,
        // or we can also set up a cache to avoid repeated query to STA.

        sta::MODS *res = nullptr;

        DBasicBlock *p = DM.Address2BB[u->condition_address()]->parent;
        p->dump();

        llvm::BasicBlock *b = dra::getFinalBB(p->basicBlock);

        this->current_time = std::time(nullptr);
        std::cout << std::ctime(&current_time);
        std::cout << "GetAllGlobalWriteBBs : " << std::endl;

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
            res = write_basicblock;
            this->current_time = std::time(nullptr);
            std::cout << std::ctime(&current_time);

            if (write_basicblock == nullptr) {
                // no taint or out side
                std::cout << "allBasicblock == nullptr" << std::endl;
            } else if (write_basicblock->size() == 0) {
                // unrelated to gv
                std::cout << "allBasicblock->size() == 0" << std::endl;
            } else if (!write_basicblock->empty()) {
                std::cout << "get useful static analysis result : " << std::dec << write_basicblock->size()
                          << std::endl;
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

    void DependencyControlCenter::outputTime(std::string s) {
        this->current_time = std::time(nullptr);
        std::cout << std::ctime(&current_time);
        std::cout << "#time : " << s << std::endl;
    }

} /* namespace dra */
