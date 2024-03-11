package main

import (
	"net"
	"github.com/joho/godotenv"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	svc "github.com/Minerva-System/minerva-go/internal/svc/user/service"
)

func main() {
	godotenv.Load()
	log.Init()
	
	log.Info("Minerva System: USER service (Go port)")
	log.Info("Copyright (c) 2022-2024 Lucas S. Vieira")

	server := svc.CreateServer()

	log.Info("Binding TCP port...")
	listener, err := net.Listen("tcp", ":9010")
	if err != nil {
		log.Fatal("Failed to bind gRPC server port: %v", err)
	}

	log.Info("Initializing gRPC server.")
	server.Serve(listener)
}
