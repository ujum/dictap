package main

import (
	"flag"
	"github.com/ujum/dictap/internal/app"
	"log"
)

const configDir = "configs"

func main() {
	path := configPath()
	err := app.Run(path)
	if err != nil {
		log.Fatal(err)
	}
}

func configPath() string {
	configPath := flag.String("cfg", configDir, "Config directory path")
	flag.Parse()

	return *configPath
}
