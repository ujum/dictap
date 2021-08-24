package server

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	iris   *iris.Application
	runner iris.Runner
	cfg    *config.ServerConfig
	logger logger.Logger
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

func New(cfg *config.ServerConfig, appLogger logger.Logger) *Server {
	irisApp := iris.New()
	irisApp.Logger().Install(appLogger)

	route(irisApp)
	srv := &http.Server{
		Addr:           ":" + strconv.Itoa(cfg.Port),
		Handler:        irisApp,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		ErrorLog:       log.New(&logWriter{appLogger}, "", 0),
	}
	runner := iris.Server(srv)

	return &Server{
		iris:   irisApp,
		runner: runner,
		cfg:    cfg,
		logger: appLogger,
	}
}

func (appSrv *Server) Start() {
	go func() {
		if err := appSrv.iris.Run(appSrv.runner); err != nil {
			appSrv.logger.Errorf("error during server starting, %+v", err)
		}
	}()
}

func (appSrv *Server) Stop() {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := appSrv.iris.Shutdown(ctxShutDown); err != nil {
		appSrv.logger.Warnf("server shutdown failed: %+v", err)
	}
	appSrv.logger.Info("http: Server closed")
}

func route(app *iris.Application) {
	app.Get("/", func(ctx iris.Context) {
		_, err := ctx.JSON(iris.Map{"app": "dictup"})
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
		}
	})
}
