package view

import (
	"github.com/gin-gonic/gin"
	
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

func InstallRoutes(router *gin.Engine) {
	log.Info("Installing company")
	router.GET("/", CompanyIndexView)
	router.GET("/companies", CompanyTableLinesView)
	router.POST("/companies/new", NewCompanyView)
	router.DELETE("/companies/:id", DeleteCompanyView)
	router.GET("/:company/menu", CompanyMenuView)
	router.GET("/:company/detail", CompanyDetailView)

	log.Info("Installing User")
	router.GET("/:company/users", UsersIndexView)
	router.GET("/:company/users/list", UsersListView)
	router.DELETE("/:company/users/:id", UsersDeleteView)
}
