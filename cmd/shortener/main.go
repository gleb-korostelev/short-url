package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db/dbimpl"
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/gleb-korostelev/short-url.git/internal/service/router"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/internal/storage/filecache"
	"github.com/gleb-korostelev/short-url.git/internal/storage/inmemory"
	"github.com/gleb-korostelev/short-url.git/internal/storage/repository"
	"github.com/gleb-korostelev/short-url.git/internal/worker"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"go.uber.org/zap"
)

func main() {

	config.ConfigInit()
	log, _ := zap.NewProduction()

	store, err := storageInit()
	if err != nil {
		return
	}
	defer store.Close()

	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	defer workerPool.Shutdown()
	svc := handler.NewAPIService(store, workerPool)

	r := router.RouterInit(svc, log)

	// go func() {
	// 	logger.Infof("Starting pprof server on :6060")
	// 	if err := http.ListenAndServe(":6060", nil); err != nil {
	// 		logger.Fatal("pprof server failed", zap.Error(err))
	// 	}
	// }()

	logger.Infof("Base URL for shortened links: %s", config.BaseURL)

	logger.Infof("Server is listening on: %s", config.ServerAddr)
	if err := http.ListenAndServe(config.ServerAddr, r); err != nil {
		logger.Fatal("Error starting server: %v", zap.Error(err))
	}

}

func storageInit() (storage.Storage, error) {
	if config.DBDSN != "" {
		database, err := dbimpl.InitDB()
		if err != nil {
			return nil, err
		}
		store := repository.NewDBStorage(database)
		logger.Infof("Using database storage")
		return store, nil
	} else if config.BaseFilePath != "" {
		store := filecache.NewFileStorage(config.BaseFilePath)
		logger.Infof("Using file storage with base file path %s", config.BaseFilePath)
		return store, nil
	} else {
		store := inmemory.NewMemoryStorage(cache.Cache)
		logger.Infof("Using inmemory storage")
		return store, nil
	}
}
