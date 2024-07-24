// Package service defines interfaces that specify the operations for managing URLs via HTTP.
// This package facilitates the decoupling of HTTP handling logic from actual service implementation,
// promoting a more modular architecture.
package service

import "net/http"

// APIServiceI defines the interface for handling HTTP requests related to URL operations.
// Implementations of this interface are responsible for handling all web interactions,
// including URL creation, retrieval, and deletion.
type APIServiceI interface {
	// GetOriginal retrieves the original URL corresponding to a given shortened URL ID.
	// It writes the result to the HTTP response.
	GetOriginal(w http.ResponseWriter, r *http.Request)

	// Ping checks the connectivity and readiness of the underlying services or databases.
	// It is typically used for health checks and monitoring.
	Ping(w http.ResponseWriter, r *http.Request)

	// PostShorter handles the creation of a shortened URL from a plain text input received in the HTTP request body.
	// It writes the shortened URL or an error message to the HTTP response.
	PostShorter(w http.ResponseWriter, r *http.Request)

	// PostShorterJSON handles the creation of a shortened URL from JSON input received in the HTTP request body.
	// It writes the shortened URL or an error message in JSON format to the HTTP response.
	PostShorterJSON(w http.ResponseWriter, r *http.Request)

	// ShortenBatchHandler handles requests for creating multiple shortened URLs in a batch from a JSON array input.
	// It writes the results or an error message in JSON format to the HTTP response.
	ShortenBatchHandler(w http.ResponseWriter, r *http.Request)

	// GetUserURLs retrieves all URLs associated with the authenticated user.
	// It writes the list of URLs or an error message in JSON format to the HTTP response.
	GetUserURLs(w http.ResponseWriter, r *http.Request)

	// DeleteURLsHandler handles the deletion of one or more URLs specified in the request body.
	// It writes the result status to the HTTP response.
	DeleteURLsHandler(w http.ResponseWriter, r *http.Request)
}
