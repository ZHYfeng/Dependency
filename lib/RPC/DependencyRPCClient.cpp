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
            return reply.address();
        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
            return 0;
        }
    }

    Inputs *DependencyRPCClient::GetNewInput() {
        Empty request;
        auto *reply = new Inputs();
        grpc::ClientContext context;
        grpc::Status status = stub_->GetNewInput(&context, request, reply);
        if (status.ok()) {
            return reply;
        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
            return nullptr;
        }
    }

    Empty *DependencyRPCClient::SendDependencyInput(const Input &request) {
        Empty *reply = new Empty();
        grpc::ClientContext context;
        grpc::Status status = stub_->SendDependencyInput(&context, request, reply);
        if (status.ok()) {
            std::cout << "SendDependencyInput : " << reply->name() << std::endl;
        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
        }
        return reply;
    }

    Inputs *DependencyRPCClient::GetDependencyInput() {
        Empty request;
        Inputs *reply = new Inputs;
        grpc::ClientContext context;
        request.set_name("vm-0");
        grpc::Status status = stub_->GetDependencyInput(&context, request, reply);
        if (status.ok()) {

        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
        }
        return reply;
    }

    Conditions *DependencyRPCClient::GetCondition() {
        Empty request;
        auto *reply = new Conditions();
        grpc::ClientContext context;
        grpc::Status status = stub_->GetCondition(&context, request, reply);
        if (status.ok()) {
            return reply;
        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
            return nullptr;
        }
    }

    Empty *DependencyRPCClient::SendWriteAddress(const WriteAddresses &request) {
        Empty *reply = new Empty();
        grpc::ClientContext context;
        grpc::Status status = stub_->SendWriteAddress(&context, request, reply);
        if (status.ok()) {
            std::cout << "SendDependencyInput : " << reply->name() << std::endl;
        } else {
            std::cout << status.error_code() << ": " << status.error_message() << std::endl;
        }
        return reply;
    }

} /* namespace dra */
