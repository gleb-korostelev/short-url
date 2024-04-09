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
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.EnsureUserCookie)
	router.Get("/ping", svc.Ping)
	router.Get("/{id}", svc.GetOriginal)
	router.Get("/api/user/urls", svc.GetUserURLs)
	router.Post("/", svc.PostShorter)
	router.Post("/api/shorten", svc.PostShorterJSON)
	router.Post("/api/shorten/batch", svc.ShortenBatchHandler)
	router.Delete("/api/user/urls", svc.DeleteURLsHandler)

	return router
}
