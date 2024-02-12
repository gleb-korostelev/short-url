package router

import (
	"github.com/gleb-korostelev/short-url.git/internal/handler/business"
	"github.com/go-chi/chi/v5"
)

func RouterInit() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/{id}", business.GetOriginal)
	router.Post("/", business.PostShorter)
	return router
}
