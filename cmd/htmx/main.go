package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sloggin "github.com/samber/slog-gin"
	"github.com/gin-contrib/cors"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	view "github.com/Minerva-System/minerva-go/internal/htmx/view"
)

func main() {
	godotenv.Load()
	log.Init()
	gin.SetMode(gin.ReleaseMode)

	log.Info("Minerva System: HTMX Server, Copyright (c) 2022-2024 Lucas S. Vieira")

	host := ":5090"
	
	router := gin.New()
	router.Use(sloggin.New(slog.Default()))
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	view.InstallRoutes(router)

	router.Run(host)
}
