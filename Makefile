# Default directories
RPC_DIR=internal/rpc
PROTO_DIR=proto
MODULES_DIR=cmd
BIN_DIR=bin
BIN9_DIR=9bin

# Wildcards for generating modules
MODULES=$(patsubst $(MODULES_DIR)/%,$(BIN_DIR)/%,$(wildcard $(MODULES_DIR)/*))
MODULES9=$(patsubst $(MODULES_DIR)/%,$(BIN9_DIR)/%,$(wildcard $(MODULES_DIR)/*))
PROTO=$(wildcard $(PROTO_DIR)/*.proto)
PROTO_GEN=$(patsubst $(PROTO_DIR)/%.proto,$(RPC_DIR)/%.pb.go,$(PROTO))
PROTO_GRPC_GEN=$(patsubst $(PROTO_DIR)/%.proto,$(RPC_DIR)/%_grpc.pb.go,$(filter-out $(PROTO_DIR)/messages.proto,$(PROTO)))

# Docker image names
DOCKER_IMGS=minerva_go_rest minerva_go_user minerva_go_session minerva_go_products minerva_go_tenant minerva_go_migrate

# Golang build flags
export CGO_ENABLED := 0

.PHONY: all clean purge docker gen-migration migrate migrate-docker

all: protobufs modules

# Build all services
modules: $(MODULES)

# Generate protocol buffers code
protobufs: $(PROTO_GEN) $(PROTO_GRPC_GEN)

# Generate Docker images
docker: $(DOCKER_IMGS)

# Clean service builds
clean:
	rm -rf $(MODULES)

# Purge all generated protobuf implementations.
# WARNING -- This impacts on source code and commits!
purge:
	rm -f $(PROTO_GEN) $(PROTO_GRPC_GEN)


# (Experimental) Generate Plan 9 ARM binaries
plan9: export GOARCH := arm
plan9: export GOOS   := plan9
plan9: BIN_DIR       := 9bin
plan9: $(MODULES9)

# ===================

# Generation of Minerva modules
$(BIN_DIR)/%: $(MODULES_DIR)/%/main.go
	go generate $<
	go build -o $@ $<

$(BIN9_DIR)/%: $(MODULES_DIR)/%/main.go
	go generate $<
	go build -o $@ $<

# ===================

# Generation of Protocol Buffer types
$(RPC_DIR)/%.pb.go: $(PROTO_DIR)/%.proto
	protoc \
		-I=${PROTO_DIR} \
		--go_out=$(RPC_DIR) \
		--go_opt=paths=source_relative \
		$< \
		--experimental_allow_proto3_optional

# Generation of gRPC server and client boilerplate
$(RPC_DIR)/%_grpc.pb.go: $(PROTO_DIR)/%.proto
	protoc \
		-I=${PROTO_DIR} \
		--go-grpc_out=$(RPC_DIR) \
		--go-grpc_opt=paths=source_relative \
		$< \
		--experimental_allow_proto3_optional

# ============

# Generate and push Docker images for AMD64 and ARM64
minerva_go_%:
	docker buildx build \
		-f _deploy/Dockerfile \
		--platform=linux/amd64,linux/arm64 \
		--target $@ \
		-t luksamuk/$@:latest \
		--push \
		.


# ============

# Execution of modules for debug purpose4s

run-%:
	go generate cmd/$(subst run-,,$@)/main.go
	go run cmd/$(subst run-,,$@)/main.go


# ============

# Execution of Docker Compose

minerva-up:
	docker compose up

database-up:
	docker compose up -d --no-deps mariadb

database-down:
	docker compose down mariadb

# ============

# Migrations

gen-migration:
	atlas migrate diff --env gorm

migrate:
	atlas migrate apply --url "maria://mysql:mysql@localhost:3306/minerva"
migrate-k8s:
	kubectl apply -n minerva-system -f _deploy/k8s/jobs.yml
