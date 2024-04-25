package view

import (
	"github.com/gin-gonic/gin"
	
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

func InstallRoutes(router *gin.Engine) {
	log.Info("Installing HTML templates")
	router.LoadHTMLGlob("templates/*/*.html")

	log.Info("Installing company")
	router.GET("/", CompanyIndexView)
	router.GET("/companies", CompanyTableLinesView)
	router.POST("/companies/new", NewCompanyView)
	router.DELETE("/companies/:id", DeleteCompanyView)
}
