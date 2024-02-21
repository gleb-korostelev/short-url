package handler

import (
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/go-chi/chi/v5"
)

func GetOriginal(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}
	cache.Mu.RLock()
	originalURL, exists := cache.Cache[id]
	cache.Mu.RUnlock()
	if !exists {
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", string(originalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
