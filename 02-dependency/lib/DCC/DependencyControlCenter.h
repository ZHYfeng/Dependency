/*
 * DependencyControlCenter.h
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#ifndef LIB_DCC_DEPENDENCYCONTROLCENTER_H_
#define LIB_DCC_DEPENDENCYCONTROLCENTER_H_

#include <string>

#include "../DMM/DataManagement.h"
#include "../RPC/DependencyRPCClient.h"
#include "../STA/StaticAnalysisResult.h"

namespace dra {

    class DependencyControlCenter {
    public:
        DependencyControlCenter();

        virtual ~DependencyControlCenter();

        void init(const std::string &obj_dump, const std::string &assembly, const std::string &bit_code,
                  const std::string &config, const std::string &port_address = "");

        void run();

        void check_uncovering_addresses_dependnency(const std::string &file);

        void check_write_addresses_dependency(const std::string &file);

        void check_all_condition_();

        void setRPCConnection(const std::string &grpc_port);

        void send_number_basicblock_covered();

        void check_input(DInput *dInput);

        sta::MODS *get_write_basicblock(Condition *u);

        sta::MODS *get_write_basicblock(u_int64_t address, u_int32_t idx = 0);

        sta::MODS *get_write_basicblock(dra::DBasicBlock *db, u_int32_t idx = 0);

        void write_basic_block_to_address(sta::Mod *write_basicblock, Condition *condition, WriteAddress *writeAddress);

        void write_basic_block_to_adttributes(sta::Mod *write_basicblock, writeAddressAttributes *waa);

        static void set_runtime_data(runTimeData *r, const std::string &program, uint32_t idx, uint32_t condition,
                                     uint32_t address);

        void send_dependency(Dependency *dependency);

        void check_condition();

        void send_write_address(WriteAddresses *writeAddress);

        void test_sta();

        void test_rpc();

        void test();

        void getFileOperations(std::string *function_name, std::string *file_operations, std::string *kind);

        sta::StaticAnalysisResult* getStaticAnalysisResult(const std::string& path);

        bool is_dependency(dra::DBasicBlock *db);

    private:
        DependencyRPCClient *client{};
        std::string port;
        DataManagement DM;
        std::map<std::string, sta::StaticAnalysisResult*> STA_map;
        nlohmann::json config_json;

        std::time_t start_time{};

        std::map<llvm::BasicBlock *, std::map<uint64_t, sta::MODS *>> staticResult;

    };

} /* namespace dra */

#endif /* LIB_DCC_DEPENDENCYCONTROLCENTER_H_ */
