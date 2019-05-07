/*
 * DependencyControlCenter.h
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#ifndef LIB_DCC_DEPENDENCYCONTROLCENTER_H_
#define LIB_DCC_DEPENDENCYCONTROLCENTER_H_

#include "../DRA/DataManagement.h"
#include "../RPC/DependencyRPCClient.h"
#include "../JSON/json.cpp"

typedef std::vector<std::string> LOC_INF;

namespace dra {

    class DependencyControlCenter {
    public:
        DependencyControlCenter();

        virtual ~DependencyControlCenter();

        void init(std::string objdump, std::string AssemblySourceCode, std::string InputFilename, const std::string &staticRes);

        void run();

    private:
        DependencyRPCClient client;
        DataManagement DM;

        nlohmann::json j_taintedBrs, j_analysisCtxMap, j_tagMap, j_modInstCtxMap;

    public:
        int initStaticRes(const std::string &staticRes);

        LOC_INF *getLocInf(llvm::Instruction *);

        LOC_INF *getLocInf(llvm::BasicBlock *);
    };

} /* namespace dra */

#endif /* LIB_DCC_DEPENDENCYCONTROLCENTER_H_ */