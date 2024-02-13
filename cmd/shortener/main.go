package main

import (
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/handler"
)

func main() {
	http.HandleFunc(`/`, handler.HandleRequest)
	fmt.Println("Server is listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
