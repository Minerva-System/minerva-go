package minerva_api_v1

import (
	// "strconv"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"

	// log "github.com/Minerva-System/minerva-go/pkg/log"
	// model "github.com/Minerva-System/minerva-go/internal/model"
	// rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	// schema "github.com/Minerva-System/minerva-go/internal/schema"
)

// @Summary Company list
// @Description Get a list of companies per page
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    []model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     page    query    int    false    "page number (0 or more)"
// @Router  /companies [get]
func (self *Server) GetCompanies(ctx *gin.Context) {}

// @Summary
// @Description
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     id    path    string    true    "company UUID"
// @Router  /companies/{id} [get]
func (self *Server) GetCompany(ctx *gin.Context) {}

// @Summary
// @Description
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     slug    path    string    true    "company slug"
// @Router  /companies/by-slug/{slug} [get]
func (self *Server) GetCompanyBySlug(ctx *gin.Context) {}

// @Summary
// @Description
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    schema.BooleanResponse
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param    id    query    string    true    "company UUID"
// @Router  /companies/exists [get]
func (self *Server) GetCompanyExists(ctx *gin.Context) {}

// @Summary
// @Description
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router  /companies [post]
func (self *Server) CreateCompany(ctx *gin.Context) {}

// @Summary
// @Description
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     id    path    string    true    "company UUID"
// @Router /companies/{id} [put]
func (self *Server) UpdateCompany(ctx *gin.Context) {}

// @Summary
// @Description
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200   "disabled successfully"
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     id    path    string    true    "company UUID"
// @Router  /companies/{id} [delete]
func (self *Server) DisableCompany(ctx *gin.Context) {}
