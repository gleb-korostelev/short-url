package handler

import (
	"context"
	"fmt"
	"net/http"
)

func (svc *APIService) Ping(w http.ResponseWriter, r *http.Request) {
	err := svc.data.Ping(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to connect to the database")
		return
	}

	w.WriteHeader(http.StatusOK)
}
