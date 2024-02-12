package business

import (
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
)

func GetOriginal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get method is allowed", http.StatusBadRequest)
		return
	}

	id := r.URL.Path[1:]

	config.Mu.Lock()
	originalURL, exists := config.Cache[id]
	config.Mu.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	// http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect) // Использование кода 307
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("content-type", "application/json")
	fmt.Fprint(w, string(originalURL))
}
