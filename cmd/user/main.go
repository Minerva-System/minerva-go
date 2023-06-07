package main

import (
	"log"
	"net"

	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	rpc "minervarpc"
)

type UserServerImpl struct {
	rpc.UnimplementedUserServer
}

func (UserServerImpl) Index(ctx context.Context, idx *rpc.PageIndex) (*rpc.UserList, error) {
	log.Print("Index method called")
	log.Printf("Payload: %s", idx)
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}

func (UserServerImpl) Show(ctx context.Context, idx *rpc.EntityIndex) (*rpc.User, error) {
	log.Print("Show method called")
	log.Printf("Payload: %s", idx)
	return nil, status.Errorf(codes.Unimplemented, "method Show not implemented")
}

func (UserServerImpl) Store(ctx context.Context, user *rpc.User) (*rpc.User, error) {
	log.Print("Store method called")
	log.Printf("Payload: %s", user)
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}

func (UserServerImpl) Update(ctx context.Context, user *rpc.User) (*rpc.User, error) {
	log.Print("Update method called")
	log.Printf("Payload: %s", user)
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (UserServerImpl) Delete(ctx context.Context, idx *rpc.EntityIndex) (*emptypb.Empty, error) {
	log.Print("Delete method called")
	log.Printf("Payload: %s", idx)
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func createServer() *UserServerImpl {
	s := &UserServerImpl{}
	return s
}

func main() {
	log.Print("Hello world!")
	listener, err := net.Listen("tcp", "0.0.0.0:9010")
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterUserServer(grpcServer, createServer())
	grpcServer.Serve(listener)
}
