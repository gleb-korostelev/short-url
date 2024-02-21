package router

import (
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/go-chi/chi/v5"
)

func RouterInit() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{id}", handler.GetOriginal)
	router.Post("/", handler.PostShorter)

	return router
}
