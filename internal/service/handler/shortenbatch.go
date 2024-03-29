package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/models"
)

func (svc *APIService) ShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	var reqItems []models.ShortenBatchRequestItem
	if err := json.NewDecoder(r.Body).Decode(&reqItems); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	if len(reqItems) == 0 {
		http.Error(w, "Empty batch is not allowed", http.StatusBadRequest)
		return
	}
	var respItems []models.ShortenBatchResponseItem
	for _, item := range reqItems {
		shortURL, err := svc.store.SaveURL(context.Background(), item.OriginalURL)
		if err != nil {
			json.NewEncoder(w).Encode(respItems)
			http.Error(w, "Error with saving", http.StatusBadRequest)
			break
		}
		respItems = append(respItems, models.ShortenBatchResponseItem{
			CorrelationID: item.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respItems)
}
