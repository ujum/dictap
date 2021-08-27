package app

import (
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

func Run(configFilePath string) {
	appConfig, err := config.New(configFilePath)
	if err != nil {
		log.Printf("can't override default config. Using default. Reason: %+v", err)
	}
	appLogger := logger.NewLogrus(appConfig)
	deps, err := createDependencies(appConfig, appLogger)
	if err != nil {
		appLogger.Errorf("can't create app dependencies: %v", err)
		return
	}
	srv := server.New(appConfig.Server, appLogger, deps.Services)
	srv.Start()
	listenOSStopSignals(deps, srv)
}

func createDependencies(cfg *config.Config, appLogger logger.Logger) (*service.Deps, error) {
	clients, err := client.New(cfg.Datasource, appLogger)
	if err != nil {
		return nil, err
	}
	repositories := repo.New(cfg, appLogger, clients)
	return service.New(appLogger, repositories, clients), nil
}

func listenOSStopSignals(deps *service.Deps, srv *server.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	sign := <-signalChan
	deps.Logger.Infof("system call: %+v", sign)
	srv.Stop()
	deps.Clients.Mongo.Disconnect()
}
