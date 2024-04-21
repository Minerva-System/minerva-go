# minerva-go

Golang refactor of the [Minerva System](https://minerva-system.github.io/minerva-system/).

## Building

Use the `Makefile` to build. You will need a Golang compiler with minimum
version 1.20.4.

```bash
make            # To build everything
make protobufs  # To regenerate Protocol Buffers implementations
make clean      # To remove generated services
make purge      # To remove generated files such as Swagger/OpenAPI
                #    and protobufs
make docker     # To build Docker images (see next sections)
make run-<cmd>  # To run a specific service (see the cmd directory)
```

If you wish to rebuild the protocol buffer files, you'll also need `protoc`
and a few more dependencies. Furthermore, to regenerate Swagger and OpenAPI
documentation for the REST service, you'll need Swag.

To use `make run-<cmd>` commands, you need to replace `<cmd>` with the desired
service name, which must always be one of the directory names within the `cmd`
directory, e.g. to run the REST server, use `make run-rest`.

<!-- ```bash -->
<!-- go install github.com/mitranim/gow@latest -->
<!-- ``` -->

### Generating protocol buffer implementations

First of all, install `protoc`, the protobuf compiler. Seek out the best
package for your Linux distribution or, on Windows, use a package manager
such as Chocolatey.

After that, use the Go package manager to install dependencies:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Finally, use `make protobufs` for an automated generation of files.

### Generating Swagger and OpenAPI documentation

First install the `swag` tool:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Finally, `cd` to `cmd/rest`, then run `swag init`.

### Building Docker images

If you wish to build a single Docker image, use a command like the
following:

```
docker image build \
	-f _deploy/Dockerfile \
	--target MODULENAME \
	-t MODULENAME \
	.
```

Possible names for `MODULENAME`:

- `minerva_go_rest` (REST API Gateway + Swagger, runs on port 9000);
- `minerva_go_user` (gRPC, runs on port 9010);
- `minerva_go_session` (gRPC, runs on port 9011);
- `minerva_go_products` (gRPC, runs on port 9012).
- `minerva_go_tenant` (gRPC, runs on port 9013).

#### Cross-platform build

If you're using BuildKit, you can also generate cross-platform images
(and push them to DockerHub) by abusing the `--platform` argument.
The following example generates an image of `minerva_go_rest` for
Linux on 64-bit Intel and ARM platforms, then pushes them to
DockerHub under my username:

```
docker buildx build \
	-f _deploy/Dockerfile \
	--target minerva_go_rest \
	--platform=linux/amd64,linux/arm64 \
	-t luksamuk/minerva_go_rest \
	--push \
	.
```

### Building for Plan 9

There is an experimental build method which will build the modules for Plan 9
under an ARM architecture. This is due to my current setup where Plan 9 is
run in a Raspberry Pi 3 Model B+ with a 32-bit ARM kernel.

If you wish to build Plan 9 binaries, use the following:

```bash
make plan9
```

Binaries will be created in `9bin`. You may also want to take a look at the
script in `extra/runsvc.rc` for ease of use in Plan 9.

## Migrations

Migrations are executed by using the [Atlas](https://atlasgo.io/) tool.
For more information on Atlas integration with GORM, see [this link](https://atlasgo.io/guides/orms/gorm).

To install Atlas, use:

```bash
curl -sSf https://atlasgo.sh | sh
```

This also uses the Atlas GORM Provider, which is installed as a project dependency.

