RPC_DIR=internal/rpc
PROTO_DIR=internal/proto
PROTO=$(wildcard $(PROTO_DIR)/*.proto)
PROTO_GEN=$(patsubst $(PROTO_DIR)/%.proto,$(RPC_DIR)/%.pb.go,$(PROTO))
PROTO_GRPC_GEN=$(patsubst $(PROTO_DIR)/%.proto,$(RPC_DIR)/%_grpc.go,$(PROTO))

all:
	echo $(PROTO_GEN)

protobufs: $(PROTO_GEN) $(PROTO_GRPC_GEN)

$(RPC_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	protoc \
		-I=${PROTO_DIR} \
		--go_out=$(RPC_DIR) \
		--go_opt=paths=source_relative \
		$< \
		--experimental_allow_proto3_optional

$(RPC_DIR)/%_grpc.go: $(PROTO_DIR)/%.proto
	protoc \
		-I=${PROTO_DIR} \
		--go-grpc_out=$(RPC_DIR) \
		--go-grpc_opt=paths=source_relative \
		$< \
		--experimental_allow_proto3_optional
