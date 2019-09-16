/*
 * DInput.cpp
 *
 *  Created on: Mar 22, 2019
 *      Author: yhao
 */

#include "DInput.h"

#include <iostream>
#include <sstream>

namespace dra {

    DInput::DInput() {
        Number = 0;
    }

    DInput::~DInput() = default;

    void DInput::setSig(const std::string &sig) {
        this->sig = sig;
#if DEBUG_INPUT
        std::cout << this->sig << std::endl;
#endif
    }

    void DInput::setProgram(const std::string &program) {
        this->program = program;
//        std::stringstream ss;
//        std::string temp = program.substr(1, program.size() - 2) + ' ';
//        char c = 0;
//
//#if DEBUG_INPUT
//        std::cout << program << std::endl;
//#endif
//        for (auto cc : temp) {
//            if (cc != ' ') {
//                c = static_cast<char>(c * 10 + cc - '0');
//            } else {
//                ss.str("");
//                ss << c;
//                c = 0;
//                this->program += ss.str();
//            }
//        }
//#if DEBUG_INPUT
//        std::cout << this->program << std::endl;
//#endif
    }

    void DInput::setCover(const std::string &cover, unsigned long long int vmOffsets) {
        std::string temp = cover.substr(1, cover.size() - 2) + ' ';
        unsigned long long int addr = 0;
        auto *thisCover = new std::set<unsigned long long int>;
        auto *tempCover = new std::set<unsigned long long int>;
        this->AllCover.push_back(thisCover);
#if DEBUG_INPUT
        std::cout << cover << std::endl;
#endif
        for (auto cc : temp) {
            if (cc != ' ') {
                addr = addr * 10 + cc - '0';
            } else {
                auto FinalAddr = addr + vmOffsets - 5;
#if DEBUG_INPUT
                if (this->MaxCover.find(FinalAddr) == this->MaxCover.end()) {
                    std::cout << "new : " << std::hex << FinalAddr << std::endl;
                } else {
                    std::cout << "old : " << std::hex << FinalAddr << std::endl;
                }
#endif
                this->MaxCover.insert(FinalAddr);
                thisCover->insert(FinalAddr);
                addr = 0;
            }
        }
        if (this->MiniCover.empty()) {
            this->MiniCover = *thisCover;
        } else {
            for (auto cc : this->MiniCover) {
                if (thisCover->find(cc) != thisCover->end()) {
                    tempCover->insert(cc);
                } else {

                }
            }
            MiniCover = *tempCover;
        }

#if DEBUG_INPUT
        std::cout << "MiniCover:\n";
        for (auto i : this->MiniCover) {
            std::cout << std::hex << i << " ";
        }
        std::cout << "\n";
        std::cout << "MaxCover:\n";
        for (auto i : this->MaxCover) {
            std::cout << std::hex << i << " ";
        }
        std::cout << "\n";
#endif
    }

    void DInput::addUncoveredAddress(unsigned long long int uncoveredAddress, unsigned long long int conditionAddress, int i) {
        Condition *d = new Condition();
        d->set_condition_address(conditionAddress);
        d->set_uncovered_address(uncoveredAddress);
        d->set_idx(this->idx);
        d->set_successor(1 << i);
        dUncoveredAddress.push_back(d);
//
//        std::cout << "uncovered trace_pc_address : " << std::hex << trace_pc_address << std::endl;
//        std::cout << "conditionAddress : " << std::hex << conditionAddress << std::endl;
    }

} /* namespace dra */