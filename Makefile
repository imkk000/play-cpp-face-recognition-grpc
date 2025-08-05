compile:
	g++ -Wall -Werror -std=c++17 \
			$(shell pkg-config --cflags --libs grpc++ protobuf opencv4) \
			-o server server.cc service.grpc.pb.cc service.pb.cc

proto-go:
	rm -f client/*.pb.*
	protoc --go_out=client --go_opt=paths=source_relative \
    --go-grpc_out=client --go-grpc_opt=paths=source_relative \
		service.proto

proto-cpp:
	rm -f *.pb.*
	protoc --grpc_out=. --cpp_out=. service.proto \
		# --plugin=protoc-gen-grpc=/opt/homebrew/bin/grpc_cpp_plugin
		--plugin=protoc-gen-grpc=/usr/bin/grpc_cpp_plugin
