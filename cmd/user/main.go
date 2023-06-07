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

type UserServerImpl struct {
	rpc.UnimplementedUserServer
}

func (UserServerImpl) Index(context.Context, *rpc.PageIndex) (*rpc.UserList, error) {
	log.Print("Index method called")
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}
func (UserServerImpl) Show(context.Context, *rpc.EntityIndex) (*rpc.User, error) {
	log.Print("Show method called")
	return nil, status.Errorf(codes.Unimplemented, "method Show not implemented")
}
func (UserServerImpl) Store(context.Context, *rpc.User) (*rpc.User, error) {
	log.Print("Store method called")
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}
func (UserServerImpl) Update(context.Context, *rpc.User) (*rpc.User, error) {
	log.Print("Update method called")
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UserServerImpl) Delete(context.Context, *rpc.EntityIndex) (*emptypb.Empty, error) {
	log.Print("Delete method called")
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
