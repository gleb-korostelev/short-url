package main

import (
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db/impl"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/gleb-korostelev/short-url.git/internal/service/router"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"go.uber.org/zap"
)

func main() {
	config.ConfigInit()
	database, err := impl.InitDB()
	if err == nil {
		defer database.Close()
		err := impl.InitializeTables(database)
		if err != nil {
			logger.Infof("Failed to initialize tables: %v", err)
		}
	} else {
		logger.Infof("Database wasn't initialized: %v", err)
	}
	business.LoadURLs()
	svc := handler.NewAPIService(database)

	log, _ := zap.NewProduction()
	r := router.RouterInit(svc, log)

	logger.Infof("Base URL for shortened links: %s", config.BaseURL)

	logger.Infof("Server is listening on: %s", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, r); err != nil {
		logger.Fatal("Error starting server: %v", zap.Error(err))
	}
}
