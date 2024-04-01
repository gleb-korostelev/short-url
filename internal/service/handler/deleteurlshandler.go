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

	var wg sync.WaitGroup

	sem := make(chan struct{}, config.MaxConcurrentUpdates)

	for i := 0; i < config.MaxConcurrentUpdates; i++ {
		wg.Add(1)
		go func(userID string, shortURLs []string) {
			sem <- struct{}{}
			defer wg.Done()
			defer func() { <-sem }()
			err = svc.store.MarkURLsAsDeleted(context.Background(), userID, shortURLs)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusAccepted)
		}(userID, shortURLs)
		// w.WriteHeader(http.StatusAccepted)
	}

	wg.Wait()
}
