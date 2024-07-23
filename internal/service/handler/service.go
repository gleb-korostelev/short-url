// Package handler contains HTTP handlers that provide web API functionality.
// These handlers manage operations such as URL creation, deletion, and redirection.
package handler

import (
	"github.com/gleb-korostelev/short-url/internal/service"
	"github.com/gleb-korostelev/short-url/internal/storage"
	"github.com/gleb-korostelev/short-url/internal/worker"
)

// APIService implements the service.APIServiceI interface and provides methods
// to interact with the URL storage and processing tasks. It abstracts the
// details of data manipulation and task scheduling away from the HTTP interface.
type APIService struct {
	store  storage.Storage      // store is the interface to the URL storage backend.
	worker *worker.DBWorkerPool // worker handles asynchronous tasks using a worker pool.
}

// NewAPIService creates a new instance of APIService with the provided storage
// and worker pool implementations. This setup allows for flexible dependency injection
// and easier testing by decoupling the service logic from specific storage and worker implementations.
//
// store: Provides access to the URL storage and manipulation functions.
// worker: Manages asynchronous execution of background tasks that shouldn't block the HTTP handlers.
func NewAPIService(store storage.Storage, worker *worker.DBWorkerPool) service.APIServiceI {
	return &APIService{
		store:  store,
		worker: worker,
	}
}
