// Package handler contains HTTP handlers that provide web API functionality.
// These handlers manage operations such as URL creation, deletion, and redirection.
package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url/internal/config"
	"github.com/gleb-korostelev/short-url/internal/worker"
	"github.com/gleb-korostelev/short-url/tools/logger"
)

// DeleteURLsHandler handles the HTTP DELETE request for deleting one or more URLs.
// It requires the user to be authenticated and provides the functionality to mark URLs as deleted.
// This handler responds with HTTP status 202 (Accepted) to indicate that the delete request has been queued.
func (svc *APIService) DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the context, and return HTTP 401 Unauthorized if it's missing.
	userID, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Decode the request body to get a list of short URLs to be deleted.
	var shortURLs []string
	if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add the task to delete the URLs to the worker pool.
	svc.worker.AddTask(worker.Task{
		Action: func(ctx context.Context) error {
			err := svc.store.MarkURLsAsDeleted(ctx, userID, shortURLs)
			if err != nil {
				// Log the internal server error.
				logger.Errorf("Internal server error %v", err)
				return err
			}
			return nil
		},
	})

	// Respond with HTTP 202 Accepted to indicate the deletion task has been queued.
	w.WriteHeader(http.StatusAccepted)
}
