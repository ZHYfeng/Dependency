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

    }

    DInput::~DInput() {

    }

    void DInput::setCover(const std::string &cover, unsigned long long int vmOffsets) {
        std::string temp = cover.substr(1, cover.size() - 2) + ' ';
        unsigned long long int addr = 0;

#if DEBUGINPUT
        std::cout << cover << std::endl;
#endif
        for (auto cc : temp) {
            if (cc != ' ') {
                addr = addr * 10 + cc - '0';
            } else {
                this->cover.insert(addr + vmOffsets - 5);
                addr = 0;
            }
        }
#if DEBUGINPUT
        for (auto addr : this->cover) {
            std::cout << std::hex << addr << " ";
        }
        std::cout << "\n";
#endif
    }

    void DInput::setProg(const std::string &prog) {

        std::stringstream ss;
        std::string temp = prog.substr(1, prog.size() - 2) + ' ';
        char c;

#if DEBUGINPUT
        std::cout << prog << std::endl;
#endif
        for (auto cc : temp) {
            if (cc != ' ') {
                c = c * 10 + cc - '0';
            } else {
                ss.str("");
                ss << c;
                c = 0;
                this->prog += ss.str();
            }
        }
#if DEBUGINPUT
        std::cout << this->prog << std::endl;
#endif
    }

    void DInput::setSig(const std::string &sig) {
        this->sig = sig;
#if DEBUGINPUT
        std::cout << this->sig << std::endl;
#endif
    }

} /* namespace dra */
