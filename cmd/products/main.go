package main

import (
	"os"
	"fmt"
	"net"

	"github.com/joho/godotenv"
	
	log "github.com/Minerva-System/minerva-go/pkg/log"
	svc "github.com/Minerva-System/minerva-go/internal/svc/products/service"
)

func main() {
	godotenv.Load()
	log.Init()
	
	log.Info("Minerva System: PRODUCTS service (Go port), Copyright (c) 2022-2024 Lucas S. Vieira")

	server := svc.CreateServer()
	var port string = ":9012"
	if p, exists := os.LookupEnv("MINERVA_PRODUCTS_PORT"); exists {
		port = fmt.Sprintf(":%s", p)
	} else {
		log.Warn("MINERVA_PRODUCTS_PORT not defined, using default port %s", port)
	}

	log.Info("Binding TCP port...")
	listener, err := net.Listen("tcp", ":9012")
	if err != nil {
		log.Fatal("Failed to bind gRPC server port: %v", err)
	}

	log.Info("Initializing gRPC server.")
	server.Serve(listener)
}
