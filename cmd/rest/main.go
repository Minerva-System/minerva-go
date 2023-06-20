package main

import (
	"log"

	"github.com/joho/godotenv"
	
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "minervarestdocs"
)

// @title Minerva System API
// @version 1.0
// @description Minerva System API
// @BasePath /
// @query.collection.format multi

func main() {
	log.Print("Minerva System: REST gateway service (Go port)")
	log.Print("Copyright (c) 2022-2023 Lucas S. Vieira")

	godotenv.Load()
	
	host := "0.0.0.0:9000"

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(host)
}
