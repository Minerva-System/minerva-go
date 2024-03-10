package main

import (
	"net"

	context "context"

	"github.com/joho/godotenv"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"golang.org/x/crypto/bcrypt"

	
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	connection "github.com/Minerva-System/minerva-go/internal/connection"
	model "github.com/Minerva-System/minerva-go/internal/model"
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

type UserServerImpl struct {
	rpc.UnimplementedUserServer
	conn connection.Collection
}

func (self UserServerImpl) Index(ctx context.Context, idx *rpc.PageIndex) (*rpc.UserList, error) {
	log.Info("Index method called")
	log.Info("Payload: %s", idx)
	
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}

func (self UserServerImpl) Show(ctx context.Context, idx *rpc.EntityIndex) (*rpc.User, error) {
	log.Info("Show method called")
	log.Info("Payload: %s", idx)

	// Test
	log.Info("Creating user...")

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), 8)
	if err != nil {
		log.Error("Unable to generate password hash: %v", err)
		return nil, err
	}

	u := model.User{
		Login: "fulano",
		Name: "Fulano de Tal",
		Email: nil,
		Pwhash: hash,
	}

	result := self.conn.DB.Create(&u)
	if result.Error != nil {
		log.Error("Unable to create user: %v", result.Error)
		return nil, status.Errorf(codes.Aborted, "Unable to create user: %v", result.Error)
	}

	log.Info("User created with ID %d (rows affected: %d)", u.ID, result.RowsAffected)
	
	return nil, status.Errorf(codes.Unimplemented, "method Show not implemented")
}

func (self UserServerImpl) Store(ctx context.Context, user *rpc.User) (*rpc.User, error) {
	log.Info("Store method called")

	log.Info("Serializing message to model...")
	db_user, err := model.UserFromMessage(user)
	if err != nil {
		log.Error("Error while converting message to model: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Error while converting message to model: %v", err)
	}

	log.Info("Saving to database...")
	result := self.conn.DB.Create(&db_user)
	if result.Error != nil {
		log.Error("Unable to create user: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "Unable to create user: %v", result.Error)
	}

	log.Info("User created. ID: %s", db_user.ID)
	
	new_user := db_user.ToMessage()
	return &new_user, nil
}

func (UserServerImpl) Update(ctx context.Context, user *rpc.User) (*rpc.User, error) {
	log.Info("Update method called")
	log.Info("Payload: %s", user)
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (UserServerImpl) Delete(ctx context.Context, idx *rpc.EntityIndex) (*emptypb.Empty, error) {
	log.Info("Delete method called")
	log.Info("Payload: %s", idx)
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func createServer() *UserServerImpl {
	log.Info("Initializing server...")

	log.Info("Establishing connections...")
	col, err := connection.NewCollection(
		connection.CollectionOptions{
			WithDatabase: true,
			WithMessageBroker: true,
		})

	if err != nil {
		log.Fatal("Failed to establish connections: %v", err)
	}

	log.Info("Migrating user table...")
	if err := col.DB.AutoMigrate(&model.User{}); err != nil {
		log.Fatal("Error while migrating database: %v", err)
	}

	s := &UserServerImpl{
		conn: col,
	}

	log.Info("Server preparations complete.")
	return s
}

func main() {
	log.Init()
	
	log.Info("Minerva System: USER service (Go port)")
	log.Info("Copyright (c) 2022-2024 Lucas S. Vieira")

	godotenv.Load()

	server := createServer()

	log.Info("Allocating TCP port...")
	listener, err := net.Listen("tcp", ":9010")
	if err != nil {
		log.Fatal("Failed to start gRPC server: %v", err)
	}

	log.Info("Initializing gRPC server.")
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterUserServer(grpcServer, server)
	grpcServer.Serve(listener)
}
