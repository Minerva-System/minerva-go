# Default directories
RPC_DIR=internal/rpc
PROTO_DIR=internal/proto
MODULES_DIR=cmd
BIN_DIR=bin

# Wildcards for generating modules
MODULES=$(patsubst $(MODULES_DIR)/%,$(BIN_DIR)/%,$(wildcard $(MODULES_DIR)/*))
PROTO=$(wildcard $(PROTO_DIR)/*.proto)
PROTO_GEN=$(patsubst $(PROTO_DIR)/%.proto,$(RPC_DIR)/%.pb.go,$(PROTO))
PROTO_GRPC_GEN=$(patsubst $(PROTO_DIR)/%.proto,$(RPC_DIR)/%_grpc.go,$(PROTO))

.PHONY: all protobufs clean

all: $(MODULES)

clean:
	rm -rf $(MODULES)

# Generation of Minerva modules
$(BIN_DIR)/%: $(MODULES_DIR)/%/main.go
	go build -o $@ $<


# Use this target while programming only!
protobufs: $(PROTO_GEN) $(PROTO_GRPC_GEN)

# Generation of Protocol Buffer types
$(RPC_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	protoc \
		-I=${PROTO_DIR} \
		--go_out=$(RPC_DIR) \
		--go_opt=paths=source_relative \
		$< \
		--experimental_allow_proto3_optional

# Generation of gRPC server and client boilerplate
$(RPC_DIR)/%_grpc.go: $(PROTO_DIR)/%.proto
	protoc \
		-I=${PROTO_DIR} \
		--go-grpc_out=$(RPC_DIR) \
		--go-grpc_opt=paths=source_relative \
		$< \
		--experimental_allow_proto3_optional
