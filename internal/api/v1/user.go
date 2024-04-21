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

// @Summary User list
// @Description Get a list of users per page
// @Tags      User
// @Accept    json
// @Produce   json
// @Param     company    path    string    true    "company UUID"
// @Param     page    query    int    false    "page number (0 or more)"
// @Success   200     {object}    []model.User
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/users [get]
func (self *Server) GetUsers(ctx *gin.Context) {
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

	log.Debug("Retrieving a user service worker...")
	conn, err := self.Collection.UserSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a user service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to user service",
		})
		return
	}
	defer conn.Close() // Very important!

	client := rpc.NewUserClient(conn)
	response, err := client.Index(ctx, &rpc.TenantPageIndex{CompanyId: companyId, Index: &page})
	if err != nil {
		log.Error("Failed to retrieve user index: %v", err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.User{}.FromListMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved user list: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not parse retrieved user list",
		})
		return
	}

	ctx.JSON(200, res)
}

// @Summary Get user
// @Description Get data of a specific user
// @Tags      User
// @Accept    json
// @Produce   json
// @Param     company    path    string    true    "company UUID"
// @Param     id    path    string    true    "user UUID"
// @Success   200     {object}    model.User
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/users/{id} [get]
func (self *Server) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	companyId := ctx.Param("company")

	log.Debug("Retrieving a user service worker...")
	conn, err := self.Collection.UserSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a user service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to user service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewUserClient(conn)
	response, err := client.Show(ctx, &rpc.TenantEntityIndex{CompanyId: companyId, Index: id})
	if err != nil {
		log.Error("Failed to retrieve user %s: %v", id, err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.User{}.FromMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved user: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not parse retrieved user",
		})
		return
	}

	ctx.JSON(200, res)
}

// @Summary Create user
// @Description Create a new user
// @Tags      User
// @Accept    json
// @Produce   json
// @Param     company    path    string    true    "company UUID"
// @Param     data    body        schema.NewUser    true    "new user data"
// @Success   201     {object}    model.User
// @Failure   400     {object}    schema.ErrorMessage
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/users [post]
func (self *Server) CreateUser(ctx *gin.Context) {
	companyId := ctx.Param("company")

	var data schema.NewUser
	if err := ctx.BindJSON(&data); err != nil {
		log.Error("Could not parse data from JSON")
		ctx.JSON(400, schema.ErrorMessage{
			Status:  400,
			Message: "Could not parse data from JSON",
		})
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
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

	log.Debug("Retrieving a user service worker...")
	conn, err := self.Collection.UserSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a user service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to user service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewUserClient(conn)
	msg := data.ToMessage(companyId)
	response, err := client.Store(ctx, &msg)
	if err != nil {
		log.Error("Error while creating user: %v", err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.User{}.FromMessage(response)
	if err != nil {
		log.Error("Could not parse created user: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status: 500,
			Message: "Could not parse created user",
		})
		return
	}

	ctx.JSON(201, res)
}

// @Summary Update user
// @Description Update information of a specific user
// @Tags User
// @Accept json
// @Param     company    path    string    true    "company UUID"
// @Param     id         path    string    true    "user UUID"
// @Param     data       body    schema.UpdatedUser    true    "user update data"
// @Success   200     {object}    model.User
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/users/{id} [put]
func (self *Server) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	companyId := ctx.Param("company")

	var data schema.UpdatedUser
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
			Status:  400,
			Message: "Error while validating data",
			Details: errors.Error(),
		})
		return
	}

	if (len(data.Password) > 0) && (len(data.Password) < 8) {
		log.Error("Error while validating data: Password has less than 8 characters")
		ctx.JSON(400, schema.ErrorMessage{
			Status: 400,
			Message: "Error while validating data",
			Details: "Password has less than 8 characters",
		})
		return
	}

	log.Debug("Retrieving a user service worker...")
	conn, err := self.Collection.UserSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a user service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to user service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewUserClient(conn)
	msg := data.ToMessage(companyId, id)
	response, err := client.Update(ctx, &msg)
	if err != nil {
		log.Error("Failed to update user: %v", err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	res, err := model.User{}.FromMessage(response)
	if err != nil {
		log.Error("Could not parse updated user: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not parse updated user",
		})
		return
	}

	ctx.JSON(200, res)
}

// @Summary Delete user
// @Description Delete a specific user
// @Tags      User
// @Accept    json
// @Param     company    path    string    true    "company UUID"
// @Param     id    path    string    true    "user UUID"
// @Success   200   "deleted successfully"
// @Failure   404     {object}    schema.ErrorMessage
// @Failure   500     {object}    schema.ErrorMessage
// @Router    /{company}/users/{id} [delete]
func (self *Server) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	companyId := ctx.Param("company")

	log.Debug("Retrieving a user service worker...")
	conn, err := self.Collection.UserSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a user service worker: %v", err)
		ctx.JSON(500, schema.ErrorMessage{
			Status:  500,
			Message: "Could not connect to user service",
		})
		return
	}
	defer conn.Close()

	client := rpc.NewUserClient(conn)
	_, err = client.Delete(ctx, &rpc.TenantEntityIndex{CompanyId: companyId, Index: id})
	if err != nil {
		log.Error("Failed to delete user %s: %v", id, err)
		m := schema.ErrorMessage{}.FromGrpcError(err)
		ctx.JSON(m.Status, m)
		return
	}

	ctx.Status(200)
}
