package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/models"
)

func (svc *APIService) PostShorterJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var payload models.URLPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// shortURL, status, err := business.CacheURL(payload.URL, svc.data)

	shortURL, status, err := svc.store.SaveUniqueURL(context.Background(), payload.URL)
	w.WriteHeader(status)
	if err != nil {
		http.Error(w, "Error with saving", http.StatusBadRequest)
		return
	}

	response := models.ShortURLResponse{Result: shortURL}
	json.NewEncoder(w).Encode(response)
}
