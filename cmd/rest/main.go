package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sloggin "github.com/samber/slog-gin"
	"github.com/gin-contrib/cors"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Minerva-System/minerva-go/cmd/rest/docs"
	connection "github.com/Minerva-System/minerva-go/internal/connection"
	log "github.com/Minerva-System/minerva-go/pkg/log"
	api_v1 "github.com/Minerva-System/minerva-go/internal/api/v1"
)

//go:generate swag init -d ../../internal/api/v1,../../internal/model -g routes.go

func main() {
	godotenv.Load()
	log.Init()
	gin.SetMode(gin.ReleaseMode)
	
	log.Info("Minerva System: REST gateway service (Go port)")
	log.Info("Copyright (c) 2022-2024 Lucas S. Vieira")

	log.Info("Establishing connections...")

	// Test
	col, err := connection.NewCollection(connection.CollectionOptions{
		WithUserService: true,
		WithSessionService: true,
		WithProductsService: true,
	})
	if err != nil {
		log.Fatal("Failed to establish connections: %v", err)
	}

	server := api_v1.Server{
		Collection: col,
	}
	
	host := ":9000"

	router := gin.New()
	router.Use(sloggin.New(slog.Default()))
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	
	// docs.SwaggerInfo.BasePath = "/"
	router.GET("/api/v1/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api_v1.InstallRoutes(router, &server)
	
	router.Run(host)
}
