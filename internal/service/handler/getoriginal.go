package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (svc *APIService) GetOriginal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}
	// cache.Mu.RLock()
	// originalURL, exists := business.GetOriginalURL(svc.data, id)
	// cache.Mu.RUnlock()

	originalURL, err := svc.store.GetOriginalLink(context.Background(), id)

	if err != nil {
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", string(originalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
