package app

import (
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/app/config"
	"github.com/ujum/dictap/internal/app/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(configFilePath string) {
	appConfig, err := config.New(configFilePath)
	if err != nil {
		log.Printf("can't override default config. Using default. Reason: %v", err)
	}
	srv, err := server.Start(appConfig)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	listenOSStopSignals(srv)
}

func listenOSStopSignals(srv *iris.Application) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	log.Printf("system call: %+v", <-signalChan)
	server.Stop(srv)
}
