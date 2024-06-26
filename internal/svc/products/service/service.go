package service

import (
	"context"

	grpc "google.golang.org/grpc"
	status "google.golang.org/grpc/status"
	codes "google.golang.org/grpc/codes"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	connection "github.com/Minerva-System/minerva-go/internal/connection"
	log "github.com/Minerva-System/minerva-go/pkg/log"

	controller "github.com/Minerva-System/minerva-go/internal/svc/products/controller"
)

type ProductsServerImpl struct {
	rpc.UnimplementedProductsServer
	conn connection.Collection
}

func (self ProductsServerImpl) Index(ctx context.Context, idx *rpc.TenantPageIndex) (*rpc.ProductList, error) {
	log.Info("Index method called")

	if (idx == nil) || (idx.Index == nil) {
		log.Error("Index method failed: Missing page index parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing page index parameter")
	}

	return controller.ProductsIndex(self.conn.DB, idx.CompanyId, *idx.Index)
}

func (self ProductsServerImpl) Show(ctx context.Context, idx *rpc.TenantEntityIndex) (*rpc.Product, error) {
	log.Info("Show method called")

	if idx == nil {
		log.Error("Show method failed: missing id parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing id parameter")
	}

	return controller.GetProduct(self.conn.DB, idx.CompanyId, idx.Index)
}

func (self ProductsServerImpl) Store(ctx context.Context, product *rpc.Product) (*rpc.Product, error) {
	log.Info("Store method called")

	if product == nil {
		log.Error("Store method failed: missing new product data")
		return nil, status.Error(codes.InvalidArgument, "Missing new product data")
	}

	return controller.CreateProduct(self.conn.DB, product)
}

func (self ProductsServerImpl) Update(ctx context.Context, product *rpc.Product) (*rpc.Product, error) {
	log.Info("Update method called")

	if product == nil {
		log.Error("Update method failed: missing product data")
		return nil, status.Error(codes.InvalidArgument, "Missing product data")
	}
	return controller.UpdateProduct(self.conn.DB, product)
}

func (self ProductsServerImpl) Delete(ctx context.Context, idx *rpc.TenantEntityIndex) (*emptypb.Empty, error) {
	log.Info("Delete method called")

	if idx == nil {
		log.Error("Delete method failed: missing id parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing id parameter")
	}

	return &emptypb.Empty{}, controller.DeleteProduct(self.conn.DB, idx.CompanyId, idx.Index)
}

func CreateServer() *grpc.Server {
	log.Info("Initializing products server...")

	log.Info("Establishing connections...")
	col, err := connection.NewCollection(
		connection.CollectionOptions{
			WithDatabase: true,
			WithMessageBroker: true,
		})

	if err != nil {
		log.Fatal("Failed to establish connections: %v", err)
	}
	
	s := &ProductsServerImpl{
		conn: col,
	}

	log.Info("Products server preparations complete.")

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterProductsServer(grpcServer, s)
	
	return grpcServer
}
