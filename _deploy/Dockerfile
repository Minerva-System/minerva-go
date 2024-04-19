# NOTE: Actual path should be project root!
FROM golang:1.22-alpine AS builder
RUN mkdir -p /minerva
WORKDIR /minerva
ENV CGO_ENABLED=0
RUN apk add --no-cache make protobuf protobuf-dev ca-certificates && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download
COPY . .

FROM scratch as runner
ENV TZ=Etc/UTC
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Migration executor
FROM alpine:latest AS minerva_go_migrate
RUN apk add --no-cache curl &&\
    curl -sSf https://atlasgo.sh >/atlasgo.sh &&\
    sh /atlasgo.sh -y -o /usr/bin/atlas &&\
    rm /atlasgo.sh
COPY migrations /minerva/migrations
COPY atlas.hcl /minerva/atlas.hcl
WORKDIR /minerva
CMD ["sh", "-c", "exec atlas migrate apply --url \"maria://$DATABASE_SERVICE_USER:$DATABASE_SERVICE_PASSWORD@$DATABASE_SERVICE_SERVER/$DATABASE_SERVICE_DBNAME\"" ]


# REST service
FROM builder as builder_rest
RUN make bin/rest
FROM runner AS minerva_go_rest
EXPOSE 9000
COPY --from=builder_rest /minerva/bin/rest ./minerva-rest
CMD ["./minerva-rest"]


# USER service
FROM builder as builder_user
RUN make bin/user
FROM runner AS minerva_go_user
EXPOSE 9010
COPY --from=builder_user /minerva/bin/user ./minerva-user
CMD ["./minerva-user"]


# SESSION service
FROM builder as builder_session
RUN make bin/session
FROM runner AS minerva_go_session
EXPOSE 9011
COPY --from=builder_session /minerva/bin/session ./minerva-session
CMD ["./minerva-session"]


# PRODUCTS service
FROM builder as builder_products
RUN make bin/products
FROM runner AS minerva_go_products
EXPOSE 9012
COPY --from=builder_products /minerva/bin/products ./minerva-products
CMD ["./minerva-products"]
