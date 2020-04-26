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
#include "../RPC/DependencyRPC.pb.h"

#define DEBUG_INPUT 0

namespace dra {

    class DInput {
    public:
        DInput();

        virtual ~DInput();

        void setSig(const std::string &sig);

        void setProgram(const std::string &program);

        void setCover(const std::string &cover, unsigned long long int vmOffsets);

        Condition* getCondition(uint64_t condition, uint64_t uncovered, const std::vector<uint64_t>& branch, int i) const;

        void addConditionAddress(uint64_t c);

        void addUncoveredAddress(Condition *c);

    public:
        std::string sig;
        std::string program;
        unsigned long long int Number;
        std::vector<std::set<unsigned long long int> *> AllCover;
        std::set<unsigned long long int> MaxCover;
        std::set<unsigned long long int> MiniCover;

        unsigned long long int idx;
        std::set<uint32_t> dConditionAddress; // all dConditionAddress comes from llvm bc
        std::vector<Condition *> dUncoveredAddress; // all dUncoveredAddress comes from llvm bc
    };

} /* namespace dra */

#endif /* LIB_DRA_DINPUT_H_ */
