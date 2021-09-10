package app

import (
	"context"
	"github.com/ujum/dictap/internal/client"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/internal/server"
	"github.com/ujum/dictap/internal/service"
	"github.com/ujum/dictap/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(configFilePath string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appConfig, err := config.New(configFilePath)
	if err != nil {
		log.Printf("can't override default config. Reason: %+v", err)
		return err
	}
	appLogger := logger.NewLogrus(appConfig)

	go listenStopSignals(ctx, cancel, appLogger)

	deps, err := createDependencies(ctx, appConfig, appLogger)
	defer func() {
		// close all clients connections
		if deps.Clients != nil {
			deps.Clients.Disconnect()
		}
	}()
	if err != nil {
		return err
	}
	srv := server.New(appConfig.Server, appLogger, deps.Services)
	return srv.Start(ctx)
}

func createDependencies(ctx context.Context, cfg *config.Config, appLogger logger.Logger) (*service.Deps, error) {
	clients, err := client.New(ctx, cfg.Datasource, appLogger)
	if err != nil {
		return service.NewDeps(appLogger, clients, nil, nil), err
	}
	repos := repo.New(cfg, appLogger, clients)
	services, err := service.NewServices(cfg, appLogger, repos)
	if err != nil {
		return service.NewDeps(appLogger, clients, repos, services), err
	}
	return service.NewDeps(appLogger, clients, repos, services), nil
}

func listenStopSignals(ctx context.Context, cancel context.CancelFunc, appLogger logger.Logger) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ctx.Done():
		appLogger.Debug("root context has been closed.")
	case sign := <-signalChan:
		appLogger.Infof("system call: %+v", sign)
		cancel()
	}
}
