// Package handler contains HTTP handlers that provide web API functionality.
// These handlers manage operations such as URL creation, deletion, and redirection.
package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url/internal/config"
	"github.com/gleb-korostelev/short-url/internal/models"
)

// PostShorterJSON handles HTTP POST requests to create shortened URLs using JSON data.
// This method requires that the request method be POST and the content type be JSON.
// It ensures user authentication, reads the URL from the JSON payload, and saves the shortened URL.
//
// The function responds with:
// - HTTP 400 Bad Request if the request method is not POST or if there's an error parsing the request body.
// - HTTP 401 Unauthorized if the user is not authenticated.
// - HTTP 201 or other appropriate HTTP status based on the result of the URL saving operation.
//
// If the operation is successful, it returns the shortened URL in a JSON structure.
func (svc *APIService) PostShorterJSON(w http.ResponseWriter, r *http.Request) {
	// Ensure the HTTP method is POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}

	// Set the content type of the response to application/json.
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the user ID from the context, and ensure authentication.
	userID, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Decode the JSON body to get the original URL.
	var payload models.URLPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Attempt to save the URL and obtain a shortened version.
	shortURL, status, err := svc.store.SaveUniqueURL(context.Background(), payload.URL, userID)
	if err != nil {
		http.Error(w, "Error with saving", status)
		return
	}
	w.WriteHeader(status)

	// Encode the shortened URL in a JSON response.
	response := models.ShortURLResponse{Result: shortURL}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
