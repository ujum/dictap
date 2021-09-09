package config

import (
	"github.com/ujum/dictap/pkg/config/loader"
	"golang.org/x/oauth2"
)

const (
	envPrefix  = "dictup"
	configType = "yaml"
	prefix     = "config"
)

type Config struct {
	ConfigDir  string
	Server     *ServerConfig
	Logger     *LoggerConfig
	Datasource *DatasourceConfig
}

type SecurityConfig struct {
	ApiKeyAuth   *ApiKeyAuthConfig
	GoogleOAuth2 *OAuthConfig
}

type OAuthConfig struct {
	oauth2.Config
	UserInfoURL string
}

type ApiKeyAuthConfig struct {
	AccessTokenMaxAgeMin  int
	RefreshTokenMaxAgeMin int
}

type ServerConfig struct {
	Host     string
	Port     int
	Security *SecurityConfig
}

type LoggerConfig struct {
	Level string
}

type DatasourceConfig struct {
	Mongo *MongoDatasourceConfig
}

type MongoDatasourceConfig struct {
	Port     int
	Database string
	Host     string
}

func New(configDir string) (*Config, error) {
	appConfig := defaultValue()
	settings := &loader.LoadSettings{
		LoadSysEnv: true,
		ConfigFile: &loader.ConfigFileSettings{
			ConfigDir:      configDir,
			FileNamePrefix: prefix,
			ConfigType:     configType,
		},
		EnvPrefix: envPrefix,
	}
	err := loader.Load(appConfig, settings)
	appConfig.ConfigDir = configDir
	return appConfig, err
}
