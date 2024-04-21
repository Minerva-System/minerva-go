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
	controller "github.com/Minerva-System/minerva-go/internal/svc/tenant/controller"
)

type TenantServerImpl struct {
	rpc.UnimplementedTenantServer
	conn connection.Collection
}

func (self TenantServerImpl) Index(ctx context.Context, idx *rpc.PageIndex) (*rpc.CompanyList, error) {
	log.Info("Index method called")

	if (idx == nil) || (idx.Index == nil) {
		log.Error("Index method failed: Missing page index parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing page index parameter")
	}

	return controller.CompanyIndex(self.conn.DB, *idx.Index)
}

func (self TenantServerImpl) Show(ctx context.Context, idx *rpc.EntityIndex) (*rpc.Company, error) {
	log.Info("Show method called")

	if idx == nil {
		log.Error("Show method failed: missing id parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing id parameter")
	}

	return controller.GetCompany(self.conn.DB, idx.Index)
}

func (self TenantServerImpl) ShowBySlug(ctx context.Context, idx *rpc.EntityIndex) (*rpc.Company, error) {
	log.Info("ShowBySlug method called")

	if idx == nil {
		log.Error("ShowBySlug method failed: missing slug parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing slug parameter")
	}

	return controller.GetCompanyBySlug(self.conn.DB, idx.Index)
}

func (self TenantServerImpl) Exists(ctx context.Context, idx *rpc.EntityIndex) (*wrapperspb.BoolValue, error) {
	log.Info("Exists method called")

	if idx == nil {
		log.Error("Exists method failed: missing id parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing id parameter")
	}

	return controller.GetCompanyExists(self.conn.DB, idx.Index)
}

func (self TenantServerImpl) Store(ctx context.Context, company *rpc.Company) (*rpc.Company, error) {
	log.Info("Store method called")

	if company == nil {
		log.Error("Store method failed: missing new company data")
		return nil, status.Error(codes.InvalidArgument, "Missing new company data")
	}

	return controller.CreateCompany(self.conn.DB, company)
}

func (self TenantServerImpl) Update(ctx context.Context, company *rpc.Company) (*rpc.Company, error) {
	log.Info("Update method called")

	if company == nil {
		log.Error("Update method failed: missing new company data")
		return nil, status.Error(codes.InvalidArgument, "Missing new company data")
	}

	return controller.UpdateCompany(self.conn.DB, company)
}

func (self TenantServerImpl) Disable(ctx context.Context, idx *rpc.EntityIndex) (*emptypb.Empty, error) {
	log.Info("Disable method called")

	if idx == nil {
		log.Error("Disable method failed: missing id parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing id parameter")
	}

	return &emptypb.Empty{}, controller.DisableCompany(self.conn.DB, idx.Index)
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
