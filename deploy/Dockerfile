# NOTE: Actual path should be project root!
FROM golang:1.20.4-alpine3.16 AS builder
RUN mkdir /minerva
WORKDIR /minerva
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
RUN apk add --no-cache make protobuf protobuf-dev && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
COPY . .
RUN make

FROM alpine:3.16 as certs
RUN apk --no-cache add ca-certificates

FROM scratch as runner
ENV TZ=Etc/UTC
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

FROM runner AS minerva_go_rest
EXPOSE 9000
COPY --from=builder /minerva/bin/rest ./minerva-rest
CMD ["./minerva-rest"]

FROM runner AS minerva_go_user
EXPOSE 9010
COPY --from=builder /minerva/bin/user ./minerva-user
CMD ["./minerva-user"]

FROM runner AS minerva_go_session
EXPOSE 9011
COPY --from=builder /minerva/bin/session ./minerva-session
CMD ["./minerva-session"]

FROM runner AS minerva_go_products
EXPOSE 9012
COPY --from=builder /minerva/bin/products ./minerva-products
CMD ["./minerva-products"]
