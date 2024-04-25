package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	sloggin "github.com/samber/slog-gin"
	"github.com/gin-contrib/cors"

	log "github.com/Minerva-System/minerva-go/pkg/log"
	view "github.com/Minerva-System/minerva-go/internal/htmx/view"
	config "github.com/Minerva-System/minerva-go/internal/htmx/config"
)

func main() {
	godotenv.Load()
	log.Init()
	gin.SetMode(gin.ReleaseMode)

	log.Info("Minerva System: HTMX Server, Copyright (c) 2022-2024 Lucas S. Vieira")

	config.Load()
	
	router := gin.New()
	router.Use(sloggin.New(slog.Default()))
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.LoadHTMLGlob("templates/*/*.html")
	view.InstallRoutes(router)

	router.Run(config.Values.ServerHost)
}
