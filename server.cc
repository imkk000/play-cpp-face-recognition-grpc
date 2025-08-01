#include "service.grpc.pb.h"
#include <cstdlib>
#include <grpc++/grpc++.h>
#include <iostream>

class MessageServiceImpl : public MessageService::Service {
public:
  grpc::Status Ping(grpc::ServerContext *context, const PingRequest *request,
                    PingResponse *response) override {

    std::cout << "Received: " << request->message() << std::endl;
    response->set_reply("Hello from C++! Got: " + request->message());
    return grpc::Status::OK;
  }
};

int main() {
  const char *addr_env = std::getenv("ADDR");
  std::string addr = addr_env ? addr_env : "0.0.0.0:50051";
  std::string server_address(addr);
  MessageServiceImpl service;

  grpc::ServerBuilder builder;
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);

  std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
  std::cout << "Server listening on " << server_address << std::endl;
  server->Wait();
  return 0;
}
