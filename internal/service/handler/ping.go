package handler

import (
	"context"
	"fmt"
	"net/http"
)

func (svc *APIService) Ping(w http.ResponseWriter, r *http.Request) {
	status, err := svc.store.Ping(context.Background())
	if err != nil {
		w.WriteHeader(status)
		fmt.Fprintln(w, err)
		return
	}
	w.WriteHeader(status)
}
