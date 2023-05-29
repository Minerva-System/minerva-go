RPC_DIR=pkg/rpc
PROTO_DIR=pkg/proto
PROTO=$(wildcard $(PROTO_DIR)/*.proto)
PROTO_GEN=$(patsubst $(PROTO_DIR)/%.proto,$(RPC_DIR)/%.pb.go,$(PROTO))

all:
	echo $(PROTO_GEN)

protobufs: $(PROTO_GEN)

$(RPC_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	protoc \
		-I=${PROTO_DIR} \
		--go_out=$(RPC_DIR) \
		$< \
		--experimental_allow_proto3_optional
