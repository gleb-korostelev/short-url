package handler

import (
	"github.com/gleb-korostelev/short-url.git/internal/service"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
)

type APIService struct {
	store  storage.Storage
	worker *worker.DBWorkerPool
}

func NewAPIService(store storage.Storage, worker *worker.DBWorkerPool) service.APIServiceI {
	return &APIService{
		store:  store,
		worker: worker,
	}
}
