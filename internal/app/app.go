package app

import (
	"context"
	"github.com/ujum/dictap/internal/client"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/internal/server"
	"github.com/ujum/dictap/internal/service"
	"github.com/ujum/dictap/pkg/logger"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(configFilePath string) error {
	ctx, cancel := context.WithCancel(context.Background())

	appConfig, err := config.New(configFilePath)
	if err != nil {
		log.Printf("can't override default config. Using default. Reason: %+v", err)
	}
	appLogger := logger.NewLogrus(appConfig)
	depsCtx, depsCtxCancel := context.WithCancel(context.Background())

	deps, err := createDependencies(depsCtx, appConfig, appLogger)
	if err != nil {
		appLogger.Errorf("can't create app dependencies: %v", err)
		//cancel root and deps context
		cancel()
		depsCtxCancel()
		deps.Clients.WaitDisconnect()
		return err
	}

	srv := server.New(appConfig.Server, appLogger, deps.Services)

	var g errgroup.Group

	g.Go(func() error {
		servErr := srv.Start(ctx)
		// cancel deps context when server is stopped
		depsCtxCancel()
		deps.Clients.WaitDisconnect()
		return servErr
	})
	go listenStopSignals(ctx, cancel, srv, deps)
	if err = g.Wait(); err != nil {
		return err
	}
	return nil
}

func createDependencies(ctx context.Context, cfg *config.Config, appLogger logger.Logger) (*service.Deps, error) {
	clients, err := client.New(ctx, cfg.Datasource, appLogger)
	if err != nil {
		return service.NewDeps(nil, clients, nil, nil), err
	}
	repos := repo.New(cfg, appLogger, clients)
	services := service.NewServices(appLogger, repos)
	return service.NewDeps(appLogger, clients, repos, services), nil
}

func listenStopSignals(ctx context.Context, cancel context.CancelFunc, srv *server.Server, deps *service.Deps) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	select {
	case <-ctx.Done():
		deps.Logger.Debug("root context has been closed. Stopping server...")
		if err := srv.Stop(); err != nil {
			deps.Logger.Errorf("http: web server shutdown failed: %+v", err)
		}
	case sign := <-signalChan:
		deps.Logger.Infof("system call: %+v", sign)
		cancel()
	}
}
