package main

import (
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/service/router"
)

func main() {
	config.ConfigInit()
	r := router.RouterInit()

	fmt.Printf("Server will run on: %s\n", config.ServerAddr)
	fmt.Printf("Base URL for shortened links: %s\n", config.BaseURL)

	fmt.Println("Server is listening on ", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, r); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
