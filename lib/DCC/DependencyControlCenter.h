/*
 * DependencyControlCenter.h
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#ifndef LIB_DCC_DEPENDENCYCONTROLCENTER_H_
#define LIB_DCC_DEPENDENCYCONTROLCENTER_H_

#include <string>

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

        std::time_t start_time;
        std::time_t finish_init_time;
        std::time_t report_time;
    };

} /* namespace dra */

#endif /* LIB_DCC_DEPENDENCYCONTROLCENTER_H_ */
