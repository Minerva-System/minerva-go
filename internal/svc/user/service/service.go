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
	controller "github.com/Minerva-System/minerva-go/internal/svc/user/controller"
)

type UserServerImpl struct {
	rpc.UnimplementedUserServer
	conn connection.Collection
}

func (self UserServerImpl) Index(ctx context.Context, idx *rpc.TenantPageIndex) (*rpc.UserList, error) {
	log.Info("Index method called")

	if (idx == nil) || (idx.Index == nil) {
		log.Error("Index method failed: Missing page index parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing page index parameter")
	}
	
	return controller.UserIndex(self.conn.DB, idx.CompanyId, *idx.Index)
}

func (self UserServerImpl) Show(ctx context.Context, idx *rpc.TenantEntityIndex) (*rpc.User, error) {
	log.Info("Show method called")

	if idx == nil {
		log.Error("Show method failed: missing id parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing id parameter")
	}

	return controller.GetUser(self.conn.DB, idx.CompanyId, idx.Index)
}

func (self UserServerImpl) Store(ctx context.Context, user *rpc.User) (*rpc.User, error) {
	log.Info("Store method called")

	if user == nil {
		log.Error("Store method failed: missing new user data")
		return nil, status.Error(codes.InvalidArgument, "Missing new user data")
	}
	
	return controller.CreateUser(self.conn.DB, user)
}

func (self UserServerImpl) Update(ctx context.Context, user *rpc.User) (*rpc.User, error) {
	log.Info("Update method called")

	if user == nil {
		log.Error("Update method failed: missing user data")
		return nil, status.Error(codes.InvalidArgument, "Missing user data")
	}

	return controller.UpdateUser(self.conn.DB, user)
}

func (self UserServerImpl) Delete(ctx context.Context, idx *rpc.TenantEntityIndex) (*emptypb.Empty, error) {
	log.Info("Delete method called")

	if idx == nil {
		log.Error("Delete method failed: missing id parameter")
		return nil, status.Error(codes.InvalidArgument, "Missing id parameter")
	}
	
	return &emptypb.Empty{}, controller.DeleteUser(self.conn.DB, idx.CompanyId, idx.Index)
}


func ApplyMigrations(col *connection.Collection) {
	log.Info("Migrating user table...")
	if err := col.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Error while migrating database: %v", err)
	}
}

func CreateServer() *grpc.Server {
	log.Info("Initializing user server...")

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

	s := &UserServerImpl{
		conn: col,
	}

	log.Info("User server preparations complete.")

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterUserServer(grpcServer, s)
	
	return grpcServer
}
