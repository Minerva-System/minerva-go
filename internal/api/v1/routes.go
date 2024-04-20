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
// @license.name MIT
// @license.url https://choosealicense.com/licenses/mit/


func InstallRoutes(router *gin.Engine, server *Server) {
	api := router.Group("/api/v1")
	api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Users
	api.GET("/:company/users", server.GetUsers)
	api.GET("/:company/users/:id", server.GetUser)
	api.POST("/:company/users", server.CreateUser)
	api.DELETE("/:company/users/:id", server.DeleteUser)

	// Products
	api.GET("/:company/products", server.GetProducts)
	api.GET("/:company/products/:id", server.GetProduct)
	api.POST("/:company/products", server.CreateProduct)
	api.DELETE("/:company/products/:id", server.DeleteProduct)

	// Tenant
	
}
