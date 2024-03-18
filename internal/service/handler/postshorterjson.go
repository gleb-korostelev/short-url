package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
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
	shortURL, status, err := business.CacheURL(payload.URL, svc.data)
	if err != nil {
		http.Error(w, "Error with saving", http.StatusBadRequest)
		return
	}

	response := models.ShortURLResponse{Result: shortURL}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
