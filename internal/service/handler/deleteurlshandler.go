package handler

import (
	"context"
	"encoding/json"
	"net/http"

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
	err = svc.store.MarkURLsAsDeleted(context.Background(), userID, shortURLs)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusAccepted)
}
