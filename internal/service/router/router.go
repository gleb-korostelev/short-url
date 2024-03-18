package router

import (
	"github.com/gleb-korostelev/short-url.git/internal/middleware"
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RouterInit(logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.GzipCompressMiddleware)
	router.Use(middleware.GzipDecompressMiddleware)
	router.Get("/{id}", middleware.LoggingMiddleware(handler.GetOriginal, logger))
	router.Post("/", middleware.LoggingMiddleware(handler.PostShorter, logger))
	router.Post("/api/shorten", middleware.LoggingMiddleware(handler.PostShorterJSON, logger))

	return router
}
