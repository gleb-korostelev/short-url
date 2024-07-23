package handler

import (
	"context"
	"fmt"
	"net/http"
)

// Ping handles an HTTP GET request to check the database connectivity.
// It sends an HTTP status code based on the result of the ping to the database.
//
// The handler attempts to ping the database and returns the HTTP status code indicating
// the result of this operation. If the ping is unsuccessful, it returns the error from the database.
//
// Successful ping to the database results in HTTP 200 OK.
// If an error occurs during the ping, the HTTP status code corresponding to the error is returned,
// along with the error message.
func (svc *APIService) Ping(w http.ResponseWriter, r *http.Request) {
	// Attempt to ping the database.
	status, err := svc.store.Ping(context.Background())
	if err != nil {
		// If there is an error, write the status code associated with the error and print the error message.
		w.WriteHeader(status)
		fmt.Fprintln(w, err)
		return
	}
	// If the ping is successful, write HTTP 200 OK.
	w.WriteHeader(status)
}
