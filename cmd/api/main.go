package main

import (
	"log"

	"github.com/meisbokai/GolangApiTest/cmd/api/server"
	_ "github.com/meisbokai/GolangApiTest/docs"

	config "github.com/meisbokai/GolangApiTest/internal/configs"
	"github.com/meisbokai/GolangApiTest/internal/constants"
	"github.com/meisbokai/GolangApiTest/pkg/logger"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		logger.Fatal(err.Error(), logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
	}
	logger.Info("configuration loaded", logrus.Fields{constants.LoggerCategory: constants.LoggerCategoryConfig})
}

// @title           User API Server Test
// @version         1.0
// @description     This is a mock server for an assignment
// @termsOfService  http://swagger.io/terms/

// @contact.name   Neow Bo Kai
// @contact.url    www.github.com/meisbokai
// @contact.email  neow.bokai@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api

// @securityDefinitions.apikey jwtToken
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	// Start API server
	app, err := server.NewServerApp()
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
