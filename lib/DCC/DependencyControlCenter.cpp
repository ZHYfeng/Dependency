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
#include "../DRA/DFunction.h"
#include "general.h"

namespace dra {

    DependencyControlCenter::DependencyControlCenter() {
        dra::outputTime("start_time");
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
//                    std::cout << "new input : " << input.second.sig() << std::endl;
//                    std::cout << input.second.program() << std::endl;
//                    DInput *dInput = DM.getInput(&input.second);
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

//            this->DM.dump_cover();
//            this->DM.dump_uncover();

            get_write_addresses();
        }
    }


    void DependencyControlCenter::setRPCConnection() {
        this->client = new dra::DependencyRPCClient(
                grpc::CreateChannel("127.0.0.1:22223", grpc::InsecureChannelCredentials()));
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

                    Dependency *dependency = new Dependency();
                    Input *input = dependency->mutable_input();
                    input->set_sig(dInput->sig);
                    input->set_program(dInput->program);

                    unsigned long long int syzkallerConditionAddress = DM.getSyzkallerAddress(u->condition_address());
                    unsigned long long int syzkallerUncoveredAddress = DM.getSyzkallerAddress(u->uncovered_address());
#if DEBUG
                    outputTime("");
                    std::cout << "condition trace_pc_address : " << std::hex << u->condition_address() << "\n";
                    std::cout << "uncovered trace_pc_address : " << std::hex << u->uncovered_address() << "\n";
                    std::cout << "condition getSyzkallerAddress : " << std::hex << syzkallerConditionAddress << "\n";
                    std::cout << "uncovered getSyzkallerAddress : " << std::hex << syzkallerUncoveredAddress << "\n";
#endif

                    UncoveredAddress *uncoveredAddress = dependency->mutable_uncovered_address();
                    uncoveredAddress->set_condition_address(syzkallerConditionAddress);
                    uncoveredAddress->set_uncovered_address(syzkallerUncoveredAddress);

                    (*input->mutable_uncovered_address())[syzkallerUncoveredAddress] = u->idx();
                    (*uncoveredAddress->mutable_input())[dInput->sig] = u->idx();

                    set_runtime_data(uncoveredAddress->mutable_run_time_date(), input->program(), u->idx(),
                                     syzkallerConditionAddress, syzkallerUncoveredAddress);

//                    (*dependency->mutable_uncovered_address())[syzkallerUncoveredAddress] = *uncoveredAddress;

                    if (this->DM.uncover.find(u->uncovered_address()) != this->DM.uncover.end()) {
                        this->DM.uncover[u->uncovered_address()]->related_to_gv = true;
                    }

                    for (auto &x : *write_basicblock) {
//                        WriteAddress *writeAddress = new WriteAddress;
//                        (*uncoveredAddress->mutable_write_address())[syzkallerUncoveredAddress] = *writeAddress;
                        WriteAddress *writeAddress = dependency->add_write_address();

                        get_write_address(x, u, writeAddress);
                        writeAddressAttributes *waa = get_write_addresses_adttributes(x);
                        (*uncoveredAddress->mutable_write_address())[waa->write_address()] = *waa;
                        (*writeAddress->mutable_uncovered_address())[syzkallerUncoveredAddress] = *waa;

                        set_runtime_data(writeAddress->mutable_run_time_date(), input->program(), u->idx(),
                                         syzkallerConditionAddress, syzkallerUncoveredAddress);

//                        std::cout << "writeAddress->mutable_write_syscall()->size() : "
//                                  << writeAddress->mutable_write_syscall()->size() << std::endl;
//
//                        for(auto &wc : *writeAddress->mutable_write_syscall()){
////                            set_runtime_data(wc.second.mutable_run_time_date(), dependency->program(),
////                                             u->idx(), syzkallerConditionAddress, writeAddress->write_address());
////                            std::cout << "write_syscall.run_time_date().program() : "
////                                      << wc.second.run_time_date().program() << std::endl;
//                            set_runtime_data(wc.mutable_run_time_date(), input->program(),
//                                             u->idx(), syzkallerConditionAddress, writeAddress->write_address());
//                            std::cout << "write_syscall.run_time_date().program() : "
//                                      << wc.run_time_date().program() << std::endl;
//                        }

                    }

                    this->send_dependency(dependency);
                    dependency->Clear();
                }
            }
        }


    }

    void DependencyControlCenter::send_dependency(Dependency *dependency) {
        if (dependency != nullptr) {
#if DEBUG_RPC
            auto ua = dependency->uncovered_address();
            std::cout << "uncover condition address : " << ua.condition_address() << std::endl;
            for (auto wa : dependency->write_address()) {
                std::cout << "wa program : " << std::endl;
                std::cout << wa.run_time_date().program();
            }
            std::cout << "dependency size : " << dependency->ByteSizeLong() << std::endl;
#endif
            if (dependency->ByteSizeLong() < 0x7fffffff) {
                auto reply = client->SendDependency(*dependency);
            } else {
                std::cout << "dependency is too big : " << dependency->ByteSizeLong() << std::endl;
            }

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
//            sta::MODS *write_basicblock = get_write_basicblock(&condition.second);
            sta::MODS *write_basicblock = get_write_basicblock(&condition);
            if (write_basicblock == nullptr) {

            } else {

                WriteAddresses *wa = new WriteAddresses();
//                wa->set_allocated_condition(&condition.second);
//                for (auto &x : *write_basicblock) {
//                    WriteAddress *writeAddress = new WriteAddress;
//                    get_write_address(x, &condition.second, writeAddress);
//                    (*wa->mutable_write_address())[condition.second.syzkaller_uncovered_address()] = *writeAddress;
//                }

                wa->set_allocated_condition(&condition);
                for (auto &x : *write_basicblock) {
                    WriteAddress *writeAddress = wa->add_write_address();
                    get_write_address(x, &condition, writeAddress);
                }

                send_write_address(wa);
            }
        }
        cs->Clear();
    }

    void DependencyControlCenter::send_write_address(WriteAddresses *writeAddress) {
        if (writeAddress != nullptr) {
#if DEBUG_RPC
            std::cout << "send_write_address : " << std::hex << writeAddress->condition().condition_address()
                      << std::endl;
#endif
            auto reply = client->SendWriteAddress(*writeAddress);
        } else {

        }

    }

    sta::MODS *DependencyControlCenter::get_write_basicblock(Condition *u) {

        sta::MODS *res = nullptr;
        llvm::BasicBlock *b;
        if (this->DM.Address2BB.find(u->condition_address()) != this->DM.Address2BB.end()) {
            DBasicBlock *p = DM.Address2BB[u->condition_address()]->parent;
            p->dump();
            b = dra::getFinalBB(p->basicBlock);
        } else {
            return res;
        }

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
        r->set_task_status(taskStatus::untested);
        r->set_rcursive_count(0);
        r->set_idx(idx);
        r->set_checkcondition(false);
        r->set_condition_address(condition);
        r->set_checkaddress(false);
        r->set_address(address);
        r->set_checkrightbranchaddress(false);
        r->mutable_right_branch_address();
    }

    void DependencyControlCenter::get_write_address(sta::Mod *write_basicblock, Condition *condition,
                                                    WriteAddress *writeAddress) {


        DBasicBlock *db = this->DM.get_DB_from_bb(write_basicblock->B);
        unsigned int write_address = DM.getSyzkallerAddress(db->trace_pc_address);
#if DEBUG
        dra::outputTime("write basicblock : ");
        db->dump();
#endif

        std::vector<sta::cmd_ctx *> *cmd_ctx = write_basicblock->get_cmd_ctx();
#if DEBUG
        for (auto c: *cmd_ctx) {
            std::cout << "cmd hex: " << std::hex << c->cmd << "\n";
            this->DM.dump_ctxs(&c->ctx);
        }
#endif

        writeAddress->set_write_address(write_address);
        writeAddress->set_condition_address(condition->syzkaller_condition_address());
        writeAddress->mutable_run_time_date();

//        auto function_name = "ioctl";
//        std::cout << "for (auto c : *cmd_ctx) {" << std::endl;
//        for (auto c : *cmd_ctx) {
//            std::cout << "for (auto c : *cmd_ctx) {" << std::endl;
//            Syscall *write_syscall = writeAddress->add_write_syscall();
//
////            (*writeAddress->mutable_write_syscall())[write_address] = *write_syscall;
//
//            write_syscall->set_name(function_name);
//            write_syscall->set_cmd(c->cmd);
//            write_syscall->mutable_run_time_date();
//
//
//            auto mm = write_syscall->mutable_critical_condition();
//            bool parity = false;
//            Condition *indirect_call = nullptr;
//            for (auto i : c->ctx) {
//                parity = !parity;
//                if (parity) {
//                    auto db = this->DM.get_DB_from_i(i);
//                    if (db != nullptr) {
//                        db->parent->compute_arrive();
//                        if (indirect_call != nullptr) {
//                            indirect_call->add_right_branch_address(db->trace_pc_address);
////                            (*indirect_call->mutable_right_branch_address())[db->trace_pc_address] = 0;
//
//                            this->DM.set_condition(indirect_call);
//                            auto ca = indirect_call->syzkaller_condition_address();
//                            (*mm)[ca] = *indirect_call;
//                        }
//                    }
//                } else {
//                    auto db = this->DM.get_DB_from_i(i);
//                    if (db != nullptr) {
//                        auto cc = db->critical_condition;
//                        for (auto ccc : cc) {
//                            this->DM.set_condition(ccc.second);
//                            auto ca = ccc.second->syzkaller_condition_address();
//                            (*mm)[ca] = *ccc.second;
//                        }
//                        indirect_call = new Condition();
//                        indirect_call->set_condition_address(db->trace_pc_address);
//                    }
//                }
//            }
//        }

        for (auto i : db->input) {
            (*writeAddress->mutable_input())[i.first->sig] = i.second;
        }
    }

    writeAddressAttributes *DependencyControlCenter::get_write_addresses_adttributes(sta::Mod *write_basicblock) {
        writeAddressAttributes *res = new writeAddressAttributes();
        DBasicBlock *db = this->DM.get_DB_from_bb(write_basicblock->B);
        unsigned int write_address = DM.getSyzkallerAddress(db->trace_pc_address);
        res->set_write_address(write_address);
        res->set_repeat(write_basicblock->repeat);
        res->set_prio(write_basicblock->prio);
        return res;
    }

} /* namespace dra */
