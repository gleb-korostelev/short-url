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
	var shortURLs []string
	if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(config.UserContextKey).(string)

	doneChan := make(chan struct{})
	svc.worker.AddTask(worker.Task{
		Action: func(ctx context.Context) error {
			err := svc.store.MarkURLsAsDeleted(ctx, userID, shortURLs)
			if err != nil {
				logger.Errorf("Internal server error %v", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			// w.WriteHeader(http.StatusAccepted)
			return nil
		},
		Done: doneChan,
	})
	w.WriteHeader(http.StatusAccepted)
	<-doneChan
}
