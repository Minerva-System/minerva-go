package main

import (
	"log"
	"net"

	"github.com/joho/godotenv"

	context "context"
	grpc "google.golang.org/grpc"
	status "google.golang.org/grpc/status"
	codes "google.golang.org/grpc/codes"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	
	rpc "minervarpc"
)

type SessionServerImpl struct {
	rpc.UnimplementedSessionServer
}

func (SessionServerImpl) Generate(ctx context.Context, session *rpc.SessionCreationData) (*rpc.SessionToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Generate not implemented")
}

func (SessionServerImpl) Retrieve(ctx context.Context, token *rpc.SessionToken) (*rpc.SessionData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Retrieve not implemented")
}

func (SessionServerImpl) Remove(ctx context.Context, token *rpc.SessionToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

func createServer() *SessionServerImpl {
	s := &SessionServerImpl{}
	return s
}

func main() {
	log.Print("Minerva System: SESSION service (Go port)")
	log.Print("Copyright (c) 2022-2023 Lucas S. Vieira")

	godotenv.Load()
	
	listener, err := net.Listen("tcp", "0.0.0.0:9011")
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterSessionServer(grpcServer, createServer())
	grpcServer.Serve(listener)
}
