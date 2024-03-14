package minerva_api_v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	
	log "github.com/Minerva-System/minerva-go/pkg/log"
	model "github.com/Minerva-System/minerva-go/internal/model"
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
)

// @Summary User list
// @Description Get a list of users per page
// @Tags      User
// @Accept    json
// @Produce   json
// @Param     page    query    int    false    "page number (0 or more)"
// @Success   200     {object}    []model.User
// @Failure   400     {object}    model.ErrorMessage
// @Failure   500     {object}    model.ErrorMessage
// @Router    /user [get]
func (self *Server) GetUsers(ctx *gin.Context) {
	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "0"), 10, 64)
	if err != nil || page < 0 {
		log.Error("Could not parse page size")
		ctx.JSON(400, model.ErrorMessage{
			Status: 400,
			Message: "Could not parse page index",
		})
		return
	}
	
	log.Debug("Retrieving a user service worker...")
	conn, err := self.Collection.UserSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a user service worker: %v", err)
		ctx.JSON(500, model.ErrorMessage{
			Status: 500,
			Message: "Could not connect to user service",
		})
		return
	}
	defer conn.Close() // Very important!

	client := rpc.NewUserClient(conn)
	response, err := client.Index(ctx, &rpc.PageIndex{ Index: &page })
	if err != nil {
		log.Error("Failed to retrieve user index: %v", err)
		ctx.JSON(500, model.ErrorMessage{
			Status: 500,
			Message: "Could not connect to user service",
		})
		return
	}

	res, err := model.User{}.FromListMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved user list: %v", err)
		ctx.JSON(500, model.ErrorMessage{
			Status: 500,
			Message: "Could not parse retrieved user list",
		})
	}
	
	ctx.JSON(200, res)
}

