package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sloggin "github.com/samber/slog-gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/Minerva-System/minerva-go/cmd/rest/docs"
	connection "github.com/Minerva-System/minerva-go/internal/connection"
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

// @title Minerva System API
// @version 1.0
// @description Minerva System API
// @BasePath /
// @query.collection.format multi

func main() {
	godotenv.Load()
	log.Init()
	gin.SetMode(gin.ReleaseMode)
	
	log.Info("Minerva System: REST gateway service (Go port)")
	log.Info("Copyright (c) 2022-2023 Lucas S. Vieira")

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
	
	host := ":9000"

	router := gin.New()
	router.Use(sloggin.New(slog.Default()))
	router.Use(gin.Recovery())
	
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


	// Example route for managing gRPC connections
	router.GET("/test", func(ctx *gin.Context) {
		log.Debug("Retrieving a user service worker...")
		conn, err := col.UserSvc.Get(ctx)
		if err != nil {
			log.Error("Failed to retrieve a user service worker: %v", err)
			ctx.JSON(500, gin.H{"status": 500, "message": "Could not connect to user service"})
			return
		}
		defer conn.Close() // Very important!

		client := rpc.NewUserClient(conn)

		var page int64 = 0
		response, err := client.Index(ctx, &rpc.PageIndex{ Index: &page })
		if err != nil {
			log.Error("Failed to retrieve user index: %v", err)
			ctx.JSON(500, gin.H{"status": 500, "message": "Could not connect to user service"})
			return
		}

		ctx.JSON(200, response)
	})
	
	router.Run(host)
}
