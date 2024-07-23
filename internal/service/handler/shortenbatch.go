package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/models"
)

// ShortenBatchHandler processes HTTP POST requests to shorten multiple URLs simultaneously.
// It reads a list of URLs from the request body in JSON format and attempts to save each one,
// returning a corresponding list of shortened URLs or error messages.
//
// The function checks for the POST method and expects the user to be authenticated.
// If the request body does not contain valid JSON, or the batch is empty, it responds with
// HTTP 400 Bad Request. Each URL saving operation is performed, and the results are accumulated
// and returned as JSON with HTTP 201 Created on success.
//
// If an error occurs during the saving of any URL, it stops processing further and returns the results
// obtained until the error occurred.
func (svc *APIService) ShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the user is authenticated.
	userID, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Decode the JSON body to get a list of URLs to shorten.
	var reqItems []models.ShortenBatchRequestItem
	if err := json.NewDecoder(r.Body).Decode(&reqItems); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Set response content type to JSON.
	w.Header().Set("Content-Type", "application/json")

	// Reject empty batch requests.
	if len(reqItems) == 0 {
		http.Error(w, "Empty batch is not allowed", http.StatusBadRequest)
		return
	}

	// Process each URL in the batch and collect the results.
	var respItems []models.ShortenBatchResponseItem
	for _, item := range reqItems {
		shortURL, err := svc.store.SaveURL(context.Background(), item.OriginalURL, userID)
		if err != nil {
			// Return the results obtained until the error occurred.
			json.NewEncoder(w).Encode(respItems)
			http.Error(w, "Error with saving", http.StatusBadRequest)
			break
		}
		respItems = append(respItems, models.ShortenBatchResponseItem{
			CorrelationID: item.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	// Successfully encode and return the full list of shortened URLs.
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respItems)
}
