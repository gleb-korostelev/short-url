package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/service/business"
)

func PostShorter(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(body)

	shortURL := business.CacheURL(originalURL)

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, shortURL)
}
