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
#if DEBUGINPUT
        std::cout << this->sig << std::endl;
#endif
    }

    void DInput::setProg(const std::string &prog) {

        std::stringstream ss;
        std::string temp = prog.substr(1, prog.size() - 2) + ' ';
        char c = 0;

#if DEBUGINPUT
        std::cout << prog << std::endl;
#endif
        for (auto cc : temp) {
            if (cc != ' ') {
                c = static_cast<char>(c * 10 + cc - '0');
            } else {
                ss.str("");
                ss << c;
                c = 0;
                this->progam += ss.str();
            }
        }
#if DEBUGINPUT
        std::cout << this->progam << std::endl;
#endif
    }

    void DInput::setCover(const std::string &cover, unsigned long long int vmOffsets) {
        std::string temp = cover.substr(1, cover.size() - 2) + ' ';
        unsigned long long int addr = 0;
        auto *thisCover = new std::set<unsigned long long int>;
        auto *tempCover = new std::set<unsigned long long int>;
        this->AllCover.push_back(thisCover);
#if DEBUGINPUT
        std::cout << cover << std::endl;
#endif
        for (auto cc : temp) {
            if (cc != ' ') {
                addr = addr * 10 + cc - '0';
            } else {
                auto FinalAddr = addr + vmOffsets - 5;
#if DEBUGINPUT
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

#if DEBUGINPUT
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

    void DInput::addUncoveredAddress(unsigned long long int address, unsigned long long int condition_address, int i) {
        DUncoveredAddress *d = new DUncoveredAddress();
        d->address = address;
        d->successor_idx = i;
        d->idx = this->idx;
        d->condition_address = condition_address;
        dUncoveredAddress.push_back(d);
//
//        std::cout << "uncovered address : " << std::hex << address << std::endl;
//        std::cout << "condition_address : " << std::hex << condition_address << std::endl;
    }

} /* namespace dra */
