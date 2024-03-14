package minerva_api_v1

import (
	"github.com/gin-gonic/gin"
	_ "github.com/Minerva-System/minerva-go/cmd/rest/docs"
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
	router.GET("/api/v1/user", server.GetUsers)
}
