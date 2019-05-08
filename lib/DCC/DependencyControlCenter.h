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
#include "../STA/StaticAnalysisResult.h"

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
        sta::StaticAnalysisResult STA;
    };

} /* namespace dra */

#endif /* LIB_DCC_DEPENDENCYCONTROLCENTER_H_ */
