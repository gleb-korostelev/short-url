// Package handler contains HTTP handlers that provide web API functionality.
// These handlers manage operations such as URL creation, deletion, and redirection.
package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

// PostShorter handles the HTTP POST requests for creating shortened URLs.
// It reads the original URL from the request body, validates the user's identity,
// and submits a task to asynchronously save the URL and generate a shortened version.
//
// The function ensures that the request uses the POST method. If not, it responds with HTTP 400 Bad Request.
// It requires user authentication, responding with HTTP 401 Unauthorized if the user ID is not found in the context.
// If the request body cannot be read, it responds with HTTP 400 Bad Request.
// The response includes the shortened URL on success or appropriate error messages.
func (svc *APIService) PostShorter(w http.ResponseWriter, r *http.Request) {
	// Validate the request method.
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}

	// Validate user identity from the context.
	userID, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Set the content type of the response.
	w.Header().Set("content-type", "text/plain")

	// Read the original URL from the request body.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(body)

	// Create a channel to wait for the asynchronous task to complete.
	doneChan := make(chan struct{})

	// Submit the task to the worker pool.
	svc.worker.AddTask(worker.Task{
		Action: func(ctx context.Context) error {
			shortURL, status, err := svc.store.SaveUniqueURL(ctx, originalURL, userID)
			w.WriteHeader(status)
			if err != nil {
				logger.Errorf("Error with saving data: %v", err)
				return nil
			}
			// Write the shortened URL to the response.
			fmt.Fprint(w, shortURL)
			return nil
		},
		Done: doneChan,
	})

	// Wait for the task to complete.
	<-doneChan
}
