package main

import (
	"flag"
	"github.com/ujum/dictap/internal/app"
	"os"
)

const configDir = "configs"

func main() {
	path := configPath()
	err := app.Run(path)
	if err != nil {
		os.Exit(1)
	}
}

func configPath() string {
	configPath := flag.String("cfg", configDir, "Config directory path")
	flag.Parse()

	return *configPath
}
