package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/service/utils"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

// GetUserURLs handles the HTTP GET request to retrieve all URLs associated with the authenticated user.
// The user's ID is extracted from the cookie, and the URLs are fetched from the storage system.
//
// If the user ID cannot be validated or is missing from the cookie, the handler responds with HTTP 401 Unauthorized.
// In case of any internal errors during URL retrieval, it responds with HTTP 500 Internal Server Error.
// If no URLs are associated with the user, it responds with HTTP 204 No Content.
// On successful data retrieval, it returns a list of URLs in JSON format with HTTP 200 OK.
func (svc *APIService) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the cookie; return HTTP 401 if extraction fails.
	userID, err := utils.GetUserIDFromCookie(r)
	if err != nil {
		logger.Errorf("Error retrieving user ID from cookie: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// logger.Infof("Retrieved userID is: %s", userID)

	// Fetch all URLs associated with the user ID from the storage.
	urls, err := svc.store.GetAllURLS(context.Background(), userID, config.BaseURL)
	if err != nil {
		logger.Errorf("Error retrieving URLs from store: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return HTTP 204 No Content if no URLs are found.
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Return the list of URLs in JSON format with HTTP 200 OK.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		logger.Errorf("Error encoding URLs to JSON: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
