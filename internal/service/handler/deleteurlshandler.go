package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/service/business"
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

	// svc.worker.AddTask(worker.Task{
	// 	Action: func(ctx context.Context) error {
	// 		var shortURLs []string
	// 		if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
	// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 			return nil
	// 		}
	// 		userID, err := business.GetUserIDFromCookie(r)
	// 		if err != nil {
	// 			w.WriteHeader(http.StatusUnauthorized)
	// 			return nil
	// 		}
	// 		err = svc.store.MarkURLsAsDeleted(context.Background(), userID, shortURLs)
	// 		if err != nil {
	// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 		}
	// 		w.WriteHeader(http.StatusAccepted)
	// 		return nil
	// 	},
	// })
	// err = svc.store.MarkURLsAsDeleted(context.Background(), userID, shortURLs)
	// if err != nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// }
	// w.WriteHeader(http.StatusAccepted)

	var wg sync.WaitGroup

	sem := make(chan struct{}, config.MaxConcurrentUpdates)

	wg.Add(1)
	go func() {
		sem <- struct{}{}
		defer wg.Done()
		defer func() { <-sem }()
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
		err = svc.store.MarkURLsAsDeleted(context.Background(), userID, shortURLs)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusAccepted)

	}()
	// w.WriteHeader(http.StatusAccepted)

	wg.Wait()
}
