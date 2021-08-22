package app

import (
	"github.com/ujum/dictap/internal/app/config"
	"github.com/ujum/dictap/internal/app/server"
	"log"
)

func Run(configFilePath string) {
	appConfig, err := config.New(configFilePath)
	if err != nil {
		log.Printf("Can't override default config. Using default. Reason: %v", err)
	}
	if err := server.Start(appConfig); err != nil {
		log.Fatalf("%v", err)
	}
}
