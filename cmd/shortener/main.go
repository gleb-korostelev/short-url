package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

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
	grpcservice "github.com/gleb-korostelev/short-url.git/pkg/grpc-service"
	pb "github.com/gleb-korostelev/short-url.git/proto"
	"github.com/gleb-korostelev/short-url.git/tools/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	// Build info
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	// Configuration initialization
	config.ConfigInit()

	// Logger initialization
	log, _ := zap.NewProduction()

	// Storage initialization
	store, err := storageInit()
	if err != nil {
		return
	}
	defer store.Close()

	// Worker Pool initialization
	workerPool := worker.NewDBWorkerPool(config.MaxConcurrentUpdates)
	defer workerPool.Shutdown()

	// API service initialization
	svc := handler.NewAPIService(store, workerPool)

	// Router initialization
	r := router.RouterInit(svc, log)

	// pprof server initialization
	// go func() {
	// 	logger.Infof("Starting pprof server on :6060")
	// 	if err := http.ListenAndServe(":6060", nil); err != nil {
	// 		logger.Fatal("pprof server failed", zap.Error(err))
	// 	}
	// }()

	// grpc Server initialization
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Infof("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterURLServiceServer(grpcServer, &grpcservice.URLServiceServerImpl{})
	if err := grpcServer.Serve(lis); err != nil {
		logger.Infof("Failed to serve: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	server := http.Server{Addr: config.ServerAddr, Handler: r}

	g, gCtx := errgroup.WithContext(ctx)

	if config.EnableHTTPS {
		logger.Infof("Starting HTTPS server on %s\n", config.ServerAddr)
		g.Go(func() error { return server.ListenAndServeTLS(config.CertFilePath, config.KeyFilePath) })
		g.Go(func() error {
			<-gCtx.Done()
			return server.Shutdown(context.Background())
		})

	} else {
		logger.Infof("Starting HTTP server on %s\n", config.ServerAddr)
		g.Go(func() error { return server.ListenAndServe() })
		g.Go(func() error {
			<-gCtx.Done()
			return server.Shutdown(context.Background())
		})
	}
	if err := g.Wait(); err != nil {
		logger.Infof("Exit with: %v\n", err)
	}
}

// Storage initialization for one of inmemory\file\database
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
