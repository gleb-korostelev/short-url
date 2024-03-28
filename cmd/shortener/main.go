package main

import (
	"net/http"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db/dbimpl"
	"github.com/gleb-korostelev/short-url.git/internal/service/handler"
	"github.com/gleb-korostelev/short-url.git/internal/service/router"
	"github.com/gleb-korostelev/short-url.git/internal/storage"
	"github.com/gleb-korostelev/short-url.git/internal/storage/filecache"
	"github.com/gleb-korostelev/short-url.git/internal/storage/inmemory"
	"github.com/gleb-korostelev/short-url.git/internal/storage/repository"
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
	svc := handler.NewAPIService(store)

	r := router.RouterInit(svc, log)

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
		store := inmemory.NewMemoryStorage(cache.Cache, &cache.Mu)
		logger.Infof("Using inmemory storage")
		return store, nil
	}
}
