package main

import (
	"log"
	
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
	log.Print("Hello world!")
	host := "0.0.0.0:9000"

	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run(host)
}
