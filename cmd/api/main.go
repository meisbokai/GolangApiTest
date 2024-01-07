package main

import (
	"log"

	"github.com/meisbokai/GolangApiTest/cmd/api/server"
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
