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
                    isDriver(final_address);
                } else {
                    std::cerr << "un find address " << std::hex << final_address << "\n";
                }

                if (this->cover.find(final_address) == this->cover.end()) {
                    std::time_t t = std::time(NULL);
                    coverage *c = new coverage();
                    c->time = t;
                    c->address = final_address;
                    this->cover[final_address] = t;
                    this->time.push_back(c);
                    std::cerr << "new cover address " << std::hex << final_address << "\n";
                } else {
                }
            }
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
        if(this->Address2BB[address]->parent != nullptr ){
            auto b = this->Address2BB[address]->parent;
            if (b->parent!= nullptr){
                auto f = b->parent;
                std::cout << "isDriver path : " << f->Path << "\n";
                std::cout << "isDriver address : " << address << "\n";
                if (f->Path.find("block/") == 0) {
                    return true;
                } else if (f->Path.find("drivers/") == 0) {
                    return true;
                }
            }
        }
        return false;
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

} /* namespace dra */

