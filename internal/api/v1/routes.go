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

	// Tenant

	/* Tenant-specific routes */
	tenant := api.Group("/:company", server.TenantCheckMiddleware())

	// Users
	tenant.GET("/users", server.GetUsers)
	tenant.GET("/users/:id", server.GetUser)
	tenant.POST("/users", server.CreateUser)
	tenant.DELETE("/users/:id", server.DeleteUser)

	// Products
	tenant.GET("/products", server.GetProducts)
	tenant.GET("/products/:id", server.GetProduct)
	tenant.POST("/products", server.CreateProduct)
	tenant.DELETE("/products/:id", server.DeleteProduct)
}
