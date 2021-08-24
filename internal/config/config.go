package config

import (
	"github.com/ujum/dictap/pkg/config/loader"
)

const (
	envPrefix  = "dictup"
	configType = "yaml"
	prefix     = "config"
)

type Config struct {
	Server   *ServerConfig
	LogLevel string
}

type ServerConfig struct {
	Port int
}

func New(configDir string) (*Config, error) {
	appConfig := defaultValue()
	settings := &loader.LoadSettings{
		LoadSysEnv:     true,
		ConfigDir:      configDir,
		FileNamePrefix: prefix,
		ConfigType:     configType,
		EnvPrefix:      envPrefix,
	}
	err := loader.Load(appConfig, settings)
	return appConfig, err
}
