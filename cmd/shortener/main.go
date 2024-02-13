package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/handler/router"
)

func main() {
	r := router.RouterInit()

	flag.Parse()

	fmt.Printf("Server will run on: %s\n", config.ServerAddr)
	fmt.Printf("Base URL for shortened links: %s\n", config.BaseURL)

	fmt.Println("Server is listening on ", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, r); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
