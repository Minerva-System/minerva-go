package main

import (
	"log"
	"net"

	context "context"
	grpc "google.golang.org/grpc"
	status "google.golang.org/grpc/status"
	codes "google.golang.org/grpc/codes"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	
	rpc "minervarpc"
)

type ProductsServerImpl struct {
	rpc.UnimplementedProductsServer
}

func (ProductsServerImpl) Index(ctx context.Context, idx *rpc.PageIndex) (*rpc.ProductList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}

func (ProductsServerImpl) Show(ctx context.Context, idx *rpc.EntityIndex) (*rpc.Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Show not implemented")
}

func (ProductsServerImpl) Store(ctx context.Context, product *rpc.Product) (*rpc.Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}

func (ProductsServerImpl) Update(ctx context.Context, product *rpc.Product) (*rpc.Product, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (ProductsServerImpl) Delete(ctx context.Context, idx *rpc.EntityIndex) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func createServer() *ProductsServerImpl {
	s := &ProductsServerImpl{}
	return s
}

func main() {
	log.Print("Hello world!")
	listener, err := net.Listen("tcp", "0.0.0.0:9012")
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterProductsServer(grpcServer, createServer())
	grpcServer.Serve(listener)
}
