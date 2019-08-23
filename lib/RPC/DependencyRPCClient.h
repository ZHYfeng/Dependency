/*
 * DependencyRPCClient.h
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#ifndef LIB_RPC_DEPENDENCYRPCCLIENT_H_
#define LIB_RPC_DEPENDENCYRPCCLIENT_H_

#include <memory>

#include "DependencyRPC.grpc.pb.h"

#define DEBUG_RPC 1

namespace dra {

    class DependencyRPCClient {
    public:
        explicit DependencyRPCClient(const std::shared_ptr<grpc::Channel> &channel);

        virtual ~DependencyRPCClient();

        uint32_t GetVmOffsets();

        void SendBasicBlockNumber(uint32_t BasicBlockNumber);

        Inputs *GetNewInput();

        Empty *SendDependency(const Dependency &request);

        Conditions *GetCondition();

        Empty *SendWriteAddress(const WriteAddresses &request);

    private:
        std::unique_ptr<DependencyRPC::Stub> stub_;
    };

} /* namespace dra */

#endif /* LIB_RPC_DEPENDENCYRPCCLIENT_H_ */
