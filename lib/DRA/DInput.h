/*
 * DInput.h
 *
 *  Created on: Mar 22, 2019
 *      Author: yhao
 */

#ifndef LIB_DRA_DINPUT_H_
#define LIB_DRA_DINPUT_H_

#include <set>
#include <vector>
#include <string>

#define DEBUGINPUT 0

namespace dra {

    class DUncoveredAddress {
    public:
        unsigned long long int address;
        unsigned long long int idx;
        unsigned long long int condition_address;
    };

    class DInput {
    public:
        DInput();

        virtual ~DInput();

        void setSig(const std::string &sig);

        void setProg(const std::string &prog);

        void setCover(const std::string &cover, unsigned long long int vmOffsets);

        void addUncoveredAddress(unsigned long long int address, unsigned long long int condition_address);

    public:
        std::string sig;
        std::string progam;
        unsigned long long int Number;
        std::vector<std::set<unsigned long long int> *> AllCover;
        std::set<unsigned long long int> MaxCover;
        std::set<unsigned long long int> MiniCover;

        unsigned long long int idx;
        std::vector<DUncoveredAddress *> dUncoveredAddress;
    };

} /* namespace dra */

#endif /* LIB_DRA_DINPUT_H_ */
