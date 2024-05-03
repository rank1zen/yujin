GO := go
PROTOC := protoc

gen:
	$(PROTOC) "--go_out=." "--go-grpc_out=." "proto/profile.proto"
