package config

func defaultValue() *Config {
	return &Config{
		Server: &ServerConfig{
			Port: 8080,
		},
		Datasource: &DatasourceConfig{
			Mongo: &MongoDatasourceConfig{
				Port: 27017,
				Host: "localhost",
			},
		},
		Logger: &LoggerConfig{Level: "info"},
	}
}
