package server

import (
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/app/config"
	"net/http"
	"strconv"
	"time"
)

func Start(cfg *config.Config) error {
	app := iris.New()
	route(app)
	srv := &http.Server{
		Addr:           ":" + strconv.Itoa(cfg.Server.Port),
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := app.Run(iris.Server(srv))
	if err != nil {
		return err
	}
	return nil
}

func route(app *iris.Application) {
	app.Get("/", func(ctx iris.Context) {
		_, err := ctx.JSON(iris.Map{"app": "dictup"})
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
		}
	})
}
