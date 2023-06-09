package main

import (
	"fmt"
	"log"
	"net"
	"os"

	context "context"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"golang.org/x/crypto/bcrypt"

	model "minervamodel"
	rpc "minervarpc"
)

type UserServerImpl struct {
	rpc.UnimplementedUserServer
	db *gorm.DB
	// Pools for database and rabbitmq
	// connections map[string](DB, RABBITMQ)
}

func (self UserServerImpl) Index(ctx context.Context, idx *rpc.PageIndex) (*rpc.UserList, error) {
	log.Print("Index method called")
	log.Printf("Payload: %s", idx)

	// Test
	log.Print("Migrating user table...")
	if err := self.db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Error while migrating database: %v", err)
	}
	
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}

func (self UserServerImpl) Show(ctx context.Context, idx *rpc.EntityIndex) (*rpc.User, error) {
	log.Print("Show method called")
	log.Printf("Payload: %s", idx)

	// Test
	log.Print("Creating user...")

	hash, err := bcrypt.GenerateFromPassword([]byte("123456"), 8)
	if err != nil {
		log.Fatalf("Unable to generate password hash: %v", err)
	}

	u := model.User{
		Login: "fulano",
		Name: "Fulano de Tal",
		Email: nil,
		Pwhash: hash,
	}

	result := self.db.Create(&u)
	if result.Error != nil {
		log.Fatalf("Unable to create user: %v", result.Error)
	}

	log.Printf("User created with ID %d (rows affected: %d)", u.ID, result.RowsAffected)
	
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
	log.Print("Initializing server...")
	var dbsrv string
	// var amqpsrv string
	var exists bool
	
	if dbsrv, exists = os.LookupEnv("DATABASE_SERVICE_SERVER"); !exists {
		log.Fatal("Unable to read DATABASE_SERVICE_SERVER")
	}
	
	// if amqpsrv, exists = os.LookupEnv("RABBITMQ_SERVICE_SERVER"); !exists {
	// 	log.Fatal("Unable to read RABBITMQ_SERVICE_SERVER")
	// }

	dsn := fmt.Sprintf(
		"minerva:mysql@tcp(%s)/minerva?charset=utf8&parseTime=True&loc=Local",
		dbsrv)

	// Connect to database
	log.Printf("Connecting to database...")
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
		DefaultStringSize: 256,
		DontSupportRenameIndex: true,
		DontSupportRenameColumn: true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				Colorful: false,
			},
		),
	})

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Print("Connected to database!")

	s := &UserServerImpl{}
	s.db = db
	
	return s
}

func main() {
	log.Print("Minerva System: USER service (Go port)")
	log.Print("Copyright (c) 2022-2023 Lucas S. Vieira")
	log.Print()

	godotenv.Load()

	server := createServer()
	
	listener, err := net.Listen("tcp", "0.0.0.0:9010")
	if err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterUserServer(grpcServer, server)
	grpcServer.Serve(listener)
}
