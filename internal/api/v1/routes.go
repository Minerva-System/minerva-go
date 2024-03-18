package minerva_api_v1

import (
	"github.com/gin-gonic/gin"
	_ "github.com/Minerva-System/minerva-go/cmd/rest/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Minerva System API
// @version 1.0
// @description Minerva System API (v1)
// @contact.name Lucas S. Vieira
// @contact.url https://luksamuk.codes
// @contact.email lucasvieira@protonmail.com
// @BasePath /api/v1
// @query.collection.format multi

// TODO: license.name, license.url

func InstallRoutes(router *gin.Engine, server *Server) {
	api := router.Group("/api/v1")
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	
	api.GET("/users", server.GetUsers)
	api.GET("/users/:id", server.GetUser)
	api.POST("/users", server.CreateUser)
}
