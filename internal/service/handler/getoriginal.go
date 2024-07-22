// Package handler contains HTTP handlers that provide web API functionality.
// These handlers manage operations such as URL creation, deletion, and redirection.
package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/go-chi/chi/v5"
)

// GetOriginal handles HTTP requests to retrieve the original URL based on a shortened URL identifier.
// The shortened URL ID is expected as a URL parameter.
//
// If the ID is not provided or the shortened URL cannot be found, it responds with HTTP 400 Bad Request.
// If the shortened URL has been marked as deleted, it responds with HTTP 410 Gone.
// Upon successful retrieval of the original URL, it sets the HTTP Location header with the original URL
// and responds with HTTP 307 Temporary Redirect.
func (svc *APIService) GetOriginal(w http.ResponseWriter, r *http.Request) {
	// Extract the 'id' URL parameter using the chi router.
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the store using the provided ID.
	originalURL, err := svc.store.GetOriginalLink(context.Background(), id)
	if err != nil {
		// Handle specific known errors, such as when the URL has been marked as deleted.
		if errors.Is(err, config.ErrGone) {
			http.Error(w, config.ErrGone.Error(), http.StatusGone)
			return
		}
		// Respond with Bad Request if the URL cannot be found or other errors occur.
		http.Error(w, "This URL doesn't exist", http.StatusBadRequest)
		return
	}

	// Set the Location header with the retrieved original URL.
	w.Header().Set("Location", string(originalURL))
	// Respond with Temporary Redirect to the original URL.
	w.WriteHeader(http.StatusTemporaryRedirect)
}
