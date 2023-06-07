# minerva-go

## Generating protocol buffer implementations

Install `protoc`.

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Run `swag init` on `cmd/rest`.


