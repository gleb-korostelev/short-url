package main

import (
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/gleb-korostelev/short-url.git/internal/service/router"
	"go.uber.org/zap"
)

func main() {
	config.ConfigInit()
	logger, _ := zap.NewProduction()
	r := router.RouterInit(logger)

	business.LoadURLs()

	fmt.Printf("Server will run on: %s\n", config.ServerAddr)
	fmt.Printf("Base URL for shortened links: %s\n", config.BaseURL)

	fmt.Println("Server is listening on ", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, r); err != nil {
		logger.Fatal("Error starting server:", zap.Error(err))
	}
}
