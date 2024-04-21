package minerva_api_v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	model "github.com/Minerva-System/minerva-go/internal/model"
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	schema "github.com/Minerva-System/minerva-go/internal/schema"
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
func (self *Server) GetCompanies(ctx *gin.Context) {
	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "0"), 10, 64)
	if err != nil || page < 0 {
		log.Error("Could not parse page size")
		ctx.JSON(400, schema.ErrorMessage{
			Status:  400,
			Message: "Could not parse page index",
		})
		return
	}

	log.Debug("Retrieving a tenant service worker...")
	conn, err := self.Collection.TenantSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a tenant service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not connect to tenant service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewTenantClient(conn)
	response, err := client.Index(ctx, &rpc.PageIndex{Index: &page})
	if err != nil {
		log.Error("Failed to retrieve company index: %v", err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.Company{}.FromListMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved companies list: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not parse retrieved companies list",
		})
		return
	}

	ctx.JSON(200, res)
}

// @Summary Get company
// @Description Get data of a specific company by id
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     id    path    string    true    "company UUID"
// @Router  /companies/{id} [get]
func (self *Server) GetCompany(ctx *gin.Context) {
	id := ctx.Param("id")

	log.Debug("Retrieving a tenant service worker...")
	conn, err := self.Collection.TenantSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a tenant service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not connect to tenant service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewTenantClient(conn)
	response, err := client.Show(ctx, &rpc.EntityIndex{Index: id})
	if err != nil {
		log.Error("Failed to retrieve company %s: %v", id, err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.Company{}.FromMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved company: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not parse retrieved company",
		})
		return
	}

	ctx.JSON(200, res)
}

// @Summary Get company by slug
// @Description Get data of a specific company by unique identifier (slug)
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     slug    path    string    true    "company slug"
// @Router  /companies/by-slug/{slug} [get]
func (self *Server) GetCompanyBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	
	log.Debug("Retrieving a tenant service worker...")
	conn, err := self.Collection.TenantSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a tenant service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not connect to tenant service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewTenantClient(conn)
	response, err := client.ShowBySlug(ctx, &rpc.EntityIndex{Index: slug})
	if err != nil {
		log.Error("Failed to retrieve company with slug \"%s\": %v", slug, err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.Company{}.FromMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved company: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not parse retrieved company",
		})
		return
	}

	ctx.JSON(200, res)
}

// @Summary Check company existence
// @Description Check whether a specific company exists
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200     {object}    schema.BooleanResponse
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param    id    query    string    true    "company UUID"
// @Router  /companies/exists [get]
func (self *Server) GetCompanyExists(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		log.Error("No company id was informed")
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "No company id was informed",
		})
		return
	}
	
	log.Debug("Retrieving a tenant service worker...")
	conn, err := self.Collection.TenantSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a tenant service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not connect to tenant service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewTenantClient(conn)
	response, err := client.Exists(ctx, &rpc.EntityIndex{Index: id})
	if err != nil {
		log.Error("Failed to retrieve whether company exists: %v", err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	ctx.JSON(200, schema.BooleanResponse{}.FromMessage(response))
}

// @Summary Create company
// @Description Create a new company
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Param     data    body        schema.NewCompany    true    "new company data"
// @Success   201     {object}    model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router  /companies [post]
func (self *Server) CreateCompany(ctx *gin.Context) {
	var data schema.NewCompany
	if err := ctx.BindJSON(&data); err != nil {
		log.Error("Could not parse data from JSON")
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "Could not parse data into JSON",
		})
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Error("Error while validating data: %s", errors)
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "Error while validating data",
			Details: errors.Error(),
		})
		return
	}

	log.Debug("Retrieving a tenant service worker...")
	conn, err := self.Collection.TenantSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a tenant service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not connect to tenant service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewTenantClient(conn)
	msg := data.ToMessage()
	response, err := client.Store(ctx, &msg)

	res, err := model.Company{}.FromMessage(response)
	if err != nil {
		log.Error("Error while creating product: %v", err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	ctx.JSON(201, res)
}

// @Summary Update company
// @Description Update information of a specific company
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Param     data    body        schema.UpdatedCompany true    "company update data"
// @Success   200     {object}    model.Company
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     id    path    string    true    "company UUID"
// @Router /companies/{id} [put]
func (self *Server) UpdateCompany(ctx *gin.Context) {
	id := ctx.Param("id")

	var data schema.UpdatedCompany
	if err := ctx.BindJSON(&data); err != nil {
		log.Error("Could not parse data from JSON")
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "Could not parse data into JSON",
		})
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(data); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Error("Error while validating data: %s", errors)
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "Error while validating data",
			Details: errors.Error(),
		})
		return
	}

	if (data.Slug != "") && (len(data.Slug) < 3) {
		log.Error("Error while validating data: Slug has less than 3 characters")
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "Error while validating data",
			Details: "Slug has less than 3 characters",
		})
		return
	}

	if (data.CompanyName != "") && (len(data.CompanyName) < 3) {
		log.Error("Error while validating data: CompanyName has less than 3 characters")
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "Error while validating data",
			Details: "CompanyName has less than 3 characters",
		})
		return
	}

	if (data.TradingName != "") && (len(data.TradingName) < 3) {
		log.Error("Error while validating data: TradingName has less than 3 characters")
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "Error while validating data",
			Details: "TradingName has less than 3 characters",
		})
		return
	}

	log.Debug("Retrieving a tenant service worker...")
	conn, err := self.Collection.TenantSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a tenant service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not connect to tenant service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewTenantClient(conn)
	msg := data.ToMessage(id)
	response, err := client.Update(ctx, &msg)

	res, err := model.Company{}.FromMessage(response)
	if err != nil {
		log.Error("Error while updating company %s: %v", id, err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	ctx.JSON(200, res)
}

// @Summary Disable company
// @Description Delete a specific company (disabling is a soft-delete).
// @Tags Tenant
// @Accept    json
// @Produce   json
// @Success   200   "disabled successfully"
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Param     id    path    string    true    "company UUID"
// @Router  /companies/{id} [delete]
func (self *Server) DisableCompany(ctx *gin.Context) {
	id := ctx.Param("id")

	log.Debug("Retrieving a tenant service worker...")
	conn, err := self.Collection.TenantSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a tenant service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not connect to tenant service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewTenantClient(conn)
	_, err = client.Disable(ctx, &rpc.EntityIndex{Index: id})
	if err != nil {
		log.Error("Failed to disable company %s: %v", id, err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	ctx.Status(200)
}
