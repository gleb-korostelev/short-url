// Package router sets up the HTTP routing for the web service, defining paths and associating them
// with handler functions and middleware.
package router

import (
	"github.com/gleb-korostelev/short-url.git/internal/middleware"
	"github.com/gleb-korostelev/short-url.git/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// RouterInit initializes the web server's routes and configures middleware. It takes a service
// interface and a logger as parameters, setting up routes that handle URL shortening operations
// and other related tasks.
//
// Parameters:
//
//	svc:    A service interface that provides methods for handling various HTTP requests related to URL management.
//	logger: A logger from the zap library used for logging within middleware.
//
// Returns:
//
//	A *chi.Mux router configured with all routes and middleware for the application.
//
// The function sets up the following routes:
//   - GET /ping: Checks database connectivity.
//   - GET /{id}: Retrieves the original URL corresponding to a shortened ID.
//   - GET /api/user/urls: Retrieves all URLs associated with the authenticated user.
//   - POST /: Creates a shortened URL from a plain text body.
//   - POST /api/shorten: Creates a shortened URL from JSON input.
//   - POST /api/shorten/batch: Handles batch creation of shortened URLs.
//   - DELETE /api/user/urls: Deletes one or more URLs associated with the user.
//
// Middleware used:
//   - GzipCompressMiddleware: Compresses response data if the client supports gzip.
//   - GzipDecompressMiddleware: Decompresses request data if compressed with gzip.
//   - LoggingMiddleware: Logs the details of the HTTP request and response.
//   - EnsureUserCookie: Ensures that a user session is valid or creates a new session.
func RouterInit(svc service.APIServiceI, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	// Register middleware that will be used across all routes.
	router.Use(middleware.GzipCompressMiddleware)
	router.Use(middleware.GzipDecompressMiddleware)
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.EnsureUserCookie)

	// Define routes and associate them with specific handler functions.
	router.Get("/ping", svc.Ping)
	router.Get("/{id}", svc.GetOriginal)
	router.Get("/api/user/urls", svc.GetUserURLs)
	router.Post("/", svc.PostShorter)
	router.Post("/api/shorten", svc.PostShorterJSON)
	router.Post("/api/shorten/batch", svc.ShortenBatchHandler)
	router.Delete("/api/user/urls", svc.DeleteURLsHandler)

	return router
}
