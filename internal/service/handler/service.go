package handler

import (
	"github.com/gleb-korostelev/short-url.git/internal/service"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
)

type APIService struct {
	store storage.Storage
}

func NewAPIService(store storage.Storage) service.APIServiceI {
	return &APIService{
		store: store,
	}
}
