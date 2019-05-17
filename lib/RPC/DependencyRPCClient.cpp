/*
 * DependencyRPCClient.cpp
 *
 *  Created on: May 1, 2019
 *      Author: yhao
 */

#include "DependencyRPCClient.h"

#include <grpcpp/grpcpp.h>

#include "DependencyRPC.pb.h"

namespace dra {

    DependencyRPCClient::DependencyRPCClient(const std::shared_ptr<grpc::Channel> &channel) :
            stub_(DependencyRPC::NewStub(channel)) {
    }

    DependencyRPCClient::~DependencyRPCClient() = default;


    unsigned long long int DependencyRPCClient::GetVmOffsets() {
        Empty request;
        Empty reply;
        grpc::ClientContext context;
        grpc::Status status = stub_->GetVmOffsets(&context, request, &reply);
        if (status.ok()) {
            std::cout << "GetVmOffsets : " << std::hex << reply.address() << std::endl;
            return reply.address();
        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
            return 0;
        }
    }

    NewInput *DependencyRPCClient::GetNewInput() {
        Empty request;
        auto *reply = new NewInput();
        grpc::ClientContext context;
        grpc::Status status = stub_->GetNewInput(&context, request, reply);
        if (status.ok()) {
            return reply;
        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
            return nullptr;
        }
    }

    void DependencyRPCClient::SendDependencyInput(const DependencyInput &request) {
        Empty reply;
        grpc::ClientContext context;
        grpc::Status status = stub_->SendDependencyInput(&context, request, &reply);
        if (status.ok()) {

        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
        }
    }

} /* namespace dra */
