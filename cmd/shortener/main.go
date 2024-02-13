package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/handler/router"
	"github.com/gleb-korostelev/short-url.git/internal/utils"
)

func main() {
	r := router.RouterInit()

	flag.Parse()

	config.ServerAddr = utils.GetEnv("SERVER_ADDRESS", config.ServerAddr)
	config.BaseURL = utils.GetEnv("BASE_URL", config.BaseURL)
	
	fmt.Printf("Server will run on: %s\n", config.ServerAddr)
	fmt.Printf("Base URL for shortened links: %s\n", config.BaseURL)

	fmt.Println("Server is listening on ", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, r); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
