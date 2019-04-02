/*
 * DataManagement.cpp
 *
 *  Created on: Feb 28, 2019
 *      Author: yhao
 */

#include "DataManagement.h"

#include <fstream>
#include <iostream>

#include "DModule.h"

#define PATH_SIZE 1000000

namespace dra {

    DataManagement::DataManagement() {
        vmOffsets = 0;
        FindNum = 0;
        UnFindNum = 0;
        SameNum = 0;
        DiffNum = 0;
        Modules = new dra::DModule();
        Address2BB.reserve(1000000);
    }

    DataManagement::~DataManagement() = default;

    void
    DataManagement::initializeModule(std::string objdump, std::string AssemblySourceCode, std::string InputFilename) {

        Modules->ReadBC(std::move(InputFilename));
        Modules->ReadObjdump(std::move(objdump));
        Modules->ReadAsmSourceCode(std::move(AssemblySourceCode));

    }

    void DataManagement::BuildAddress2BB(
            std::unordered_map<std::string, std::unordered_map<std::string, DFunction *>> Function) {
        for (auto file : Modules->Function) {
            for (auto function : file.second) {
                if (function.second->isRepeat()) {

                } else {
                    for (auto inst : function.second->InstASM) {
                        Address2BB[inst->Address] = inst;
                    }
                }
            }
        }
    }

} /* namespace dra */

void dra::DataManagement::getInput(std::string coverfile) {

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
            cover.insert(ii);
        }
    }
#if 0 && DEBUGINPUT
    std::cout << "all cover: " << std::endl;
    for(auto i : cover){
        std::cout << std::hex << i << "\n";
    }
#endif
}

void dra::DataManagement::getVmOffsets(std::string vmOffsets) {
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
