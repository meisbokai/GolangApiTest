package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meisbokai/GolangApiTest/internal/http/routes"
)

type App struct {
	HttpServer *http.Server
}

func NewServerApp() (*App, error) {
	// Setup router
	router := setupRouter()

	// API Routes
	api := router.Group("api")
	api.GET("/", routes.RootHandler)

	// http Server
	server := &http.Server{
		Addr:           ":3000",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &App{
		HttpServer: server,
	}, nil
}

func (a *App) Run() (err error) {
	if err := a.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}
	return
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)

	// Create new router
	router := gin.New()

	// Set middleware
	router.Use(gin.Recovery())

	return router
}
