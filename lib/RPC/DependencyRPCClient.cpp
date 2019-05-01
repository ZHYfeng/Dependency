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

DependencyRPCClient::DependencyRPCClient(std::shared_ptr<grpc::Channel> channel) :
		stub_(DependencyRPC::NewStub(channel)) {
	// TODO Auto-generated constructor stub
}

DependencyRPCClient::~DependencyRPCClient() {
	// TODO Auto-generated destructor stub
}

std::string DependencyRPCClient::SendDependencyInput(const DependencyInput dependencyInput) {
	// Data we are sending to the server.
	DependencyInput request;
	// Container for the data we expect from the server.
	Empty reply;
	// Context for the client. It could be used to convey extra information to
	// the server and/or tweak certain RPC behaviors.
	grpc::ClientContext context;
	// The actual RPC.
	grpc::Status status = stub_->SendDependencyInput(&context, request, &reply);
	// Act upon its status.
	if (status.ok()) {
		return reply.SerializeAsString();
	} else {
		std::cout << status.error_code() << ": " << status.error_message() << std::endl;
		return "RPC failed";
	}
}

} /* namespace dra */
