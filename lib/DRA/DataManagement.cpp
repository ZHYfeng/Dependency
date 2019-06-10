/*
 * DataManagement.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DataManagement.h"
#include "llvm/IR/CFG.h"
#include <fstream>
#include <iostream>

#define PATH_SIZE 1000000

namespace dra {

    DataManagement::DataManagement() {
        vmOffsets = 0;
        Modules = new dra::DModule();
        Address2BB.reserve(1000000);
    }

    DataManagement::~DataManagement() = default;

    void
    DataManagement::initializeModule(std::string objdump, std::string AssemblySourceCode, std::string InputFilename) {
#if DEBUGOBJDUMP
        std::string obj = objdump.substr(0, objdump.find(".objdump"));
        std::string Cmd = "addr2line -afi -e " + obj;
        std::cout << "o Cmd :" << Cmd << std::endl;
#endif
        Modules->ReadBC(std::move(InputFilename));
        Modules->ReadObjdump(std::move(objdump));
        Modules->ReadAsmSourceCode(std::move(AssemblySourceCode));
        BuildAddress2BB();

    }

    void DataManagement::BuildAddress2BB() {
        for (auto file : Modules->Function) {
            for (auto function : file.second) {
                if (function.second->isRepeat()) {

                } else {
                    for (auto inst : function.second->InstASM) {
                        Address2BB[inst->address] = inst;
                    }
                }
            }
        }
    }

    void DataManagement::getInput(std::string coverfile) {

        std::string Line;
        std::ifstream coverFile(coverfile);
        if (coverFile.is_open()) {
            while (getline(coverFile, Line)) {

                DInput *input;
                if (Inputs.find(Line) != Inputs.end()) {
                    input = Inputs[Line];

#if DEBUGINPUT
                    std::cout << "repeat sig : " << Line << std::endl;
#endif
                    getline(coverFile, Line);
                } else {
                    input = new DInput;
                    Inputs[Line] = input;
                    input->setSig(Line);
                    getline(coverFile, Line);
                    input->setProg(Line);
                }
                input->Number++;
                getline(coverFile, Line);
                input->setCover(Line, vmOffsets);
            }
        } else {
            std::cerr << "Unable to open coverfile file " << coverfile << "\n";
        }

        for (auto i : Inputs) {
            for (auto ii : i.second->MaxCover) {
                cover[ii] = std::time(NULL);
            }
        }

        setInput();

#if 0 && DEBUGINPUT
        std::cout << "all cover: " << std::endl;
        for(auto i : cover){
            std::cout << std::hex << i << "\n";
        }
#endif
    }

    void DataManagement::getVmOffsets(std::string vmOffsets) {
        std::string Line;
        std::ifstream VmOffsets(vmOffsets);
        if (VmOffsets.is_open()) {
            while (getline(VmOffsets, Line)) {
                this->vmOffsets = std::stoul(Line, nullptr, 10);
                this->vmOffsets = (this->vmOffsets << 32);
            }
        } else {
            std::cerr << "Unable to open vmOffsets file " << vmOffsets << "\n";
        }
    }

    void DataManagement::setInput() {
        for (auto it : this->Inputs) {
            std::string sig = it.first;
            for (auto addr : it.second->MaxCover) {
                if (this->Address2BB.find(addr) != this->Address2BB.end()) {
                    this->Address2BB[addr]->update(CoverKind::cover, it.second);
                } else {
                    std::cerr << "un find address " << std::hex << addr << "\n";
                }

            }
        }
    }

    void DataManagement::setVmOffsets(unsigned long long int vmOffsets) {
        this->vmOffsets = (vmOffsets << 32);
        std::cout << "GetVmOffsets : " << std::hex << this->vmOffsets << std::endl;
    }

    DInput *DataManagement::getInput(Input input) {
        std::string sig = input.sig();
#if DEBUGINPUT
        std::cout << "sig : " << sig << std::endl;
#endif
        DInput *dInput;
        if (Inputs.find(sig) != Inputs.end()) {
            dInput = Inputs[sig];
        } else {
            dInput = new DInput;
            Inputs[sig] = dInput;
            dInput->setSig(sig);
        }
        dInput->Number++;
        for (auto c : input.call()) {
            dInput->idx = c.second.idx();
            for (auto a : c.second.address()) {
                unsigned long long int address = a.first;
                auto final_address = getRealAddress(address);
                if (this->Address2BB.find(final_address) != this->Address2BB.end()) {
                    this->Address2BB[final_address]->update(CoverKind::cover, dInput);
                    this->dump_address(final_address);
                } else {
                    std::cerr << "un find address " << std::hex << final_address << "\n";
                }

                if (this->cover.find(final_address) == this->cover.end()) {
                    auto current_time = std::time(NULL);
                    coverage *c = new coverage();
                    c->time = current_time;
                    c->address = final_address;
                    this->cover[final_address] = current_time;
                    this->time.push_back(c);
                    std::cout << std::ctime(&current_time) << "new cover address " << std::hex << final_address << "\n";
                    if (this->uncover.find(final_address) == this->uncover.end()) {

                    } else {
                        this->uncover[final_address]->covered = true;
                        if (input.dependency() == true) {
                            this->uncover[final_address]->covered_by_dependency = true;
                        }
                    }

                } else {
                }
            }
        }

        std::vector<DUncoveredAddress *> temp;
        for (auto ua : dInput->dUncoveredAddress) {
            if (this->cover.find(ua->address) == this->cover.end()) {
                temp.push_back(ua);
                if (this->uncover.find(ua->address) == this->uncover.end()) {
                    auto current_time = std::time(NULL);
                    auto ui = new uncover_info();
                    ui->time = current_time;
                    ui->address = ua->address;
                    this->uncover[ua->address] = ui;
                }
            } else {
                delete ua;
            }
        }
        dInput->dUncoveredAddress.clear();
        for (auto ua : temp) {
            dInput->dUncoveredAddress.push_back(ua);
        }
        return dInput;
    }

    unsigned long long int DataManagement::getRealAddress(unsigned long long int address) {
        return address + this->vmOffsets - 5;
    }

    unsigned long long int DataManagement::getSyzkallerAddress(unsigned long long int address) {
        return address - this->vmOffsets + 5;
    }

    bool DataManagement::isDriver(unsigned long long int address) {

        if (this->Address2BB.find(address) != this->Address2BB.end()) {
            if (this->Address2BB[address]->parent != nullptr) {
                auto b = this->Address2BB[address]->parent;
                if (b->parent != nullptr) {
                    auto f = b->parent;
                    std::cout << "isDriver path : " << f->Path << "\n";
                    std::cout << "isDriver address : " << address << "\n";
                    if (f->Path.find("block/") == 0) {
                        return true;
                    } else if (f->Path.find("drivers/") == 0) {
                        return true;
                    } else if (f->Path.find("sound/") == 0) {
                        return true;
                    }
                } else {
                    std::cerr << "isDriver not have parent f : " << std::hex << address << "\n";
                }
            } else {
                std::cerr << "isDriver not have parent bb : " << std::hex << address << "\n";
            }
        } else {
            std::cerr << "isDriver not find address : " << std::hex << address << "\n";
        }
        return false;
    }

    void DataManagement::dump_address(unsigned long long int address) {

        if (this->Address2BB.find(address) != this->Address2BB.end()) {
            if (this->Address2BB[address]->parent != nullptr) {
                auto b = this->Address2BB[address]->parent;
                if (b->parent != nullptr) {
                    auto f = b->parent;
                    std::cout << "dump_address path : " << f->Path << "\n";
                    std::cout << "dump_address address : " << this->getSyzkallerAddress(address) << "\n";
                } else {
                    std::cerr << "dump_address not have parent f : " << std::hex << address << "\n";
                }
            } else {
                std::cerr << "dump_address not have parent bb : " << std::hex << address << "\n";
            }
        } else {
            std::cerr << "dump_address not find address : " << std::hex << address << "\n";
        }
    }

    llvm::BasicBlock *DataManagement::getRealBB(llvm::BasicBlock *b) {
        if (b->hasName()) {
            return b;
        } else {
            for (auto *Pred : llvm::predecessors(b)) {
                return getRealBB(Pred);
            }
        }
    }

    llvm::BasicBlock *DataManagement::getFinalBB(llvm::BasicBlock *b) {
        auto *inst = b->getTerminator();
        for (unsigned int i = 0, end = inst->getNumSuccessors(); i < end; i++) {
            std::string name = inst->getSuccessor(i)->getName().str();
            if (inst->getSuccessor(i)->hasName()) {
            } else {
                return getFinalBB(inst->getSuccessor(i));
            }
        }
        return b;
    }

    void DataManagement::dump_cover() {
        std::ofstream out_file("cover_uncover.txt",
                               std::ios_base::out | std::ios_base::app);
        auto current_time = std::time(NULL);
        out_file << std::ctime(&current_time);
        out_file << "this->cover.size() : " << this->cover.size() << "\n";
        out_file.close();
    }

    void DataManagement::dump_uncover() {
        std::ofstream out_file("cover_uncover.txt",
                               std::ios_base::out | std::ios_base::app);

        int ud = 0, ug = 0, ucc = 0, ucd = 0;
        for (auto uc : this->uncover) {
            if (uc.second->belong_to_Driver) {
                ud++;
                if (uc.second->related_to_gv) {
                    ug++;
                    if (uc.second->covered) {
                        ucc++;
                        if (uc.second->covered_by_dependency) {
                            ucd++;
                        }
                    }
                }
            }
        }

        auto current_time = std::time(NULL);
        out_file << std::ctime(&current_time);
        out_file << "this->uncover.size() : " << this->uncover.size() << "\n";
        out_file << "belong to driver : " << ud << "\n";
        out_file << "related to gv : " << ug << "\n";
        out_file << "be covered : " << ucc << "\n";
        out_file << "be covered by dependency : " << ucd << "\n";

        out_file.close();
    }

    uncover_info::uncover_info() : address(0),
                                   belong_to_Driver(false),
                                   related_to_gv(false),
                                   covered(false),
                                   covered_by_dependency(false) {}
} /* namespace dra */

