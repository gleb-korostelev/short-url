package business

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/utils"
)

func PostShorter(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(body)

	shortUrl := utils.CacheURL(originalURL)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, shortUrl)
}
