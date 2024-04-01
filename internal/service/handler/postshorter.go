package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
)

func (svc *APIService) PostShorter(w http.ResponseWriter, r *http.Request) {

	// svc.worker.AddTask(worker.Task{
	// 	Action: func(ctx context.Context) error {
	// 		if r.Method != http.MethodPost {
	// 			http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
	// 			return nil
	// 		}
	// 		userID, ok := r.Context().Value(config.UserContextKey).(string)
	// 		if !ok {
	// 			w.WriteHeader(http.StatusUnauthorized)
	// 			return nil
	// 		}
	// 		w.Header().Set("content-type", "text/plain")
	// 		body, err := io.ReadAll(r.Body)
	// 		if err != nil {
	// 			http.Error(w, "Error reading request body", http.StatusBadRequest)
	// 			return nil
	// 		}
	// 		defer r.Body.Close()

	// 		originalURL := string(body)

	// 		shortURL, status, err := svc.store.SaveUniqueURL(context.Background(), originalURL, userID)
	// 		w.WriteHeader(status)
	// 		if err != nil {
	// 			logger.Errorf("Error with saving data in here %v", err)
	// 			http.Error(w, "Error with saving data", http.StatusBadRequest)
	// 			return nil
	// 		}
	// 		fmt.Fprint(w, shortURL)
	// 		return nil
	// 	},
	// })

	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
	// 	return
	// }
	// userID, ok := r.Context().Value(config.UserContextKey).(string)
	// if !ok {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// w.Header().Set("content-type", "text/plain")
	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Error reading request body", http.StatusBadRequest)
	// 	return
	// }
	// defer r.Body.Close()

	// originalURL := string(body)

	// shortURL, status, err := svc.store.SaveUniqueURL(context.Background(), originalURL, userID)
	// w.WriteHeader(status)
	// if err != nil {
	// 	logger.Errorf("Error with saving data in here %v", err)
	// 	http.Error(w, "Error with saving data", http.StatusBadRequest)
	// 	return
	// }
	// fmt.Fprint(w, shortURL)

	var wg sync.WaitGroup

	sem := make(chan struct{}, config.MaxConcurrentUpdates)

	wg.Add(1)
	go func() {
		sem <- struct{}{}
		defer wg.Done()
		defer func() { <-sem }()
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
		shortURL, status, err := svc.store.SaveUniqueURL(context.Background(), originalURL, userID)
		w.WriteHeader(status)
		if err != nil {
			logger.Errorf("Error with saving data in here %v", err)
			http.Error(w, "Error with saving data", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, shortURL)
	}()
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
