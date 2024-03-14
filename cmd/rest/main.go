package main

import (
	"strconv"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sloggin "github.com/samber/slog-gin"
	"github.com/gin-contrib/cors"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/Minerva-System/minerva-go/cmd/rest/docs"
	connection "github.com/Minerva-System/minerva-go/internal/connection"
	rpc "github.com/Minerva-System/minerva-go/internal/rpc"
	log "github.com/Minerva-System/minerva-go/pkg/log"
	model "github.com/Minerva-System/minerva-go/internal/model"
)

//go:generate swag init --parseDependency --parseInternal --parseDepth 1 -g main.go

// @title Minerva System API
// @version 1.0
// @description Minerva System API

// @contact.name Lucas S. Vieira
// @contact.url https://luksamuk.codes
// @contact.email lucasvieira@protonmail.com

// TODO: license.name, license.url

// @BasePath /api/v1
// @query.collection.format multi

type Server struct {
	col connection.Collection
}

type ErrorMessage struct {
	Status int `json:"status"`
	Message string `json:"message"`
}

// @Summary User list
// @Description Get a list of users per page
// @Tags      User
// @Accept    json
// @Produce   json
// @Param     page    query    int    false    "page number (0 or more)"
// @Success   200     {object}    []model.User
// @Failure   400     {object}    ErrorMessage
// @Failure   500     {object}    ErrorMessage
// @Router    /user [get]
func (self *Server) getUsers(ctx *gin.Context) {
	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "0"), 10, 64)
	if err != nil || page < 0 {
		log.Error("Could not parse page size")
		ctx.JSON(400, ErrorMessage{Status: 400, Message: "Could not parse page index"})
		return
	}
	
	log.Debug("Retrieving a user service worker...")
	conn, err := self.col.UserSvc.Get(ctx)
	if err != nil {
		log.Error("Failed to retrieve a user service worker: %v", err)
		ctx.JSON(500, ErrorMessage{Status: 500, Message: "Could not connect to user service"})
		return
	}
	defer conn.Close() // Very important!

	client := rpc.NewUserClient(conn)
	response, err := client.Index(ctx, &rpc.PageIndex{ Index: &page })
	if err != nil {
		log.Error("Failed to retrieve user index: %v", err)
		ctx.JSON(500, ErrorMessage{Status: 500, Message: "Could not connect to user service"})
		return
	}

	res, err := model.User{}.FromListMessage(response)
	if err != nil {
		log.Error("Could not parse retrieved user list: %v", err)
		ctx.JSON(500, ErrorMessage{Status: 500, Message: "Could not parse retrieved user list"})
	}
	
	ctx.JSON(200, res)
}

func main() {
	godotenv.Load()
	log.Init()
	gin.SetMode(gin.ReleaseMode)
	
	log.Info("Minerva System: REST gateway service (Go port)")
	log.Info("Copyright (c) 2022-2023 Lucas S. Vieira")

	log.Info("Establishing connections...")

	// Test
	col, err := connection.NewCollection(connection.CollectionOptions{
		WithUserService: true,
		WithSessionService: true,
		WithProductsService: true,
	})
	if err != nil {
		log.Fatal("Failed to establish connections: %v", err)
	}

	server := Server{
		col: col,
	}
	
	host := ":9000"

	router := gin.New()
	router.Use(sloggin.New(slog.Default()))
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


	// Example route for managing gRPC connections
	v1 := router.Group("/api/v1")
	
	v1.GET("/user", server.getUsers)
	
	router.Run(host)
}
