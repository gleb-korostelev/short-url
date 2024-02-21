package handler

import (
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/go-chi/chi/v5"
)

func GetOriginal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}
	business.Mu.RLock()
	originalURL, exists := business.Cache[id]
	business.Mu.RUnlock()

	if !exists {
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", string(originalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
