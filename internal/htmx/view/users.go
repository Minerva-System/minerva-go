package view

import (
	// "encoding/json"
	
	"github.com/gin-gonic/gin"

	service "github.com/Minerva-System/minerva-go/internal/htmx/service"
	log "github.com/Minerva-System/minerva-go/pkg/log"
	config "github.com/Minerva-System/minerva-go/internal/htmx/config"
)

func UsersIndexView(c *gin.Context) {
	companyId := c.Param("company")
	c.HTML(200, "users/main", gin.H{
		"api_host": config.Values.FullHost,
		"company_id": companyId,
		"title":    "Users",
	})
}

func UsersListView(c *gin.Context) {
	companyId := c.Param("company")
	users, err := service.GetUsers(companyId)
	if err != nil {
		log.Error("Error fetching users: %v", err)
		c.AbortWithError(500, err)
		return
	}
	
	c.HTML(200, "users/list", gin.H{
		"api_host": config.Values.FullHost,
		"company_id": companyId,
		"users": users,
	})
}

func UsersDeleteView(c *gin.Context) {
	// TODO: inject company_id on params
	UsersListView(c)
}
