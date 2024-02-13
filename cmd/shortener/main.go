package main

import (
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/handler/router"
)

func main() {
	r := router.RouterInit()

	fmt.Println("Server is listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
