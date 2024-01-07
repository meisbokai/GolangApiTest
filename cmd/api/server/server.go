package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	config "github.com/meisbokai/GolangApiTest/internal/configs"
	"github.com/meisbokai/GolangApiTest/internal/constants"
	"github.com/meisbokai/GolangApiTest/internal/http/middlewares"
	"github.com/meisbokai/GolangApiTest/internal/http/routes"
	"github.com/meisbokai/GolangApiTest/internal/util"
	"github.com/meisbokai/GolangApiTest/pkg/jwt"
	"github.com/meisbokai/GolangApiTest/pkg/logger"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	HttpServer *http.Server
}

func NewServerApp() (*App, error) {

	// Setup database
	conn, err := util.SetupPostgresConnection()
	if err != nil {
		return nil, err
	}

	// Setup router
	router := setupRouter()

	// JWT Middleware
	jwtService := jwt.NewJWTService(config.AppConfig.JWTSecret, config.AppConfig.JWTIssuer, config.AppConfig.JWTExpired)

	// JWT for all users
	authMiddleware := middlewares.NewAuthMiddleware(jwtService, false)

	// JWT for Admin user only
	adminAuthMiddleware := middlewares.NewAuthMiddleware(jwtService, true)

	// API Routes
	api := router.Group("api")
	api.GET("/", routes.RootHandler)
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.NewUsersRoute(api, conn, jwtService, authMiddleware).Routes()
	routes.NewAdminRoute(api, conn, jwtService, adminAuthMiddleware).Routes()

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
	// Gracefull Shutdown
	go func() {
		logger.InfoF("success to listen and serve on :%d", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer}, config.AppConfig.Port)
		if err := a.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// make blocking channel and waiting for a signal
	<-quit
	logger.Info("shutdown server ...", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.HttpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error when shutdown server: %v", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	logger.Info("timeout of 5 seconds.", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
	logger.Info("server exiting", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryServer})
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
