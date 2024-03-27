package main

import (
	"net/http"
	"os"

	"github.com/gleb-korostelev/short-url.git/internal/cache"
	"github.com/gleb-korostelev/short-url.git/internal/config"
	"github.com/gleb-korostelev/short-url.git/internal/db/impl"
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
		database, err := impl.InitDB()
		if err != nil {
			return nil, err
		}
		// defer database.Close()
		store := repository.NewDBStorage(database)
		logger.Info("Using database storage")
		return store, nil
	} else if _, err := os.Stat(config.BaseFilePath); err == nil {
		store := filecache.NewFileStorage(config.BaseFilePath)
		logger.Info("Using file storage")
		return store, nil
	} else {
		store := inmemory.NewMemoryStorage(cache.Cache, &cache.Mu)
		logger.Info("Using inmemory storage")
		return store, nil
	}
}
