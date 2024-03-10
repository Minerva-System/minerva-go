package main

import (
	"log/slog"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	docs "github.com/Minerva-System/minerva-go/cmd/rest/docs"
)

// @title Minerva System API
// @version 1.0
// @description Minerva System API
// @BasePath /
// @query.collection.format multi

func main() {
	log.Init()
	gin.SetMode(gin.ReleaseMode)
	
	log.Info("Minerva System: REST gateway service (Go port)")
	log.Info("Copyright (c) 2022-2023 Lucas S. Vieira")

	godotenv.Load()
	
	host := ":9000"

	router := gin.New()
	router.Use(sloggin.New(slog.Default()))
	router.Use(gin.Recovery())
	
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run(host)
}
