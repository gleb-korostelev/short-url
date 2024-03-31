package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/service/business"
)

func (svc *APIService) PostShorter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusBadRequest)
		return
	}
	userID, err := business.GetUserIDFromCookie(r)
	if err != nil {
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
		http.Error(w, "Error with saving file", http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, shortURL)
}
