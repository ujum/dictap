package app

import (
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/server"
	"github.com/ujum/dictap/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var appLog logger.Logger

func Run(configFilePath string) {
	appConfig, err := config.New(configFilePath)
	if err != nil {
		log.Printf("can't override default config. Using default. Reason: %+v", err)
	}
	appLog = logger.New(appConfig)
	srv := server.New(appConfig.Server, appLog)
	srv.Start()

	listenOSStopSignals(srv)
}

func listenOSStopSignals(srv *server.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	sign := <-signalChan
	appLog.Infof("system call: %+v", sign)
	srv.Stop()
}
