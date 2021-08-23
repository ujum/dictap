package server

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/app/config"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Start(cfg *config.Config) (*iris.Application, error) {
	app := iris.New()
	route(app)
	srv := &http.Server{
		Addr:           ":" + strconv.Itoa(cfg.Server.Port),
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := app.Run(iris.Server(srv)); err != nil {
			log.Fatalf("%v", err)
		}
	}()

	return app, nil
}

func Stop(app *iris.Application) {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := app.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown Failed: %v", err)
	}
	log.Printf("http: Server closed")
}

func route(app *iris.Application) {
	app.Get("/", func(ctx iris.Context) {
		_, err := ctx.JSON(iris.Map{"app": "dictup"})
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
		}
	})
}
