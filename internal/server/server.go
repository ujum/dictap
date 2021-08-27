package server

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/api/v1"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/service"
	"github.com/ujum/dictap/pkg/logger"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	Logger     logger.Logger
	Iris       *iris.Application
	httpServer *http.Server
	Cfg        *config.ServerConfig
}

// logWriter implement io.Writer interface to adapt
// app logger to log.Logger (http server logger)
type logWriter struct {
	logger logger.Logger
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	lw.logger.Error(p)
	return len(p), nil
}

func New(cfg *config.ServerConfig, appLogger logger.Logger, services *service.Services) *Server {
	irisApp := iris.New()
	irisApp.Logger().Install(appLogger)
	srv := &http.Server{
		Addr:           ":" + strconv.Itoa(cfg.Port),
		Handler:        irisApp,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       log.New(&logWriter{appLogger}, "", 0),
	}

	appSrv := &Server{
		Logger:     appLogger,
		Iris:       irisApp,
		httpServer: srv,
		Cfg:        cfg,
	}
	requestHandler := v1.NewHandler(appLogger, services)
	requestHandler.RegisterRoutes(irisApp)

	return appSrv
}

func (appSrv *Server) Start() {
	go func() {
		if err := appSrv.Iris.Run(iris.Server(appSrv.httpServer)); err != nil {
			if err == http.ErrServerClosed {
				appSrv.Logger.Info("http: web server shutdown complete")
			} else {
				appSrv.Logger.Errorf("http: web server closed unexpect: %s", err)
			}
		}
	}()
}

func (appSrv *Server) Stop() {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := appSrv.Iris.Shutdown(ctxShutDown); err != nil {
		appSrv.Logger.Errorf("http: web server shutdown failed: %+v", err)
		return
	}
	appSrv.Logger.Debug("http: web server closed")
}
