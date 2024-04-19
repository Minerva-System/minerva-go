package main

import (
	"os"
	"fmt"
	"net"

	"github.com/joho/godotenv"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	svc "github.com/Minerva-System/minerva-go/internal/svc/tenant/service"
)

func main() {
	godotenv.Load()
	log.Init()

	log.Info("Minerva System: TENANT service (Go port), Copyright (c) 2022-2024 Lucas S. Vieira")

	server := svc.CreateServer()
	var port string = ":9013"

	if p, exists := os.LookupEnv("MINERVA_TENANT_PORT"); exists {
		port = fmt.Sprintf(":%s", p)
	} else {
		log.Warn("MINERVA_TENANT_PORT not defined, using default port %s", port)
	}

	log.Info("Binding TCP port...")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to bind gRPC server port: %v", err)
	}

	log.Info("Initializing gRPC server.")
	server.Serve(listener)
}
