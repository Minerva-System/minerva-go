# minerva-go

## Generating protocol buffer implementations

Install `protoc`.

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Run `swag init` on `cmd/rest`.


## Building images

```
docker image build \
	-f deploy/Dockerfile \
	--target MODULENAME \
	-t MODULENAME \
	.
```

Possible names for `MODULENAME`:

- `minerva_go_rest` (REST API Gateway + Swagger, runs on port 9000);
- `minerva_go_user` (gRPC, runs on port 9010);
- `minerva_go_session` (gRPC, runs on port 9011);
- `minerva_go_products` (gRPC, runs on port 9012).

### Cross-platform build

If you're using BuildKit, you can also generate cross-platform images
(and push them to DockerHub) by abusing the `--platform` argument.
The following example generates an image of `minerva_go_rest` for
Linux on 64-bit Intel and ARM platforms, then pushes them to
DockerHub under my username:

```
docker buildx build \
	-f deploy/Dockerfile \
	--target minerva_go_rest \
	--platform=linux/amd64,linux/arm64 \
	-t luksamuk/minerva_go_rest \
	--push \
	.
```

