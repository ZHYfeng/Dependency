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

        void setRPCConnection();

        sta::MODS *get_write_basicblock(Condition *u);

        void get_dependency_input(DInput *dInput);

        void send_dependency_input(Input *dependencyInput);

        void get_write_addresses();

        void send_write_address(WriteAddresses *writeAddress);

        void set_runtime_data(runTimeData *r, std::string program, uint32_t idx, uint32_t condition, uint32_t address);

        void outputTime(std::string);

        void test_sta();
        void test_rpc();

    private:
        DependencyRPCClient *client;
        DataManagement DM;
        sta::StaticAnalysisResult STA;

        std::time_t start_time;
        std::time_t current_time;

    };

} /* namespace dra */

#endif /* LIB_DCC_DEPENDENCYCONTROLCENTER_H_ */
