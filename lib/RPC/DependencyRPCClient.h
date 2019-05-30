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

namespace dra {

    class DependencyRPCClient {
    public:
        DependencyRPCClient(const std::shared_ptr<grpc::Channel> &channel);

        virtual ~DependencyRPCClient();

        unsigned long long int GetVmOffsets();

        NewInput *GetNewInput();

        Empty *SendDependencyInput(const DependencyInput &request);

        NewDependencyInput *GetDependencyInput();

    private:
        std::unique_ptr<DependencyRPC::Stub> stub_;
    };

} /* namespace dra */

#endif /* LIB_RPC_DEPENDENCYRPCCLIENT_H_ */
