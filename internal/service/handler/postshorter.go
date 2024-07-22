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

func (svc *APIService) PostShorter(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(config.UserContextKey).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("content-type", "text/plain")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	originalURL := string(body)

	doneChan := make(chan struct{})

	svc.worker.AddTask(worker.Task{
		Action: func(ctx context.Context) error {
			shortURL, status, err := svc.store.SaveUniqueURL(ctx, originalURL, userID)
			w.WriteHeader(status)
			if err != nil {
				logger.Errorf("Error with saving data in here %v", err)
				return nil
			}
			fmt.Fprint(w, shortURL)
			return nil
		},
		Done: doneChan,
	})

	<-doneChan
}
