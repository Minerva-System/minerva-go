package service

import (
	"context"

	grpc "google.golang.org/grpc"
	status "google.golang.org/grpc/status"
	codes "google.golang.org/grpc/codes"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	connection "github.com/Minerva-System/minerva-go/internal/connection"
	log "github.com/Minerva-System/minerva-go/pkg/log"

	model "github.com/Minerva-System/minerva-go/internal/model"
	// controller "github.com/Minerva-System/minerva-go/internal/svc/tenant/controller"
)

type TenantServerImpl struct {
	rpc.UnimplementedTenantServer
	conn connection.Collection
}

func (self TenantServerImpl) Index(context.Context, *rpc.PageIndex) (*rpc.CompanyList, error) {
	log.Info("Index method called")
	return nil, status.Error(codes.Unimplemented, "Method unimplemented")
}

func (self TenantServerImpl) Show(context.Context, *rpc.EntityIndex) (*rpc.Company, error) {
	log.Info("Show method called")
	return nil, status.Error(codes.Unimplemented, "Method unimplemented")
}

func (self TenantServerImpl) ShowBySlug(context.Context, *rpc.EntityIndex) (*rpc.Company, error) {
	log.Info("ShowBySlug method called")
	return nil, status.Error(codes.Unimplemented, "Method unimplemented")
}

func (self TenantServerImpl) Exists(context.Context, *rpc.EntityIndex) (*wrapperspb.BoolValue, error) {
	log.Info("Exists method called")
	return nil, status.Error(codes.Unimplemented, "Method unimplemented")
}

func (self TenantServerImpl) Store(context.Context, *rpc.Company) (*rpc.Company, error) {
	log.Info("Store method called")
	return nil, status.Error(codes.Unimplemented, "Method unimplemented")
}

func (self TenantServerImpl) Update(context.Context, *rpc.Company) (*rpc.Company, error) {
	log.Info("Update method called")
	return nil, status.Error(codes.Unimplemented, "Method unimplemented")
}

func (self TenantServerImpl) Disable(context.Context, *rpc.EntityIndex) (*emptypb.Empty, error) {
	log.Info("Disable method called")
	return nil, status.Error(codes.Unimplemented, "Method unimplemented")
}

func ApplyMigrations(col *connection.Collection) {
	log.Info("Migrating company table...")
	if err := col.DB.AutoMigrate(&model.Company{}); err != nil {
		log.Fatal("Error while migrating database: %v", err)
	}
}

func CreateServer() *grpc.Server {
	log.Info("Initializing tenant server...")

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

	s := &TenantServerImpl{
		conn: col,
	}

	log.Info("User server preparations complete.")

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterTenantServer(grpcServer, s)
	
	return grpcServer
}
