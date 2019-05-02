// Generated by the gRPC C++ plugin.
// If you make any local change, they will be lost.
// source: DependencyRPC.proto

#include "DependencyRPC.pb.h"
#include "DependencyRPC.grpc.pb.h"

#include <functional>
#include <grpcpp/impl/codegen/async_stream.h>
#include <grpcpp/impl/codegen/async_unary_call.h>
#include <grpcpp/impl/codegen/channel_interface.h>
#include <grpcpp/impl/codegen/client_unary_call.h>
#include <grpcpp/impl/codegen/client_callback.h>
#include <grpcpp/impl/codegen/method_handler_impl.h>
#include <grpcpp/impl/codegen/rpc_service_method.h>
#include <grpcpp/impl/codegen/server_callback.h>
#include <grpcpp/impl/codegen/service_type.h>
#include <grpcpp/impl/codegen/sync_stream.h>
namespace dra {

static const char* DependencyRPC_method_names[] = {
        "/dra.DependencyRPC/GetDependencyInput",
        "/dra.DependencyRPC/GetNewInput",
        "/dra.DependencyRPC/GetVmOffsets",
        "/dra.DependencyRPC/SendDependencyInput",
        "/dra.DependencyRPC/SendInput",
};

std::unique_ptr< DependencyRPC::Stub> DependencyRPC::NewStub(const std::shared_ptr< ::grpc::ChannelInterface>& channel, const ::grpc::StubOptions& options) {
  (void)options;
  std::unique_ptr< DependencyRPC::Stub> stub(new DependencyRPC::Stub(channel));
  return stub;
}

DependencyRPC::Stub::Stub(const std::shared_ptr< ::grpc::ChannelInterface>& channel)
  : channel_(channel), rpcmethod_GetDependencyInput_(DependencyRPC_method_names[0], ::grpc::internal::RpcMethod::NORMAL_RPC, channel),
    rpcmethod_GetNewInput_(DependencyRPC_method_names[1], ::grpc::internal::RpcMethod::NORMAL_RPC, channel),
    rpcmethod_GetVmOffsets_(DependencyRPC_method_names[2], ::grpc::internal::RpcMethod::NORMAL_RPC, channel),
    rpcmethod_SendDependencyInput_(DependencyRPC_method_names[3], ::grpc::internal::RpcMethod::NORMAL_RPC, channel),
    rpcmethod_SendInput_(DependencyRPC_method_names[4], ::grpc::internal::RpcMethod::NORMAL_RPC, channel)
  {}

::grpc::Status DependencyRPC::Stub::GetDependencyInput(::grpc::ClientContext* context, const ::dra::Empty& request, ::dra::DependencyInput* response) {
  return ::grpc::internal::BlockingUnaryCall(channel_.get(), rpcmethod_GetDependencyInput_, context, request, response);
}

void DependencyRPC::Stub::experimental_async::GetDependencyInput(::grpc::ClientContext* context, const ::dra::Empty* request, ::dra::DependencyInput* response, std::function<void(::grpc::Status)> f) {
  return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_GetDependencyInput_, context, request, response, std::move(f));
}

void DependencyRPC::Stub::experimental_async::GetDependencyInput(::grpc::ClientContext* context, const ::grpc::ByteBuffer* request, ::dra::DependencyInput* response, std::function<void(::grpc::Status)> f) {
  return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_GetDependencyInput_, context, request, response, std::move(f));
}

::grpc::ClientAsyncResponseReader< ::dra::DependencyInput>* DependencyRPC::Stub::AsyncGetDependencyInputRaw(::grpc::ClientContext* context, const ::dra::Empty& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::dra::DependencyInput>::Create(channel_.get(), cq, rpcmethod_GetDependencyInput_, context, request, true);
}

::grpc::ClientAsyncResponseReader< ::dra::DependencyInput>* DependencyRPC::Stub::PrepareAsyncGetDependencyInputRaw(::grpc::ClientContext* context, const ::dra::Empty& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::dra::DependencyInput>::Create(channel_.get(), cq, rpcmethod_GetDependencyInput_, context, request, false);
}

    ::grpc::Status DependencyRPC::Stub::GetNewInput(::grpc::ClientContext *context, const ::dra::Empty &request, ::dra::NewInput *response) {
        return ::grpc::internal::BlockingUnaryCall(channel_.get(), rpcmethod_GetNewInput_, context, request, response);
    }

    void DependencyRPC::Stub::experimental_async::GetNewInput(::grpc::ClientContext *context, const ::dra::Empty *request, ::dra::NewInput *response,
                                                              std::function<void(::grpc::Status)> f) {
        return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_GetNewInput_, context, request, response, std::move(f));
    }

    void DependencyRPC::Stub::experimental_async::GetNewInput(::grpc::ClientContext *context, const ::grpc::ByteBuffer *request, ::dra::NewInput *response,
                                                              std::function<void(::grpc::Status)> f) {
        return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_GetNewInput_, context, request, response, std::move(f));
    }

    ::grpc::ClientAsyncResponseReader<::dra::NewInput> *
    DependencyRPC::Stub::AsyncGetNewInputRaw(::grpc::ClientContext *context, const ::dra::Empty &request, ::grpc::CompletionQueue *cq) {
        return ::grpc::internal::ClientAsyncResponseReaderFactory<::dra::NewInput>::Create(channel_.get(), cq, rpcmethod_GetNewInput_, context, request, true);
    }

    ::grpc::ClientAsyncResponseReader<::dra::NewInput> *
    DependencyRPC::Stub::PrepareAsyncGetNewInputRaw(::grpc::ClientContext *context, const ::dra::Empty &request, ::grpc::CompletionQueue *cq) {
        return ::grpc::internal::ClientAsyncResponseReaderFactory<::dra::NewInput>::Create(channel_.get(), cq, rpcmethod_GetNewInput_, context, request, false);
    }

    ::grpc::Status DependencyRPC::Stub::GetVmOffsets(::grpc::ClientContext *context, const ::dra::Empty &request, ::dra::Empty *response) {
        return ::grpc::internal::BlockingUnaryCall(channel_.get(), rpcmethod_GetVmOffsets_, context, request, response);
    }

    void DependencyRPC::Stub::experimental_async::GetVmOffsets(::grpc::ClientContext *context, const ::dra::Empty *request, ::dra::Empty *response,
                                                               std::function<void(::grpc::Status)> f) {
        return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_GetVmOffsets_, context, request, response, std::move(f));
    }

    void DependencyRPC::Stub::experimental_async::GetVmOffsets(::grpc::ClientContext *context, const ::grpc::ByteBuffer *request, ::dra::Empty *response,
                                                               std::function<void(::grpc::Status)> f) {
        return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_GetVmOffsets_, context, request, response, std::move(f));
    }

    ::grpc::ClientAsyncResponseReader<::dra::Empty> *
    DependencyRPC::Stub::AsyncGetVmOffsetsRaw(::grpc::ClientContext *context, const ::dra::Empty &request, ::grpc::CompletionQueue *cq) {
        return ::grpc::internal::ClientAsyncResponseReaderFactory<::dra::Empty>::Create(channel_.get(), cq, rpcmethod_GetVmOffsets_, context, request, true);
    }

    ::grpc::ClientAsyncResponseReader<::dra::Empty> *
    DependencyRPC::Stub::PrepareAsyncGetVmOffsetsRaw(::grpc::ClientContext *context, const ::dra::Empty &request, ::grpc::CompletionQueue *cq) {
        return ::grpc::internal::ClientAsyncResponseReaderFactory<::dra::Empty>::Create(channel_.get(), cq, rpcmethod_GetVmOffsets_, context, request, false);
}

::grpc::Status DependencyRPC::Stub::SendDependencyInput(::grpc::ClientContext* context, const ::dra::DependencyInput& request, ::dra::Empty* response) {
  return ::grpc::internal::BlockingUnaryCall(channel_.get(), rpcmethod_SendDependencyInput_, context, request, response);
}

void DependencyRPC::Stub::experimental_async::SendDependencyInput(::grpc::ClientContext* context, const ::dra::DependencyInput* request, ::dra::Empty* response, std::function<void(::grpc::Status)> f) {
  return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_SendDependencyInput_, context, request, response, std::move(f));
}

void DependencyRPC::Stub::experimental_async::SendDependencyInput(::grpc::ClientContext* context, const ::grpc::ByteBuffer* request, ::dra::Empty* response, std::function<void(::grpc::Status)> f) {
  return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_SendDependencyInput_, context, request, response, std::move(f));
}

::grpc::ClientAsyncResponseReader< ::dra::Empty>* DependencyRPC::Stub::AsyncSendDependencyInputRaw(::grpc::ClientContext* context, const ::dra::DependencyInput& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::dra::Empty>::Create(channel_.get(), cq, rpcmethod_SendDependencyInput_, context, request, true);
}

::grpc::ClientAsyncResponseReader< ::dra::Empty>* DependencyRPC::Stub::PrepareAsyncSendDependencyInputRaw(::grpc::ClientContext* context, const ::dra::DependencyInput& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::dra::Empty>::Create(channel_.get(), cq, rpcmethod_SendDependencyInput_, context, request, false);
}

    ::grpc::Status DependencyRPC::Stub::SendInput(::grpc::ClientContext *context, const ::dra::Input &request, ::dra::Empty *response) {
        return ::grpc::internal::BlockingUnaryCall(channel_.get(), rpcmethod_SendInput_, context, request, response);
    }

    void DependencyRPC::Stub::experimental_async::SendInput(::grpc::ClientContext *context, const ::dra::Input *request, ::dra::Empty *response,
                                                            std::function<void(::grpc::Status)> f) {
        return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_SendInput_, context, request, response, std::move(f));
    }

    void DependencyRPC::Stub::experimental_async::SendInput(::grpc::ClientContext *context, const ::grpc::ByteBuffer *request, ::dra::Empty *response,
                                                            std::function<void(::grpc::Status)> f) {
        return ::grpc::internal::CallbackUnaryCall(stub_->channel_.get(), stub_->rpcmethod_SendInput_, context, request, response, std::move(f));
    }

    ::grpc::ClientAsyncResponseReader<::dra::Empty> *
    DependencyRPC::Stub::AsyncSendInputRaw(::grpc::ClientContext *context, const ::dra::Input &request, ::grpc::CompletionQueue *cq) {
        return ::grpc::internal::ClientAsyncResponseReaderFactory<::dra::Empty>::Create(channel_.get(), cq, rpcmethod_SendInput_, context, request, true);
    }

    ::grpc::ClientAsyncResponseReader<::dra::Empty> *
    DependencyRPC::Stub::PrepareAsyncSendInputRaw(::grpc::ClientContext *context, const ::dra::Input &request, ::grpc::CompletionQueue *cq) {
        return ::grpc::internal::ClientAsyncResponseReaderFactory<::dra::Empty>::Create(channel_.get(), cq, rpcmethod_SendInput_, context, request, false);
}

DependencyRPC::Service::Service() {
  AddMethod(new ::grpc::internal::RpcServiceMethod(
      DependencyRPC_method_names[0],
      ::grpc::internal::RpcMethod::NORMAL_RPC,
      new ::grpc::internal::RpcMethodHandler< DependencyRPC::Service, ::dra::Empty, ::dra::DependencyInput>(
          std::mem_fn(&DependencyRPC::Service::GetDependencyInput), this)));
  AddMethod(new ::grpc::internal::RpcServiceMethod(
          DependencyRPC_method_names[1],
          ::grpc::internal::RpcMethod::NORMAL_RPC,
          new ::grpc::internal::RpcMethodHandler<DependencyRPC::Service, ::dra::Empty, ::dra::NewInput>(
                  std::mem_fn(&DependencyRPC::Service::GetNewInput), this)));
  AddMethod(new ::grpc::internal::RpcServiceMethod(
          DependencyRPC_method_names[2],
          ::grpc::internal::RpcMethod::NORMAL_RPC,
          new ::grpc::internal::RpcMethodHandler<DependencyRPC::Service, ::dra::Empty, ::dra::Empty>(
                  std::mem_fn(&DependencyRPC::Service::GetVmOffsets), this)));
    AddMethod(new ::grpc::internal::RpcServiceMethod(
            DependencyRPC_method_names[3],
            ::grpc::internal::RpcMethod::NORMAL_RPC,
            new ::grpc::internal::RpcMethodHandler< DependencyRPC::Service, ::dra::DependencyInput, ::dra::Empty>(
                    std::mem_fn(&DependencyRPC::Service::SendDependencyInput), this)));
  AddMethod(new ::grpc::internal::RpcServiceMethod(
          DependencyRPC_method_names[4],
          ::grpc::internal::RpcMethod::NORMAL_RPC,
          new ::grpc::internal::RpcMethodHandler< DependencyRPC::Service, ::dra::Input, ::dra::Empty>(
                  std::mem_fn(&DependencyRPC::Service::SendInput), this)));
}

DependencyRPC::Service::~Service() {
}

::grpc::Status DependencyRPC::Service::GetDependencyInput(::grpc::ServerContext* context, const ::dra::Empty* request, ::dra::DependencyInput* response) {
  (void) context;
  (void) request;
  (void) response;
  return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
}

    ::grpc::Status DependencyRPC::Service::GetNewInput(::grpc::ServerContext *context, const ::dra::Empty *request, ::dra::NewInput *response) {
        (void) context;
        (void) request;
        (void) response;
        return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }

    ::grpc::Status DependencyRPC::Service::GetVmOffsets(::grpc::ServerContext *context, const ::dra::Empty *request, ::dra::Empty *response) {
  (void) context;
  (void) request;
  (void) response;
  return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
}

::grpc::Status DependencyRPC::Service::SendDependencyInput(::grpc::ServerContext* context, const ::dra::DependencyInput* request, ::dra::Empty* response) {
  (void) context;
  (void) request;
  (void) response;
  return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
}

    ::grpc::Status DependencyRPC::Service::SendInput(::grpc::ServerContext *context, const ::dra::Input *request, ::dra::Empty *response) {
  (void) context;
  (void) request;
  (void) response;
  return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
}


}  // namespace dra

