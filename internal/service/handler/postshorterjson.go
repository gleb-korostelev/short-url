package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/models"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
)

func PostShorterJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}

	var payload models.URLPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := business.CacheURL(payload.URL)
	if err != nil {
		http.Error(w, "Error with saving file", http.StatusBadRequest)
		return
	}

	response := models.ShortURLResponse{Result: shortURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
