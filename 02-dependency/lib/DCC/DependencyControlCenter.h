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

        void init(std::string obj_dump, std::string AssemblySourceCode, std::string InputFilename,
                  const std::string &staticRes = "", const std::string &function = "",
                  const std::string &port_address = "");

        void run();

        void check_uncovered_addresses_depednency(const std::string &file);

        void setRPCConnection(const std::string &grpc_port);

        sta::MODS *get_write_basicblock(Condition *u);

        void ckeck_input_dependency(DInput *dInput);

        void send_dependency(Dependency *dependency);

        void get_write_address(sta::Mod *write_basicblock, Condition *condition, WriteAddress *writeAddress);

        writeAddressAttributes *get_write_addresses_adttributes(sta::Mod *write_basicblock);

        static void set_runtime_data(runTimeData *r, const std::string &program, uint32_t idx, uint32_t condition,
                                     uint32_t address);

        void check_confition_depednency();

        void send_write_address(WriteAddresses *writeAddress);

        void test_sta();

        static void test_rpc();

        void getFileOperations(std::string *function_name, std::string *file_operations, std::string *kind);

    private:
        DependencyRPCClient *client{};
        std::string port;
        DataManagement DM;
        sta::StaticAnalysisResult STA;
        nlohmann::json function_json;

        std::time_t start_time{};

        std::map<llvm::BasicBlock *, std::map<uint64_t, sta::MODS *>> staticResult;

    };

} /* namespace dra */

#endif /* LIB_DCC_DEPENDENCYCONTROLCENTER_H_ */
