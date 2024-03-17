package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/db"
)

func Ping(w http.ResponseWriter, r *http.Request, db *db.Database) {
	err := db.Conn.Ping(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Failed to connect to the database")
		return
	}

	w.WriteHeader(http.StatusOK)
}
