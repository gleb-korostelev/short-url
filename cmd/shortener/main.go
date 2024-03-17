package main

import (
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/gleb-korostelev/short-url.git/internal/service/router"
	"go.uber.org/zap"
)

func main() {
	config.ConfigInit()
	business.LoadURLs()
	database := db.InitDB()
	defer database.Close()

	logger, _ := zap.NewProduction()
	r := router.RouterInit(database, logger)

	fmt.Printf("Base URL for shortened links: %s\n", config.BaseURL)

	fmt.Println("Server is listening on ", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, r); err != nil {
		logger.Fatal("Error starting server:", zap.Error(err))
	}
}
