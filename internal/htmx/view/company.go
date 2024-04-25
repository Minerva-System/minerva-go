package view

import (
	"encoding/json"
	
	"github.com/gin-gonic/gin"

	service "github.com/Minerva-System/minerva-go/internal/htmx/service"
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

type NewCompanyForm struct {
	Slug        string `form:"slug" json:"slug"`
	CompanyName string `form:"companyName" json:"companyName"`
	TradingName string `form:"tradingName" json:"tradingName"`
}

func CompanyIndexView(c *gin.Context) {
	c.HTML(200, "company/main", gin.H{
		"api_host": "http://localhost:5090",
		"title":    "Companies",
	})
}

func CompanyTableLinesView(c *gin.Context) {
	data, err := service.GetCompanies()
	if err != nil {
		log.Error("Error fetching companies: %v", err)
		c.AbortWithError(500, err)
		return
	}

	c.HTML(200, "company/list", data)
}

func NewCompanyView(c *gin.Context) {
	var data NewCompanyForm
	c.Bind(&data)
	json, err := json.Marshal(data)
	if err != nil {
		log.Error("Error marshalling form data to JSON: %v", err)
		c.AbortWithError(500, err)
		return
	}

	_, err = service.NewCompany(string(json))
	if err != nil {
		log.Error("Error while creating company: %v", err)
		c.AbortWithError(500, err)
		return
	}
	
	// Return company table items
	CompanyTableLinesView(c)
}

func DeleteCompanyView(c *gin.Context) {
	id := c.Param("id")
	err := service.DeleteCompany(id)
	if err != nil {
		log.Error("Error deleting companies: %v", err)
		c.AbortWithError(500, err)
		return
	}

	// Return company table items
	CompanyTableLinesView(c)
}
