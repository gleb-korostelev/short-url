package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/service/business"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

func (svc *APIService) DeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	// var shortURLs []string
	// if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	// userID, err := business.GetUserIDFromCookie(r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	var shortURLs []string
	if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, err := business.GetUserIDFromCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	doneChan := make(chan struct{})
	svc.worker.AddTask(worker.Task{
		Action: func(ctx context.Context) error {
			err = svc.store.MarkURLsAsDeleted(context.Background(), userID, shortURLs)
			if err != nil {
				logger.Errorf("Internal server error %v", err)
				// http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusAccepted)
			return nil
		},
		Done: doneChan,
	})
	<-doneChan
	// err = svc.store.MarkURLsAsDeleted(context.Background(), userID, shortURLs)
	// if err != nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// }
	// w.WriteHeader(http.StatusAccepted)

	// var wg sync.WaitGroup

	// sem := make(chan struct{}, config.MaxConcurrentUpdates)

	// wg.Add(1)
	// go func() {
	// 	sem <- struct{}{}
	// 	defer wg.Done()
	// 	defer func() { <-sem }()
	// 	var shortURLs []string
	// 	if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
	// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 		return
	// 	}
	// 	userID, err := business.GetUserIDFromCookie(r)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusUnauthorized)
	// 		return
	// 	}
	// 	err = svc.store.MarkURLsAsDeleted(context.Background(), userID, shortURLs)
	// 	if err != nil {
	// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	}
	// 	w.WriteHeader(http.StatusAccepted)

	// }()
	// // w.WriteHeader(http.StatusAccepted)

	// wg.Wait()
}
