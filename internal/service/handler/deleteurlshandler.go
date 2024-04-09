package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

func (svc *APIService) DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var shortURLs []string
	if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	svc.worker.AddTask(worker.Task{
		Action: func(ctx context.Context) error {
			err := svc.store.MarkURLsAsDeleted(ctx, userID, shortURLs)
			if err != nil {
				logger.Errorf("Internal server error %v", err)
			}
			return nil
		},
	})
	w.WriteHeader(http.StatusAccepted)
}
