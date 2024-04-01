package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gleb-korostelev/short-url.git/internal/config"
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

	var wg sync.WaitGroup

	sem := make(chan struct{}, config.MaxConcurrentUpdates)

	wg.Add(1)
	go func(userID string, originalURL string) {
		sem <- struct{}{}
		defer wg.Done()
		defer func() { <-sem }()
		shortURL, status, err := svc.store.SaveUniqueURL(context.Background(), originalURL, userID)
		w.WriteHeader(status)
		if err != nil {
			http.Error(w, "Error with saving file", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, shortURL)
	}(userID, originalURL)
	// w.WriteHeader(http.StatusAccepted)

	wg.Wait()

	// shortURL, status, err := svc.store.SaveUniqueURL(context.Background(), originalURL, userID)
	// w.WriteHeader(status)
	// if err != nil {
	// 	http.Error(w, "Error with saving file", http.StatusBadRequest)
	// 	return
	// }
	// fmt.Fprint(w, shortURL)
}
