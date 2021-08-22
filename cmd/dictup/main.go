package main

import (
	"flag"
	"github.com/ujum/dictap/internal/app"
)

const configDir = "configs"

func main() {
	path := configPath()
	app.Run(path)
}

func configPath() string {
	configPath := flag.String("cfg", configDir, "Config directory path")
	flag.Parse()

	return *configPath
}
