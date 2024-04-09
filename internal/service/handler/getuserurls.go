package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
)

func (svc *APIService) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	urls, err := svc.store.GetAllURLS(context.Background(), userID, config.BaseURL)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}
