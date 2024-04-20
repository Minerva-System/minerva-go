package minerva_api_v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	model "github.com/Minerva-System/minerva-go/internal/model"
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	schema "github.com/Minerva-System/minerva-go/internal/schema"
	log "github.com/Minerva-System/minerva-go/pkg/log"
)

// @Summary Products list
// @Description Get a list of products per page
// @Tags      Products
// @Accept    json
// @Produce   json
// @Param     company    path    string    true    "company UUID"
// @Param     page    query    int    false    "page number (0 or more)"
// @Success   200     {object}    []model.Product
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/products [get]
func (self *Server) GetProducts(ctx *gin.Context) {
	companyId := ctx.Param("company")
	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "0"), 10, 64)
	if err != nil || page < 0 {
		log.Error("Could not parse page size")
		ctx.JSON(400, schema.ErrorMessage{
			Status:  400,
			Message: "Could not parse page index",
		})
		return
	}

	log.Debug("Retrieving a products service worker...")
	conn, err := self.Collection.ProductsSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a products service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to products service",
		})
		return
	}
	defer conn.Close() // Very important!

	client := rpc.NewProductsClient(conn)
	response, err := client.Index(ctx, &rpc.TenantPageIndex{CompanyId: companyId, Index: &page})
	if err != nil {
		log.Error("Failed to retrieve product index: %v", err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.Product{}.FromListMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved products list: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not parse retrieved products list",
		})
	}

	ctx.JSON(200, res)
}

// @Summary Get product
// @Description Get data of a specific product
// @Tags      Products
// @Accept    json
// @Produce   json
// @Param     company    path    string    true    "company UUID"
// @Param     id    path    string    true    "product UUID"
// @Success   200     {object}    model.Product
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/products/{id} [get]
func (self *Server) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	companyId := ctx.Param("company")

	log.Debug("Retrieving a products service worker...")
	conn, err := self.Collection.ProductsSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a products service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to products service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewProductsClient(conn)
	response, err := client.Show(ctx, &rpc.TenantEntityIndex{CompanyId: companyId, Index: id})
	if err != nil {
		log.Error("Failed to retrieve product %s: %v", id, err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.Product{}.FromMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved product: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not parse retrieved product",
		})
	}

	ctx.JSON(200, res)
}

// @Summary Create product
// @Description Create a new product
// @Tags      Products
// @Accept    json
// @Produce   json
// @Param     company    path    string    true    "company UUID"
// @Param     data    body        schema.NewProduct    true    "new product data"
// @Success   201     {object}    model.Product
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/products [post]
func (self *Server) CreateProduct(ctx *gin.Context) {
	companyId := ctx.Param("company")

	var data schema.NewProduct
	if err := ctx.BindJSON(&data); err != nil {
		ctx.JSON(400, schema.ErrorMessage{
			Status:  400,
			Message: "Could not parse data into JSON",
		})
	}

	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Error("Error while validating data: %s", errors)
		ctx.JSON(400, schema.ErrorMessage{
			Status:  400,
			Message: "Error while validating data",
			Details: errors.Error(),
		})
		return
	}

	log.Debug("Retrieving a products service worker...")
	conn, err := self.Collection.ProductsSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a products service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to products service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewProductsClient(conn)
	msg := data.ToMessage(companyId)
	response, err := client.Store(ctx, &msg)

	res, err := model.Product{}.FromMessage(response)
	if err != nil {
		log.Error("Error while creating product: %v", err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	ctx.JSON(201, res)
}

// @Summary Delete product
// @Description Delete a specific product
// @Tags      Products
// @Accept    json
// @Param     company    path    string    true    "company UUID"
// @Param     id    path    string    true    "product UUID"
// @Success   200   "deleted successfully"
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/products/{id} [delete]
func (self *Server) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	companyId := ctx.Param("company")

	log.Debug("Retrieving a products service worker...")
	conn, err := self.Collection.ProductsSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a products service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to products service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewProductsClient(conn)
	_, err = client.Delete(ctx, &rpc.TenantEntityIndex{CompanyId: companyId, Index: id})
	if err != nil {
		log.Error("Failed to delete product %s: %v", id, err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	ctx.Status(200)
}
