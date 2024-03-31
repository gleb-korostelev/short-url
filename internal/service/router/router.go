package router

import (
	"github.com/gleb-korostelev/short-url.git/internal/middleware"
	"github.com/gleb-korostelev/short-url.git/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RouterInit(svc service.APIServiceI, logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.GzipCompressMiddleware)
	router.Use(middleware.GzipDecompressMiddleware)
	router.Use(middleware.EnsureUserCookie)
	router.Get("/ping", middleware.LoggingMiddleware(svc.Ping, logger))
	router.Get("/{id}", middleware.LoggingMiddleware(svc.GetOriginal, logger))
	router.Get("/api/user/urls", middleware.LoggingMiddleware(svc.GetUserURLs, logger))
	router.Post("/", middleware.LoggingMiddleware(svc.PostShorter, logger))
	router.Post("/api/shorten", middleware.LoggingMiddleware(svc.PostShorterJSON, logger))
	router.Post("/api/shorten/batch", middleware.LoggingMiddleware(svc.ShortenBatchHandler, logger))
	router.Delete("/api/user/urls", middleware.LoggingMiddleware(svc.DeleteURLsHandler, logger))

	return router
}
