compile:
	g++ -std=c++17 -I/opt/homebrew/include \
			$(shell pkg-config --cflags grpc++ protobuf) \
			-o server server.cc service.grpc.pb.cc service.pb.cc \
			$(shell pkg-config --libs grpc++ protobuf) -lpthread -ldl

proto-go:
	rm -f client/*.pb.*
	protoc --go_out=client --go_opt=paths=source_relative \
    --go-grpc_out=client --go-grpc_opt=paths=source_relative \
		service.proto

proto-cpp:
	rm -f *.pb.*
	protoc --grpc_out=. --cpp_out=. service.proto \
		--plugin=protoc-gen-grpc=/opt/homebrew/bin/grpc_cpp_plugin
