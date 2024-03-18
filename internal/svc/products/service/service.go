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

	model "github.com/Minerva-System/minerva-go/internal/model"
	// controller "github.com/Minerva-System/minerva-go/internal/svc/products/controller"
)

type ProductsServerImpl struct {
	rpc.UnimplementedProductsServer
	conn connection.Collection
}

func (ProductsServerImpl) Index(ctx context.Context, idx *rpc.PageIndex) (*rpc.ProductList, error) {
	log.Info("Index method called")
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}

func (ProductsServerImpl) Show(ctx context.Context, idx *rpc.EntityIndex) (*rpc.Product, error) {
	log.Info("Show method called")
	return nil, status.Errorf(codes.Unimplemented, "method Show not implemented")
}

func (ProductsServerImpl) Store(ctx context.Context, product *rpc.Product) (*rpc.Product, error) {
	log.Info("Store method called")
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}

func (ProductsServerImpl) Update(ctx context.Context, product *rpc.Product) (*rpc.Product, error) {
	log.Info("Update method called")
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (ProductsServerImpl) Delete(ctx context.Context, idx *rpc.EntityIndex) (*emptypb.Empty, error) {
	log.Info("Delete method called")
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func ApplyMigrations(col *connection.Collection) {
	log.Info("Migrating product table...")
	if err := col.DB.AutoMigrate(&model.Product{}); err != nil {
		log.Fatal("Error while migrating database: %v", err)
	}
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

	ApplyMigrations(&col)
	
	s := &ProductsServerImpl{
		conn: col,
	}

	log.Info("Products server preparations complete.")

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterProductsServer(grpcServer, s)
	
	return grpcServer
}
