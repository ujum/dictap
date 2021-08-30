package config

func defaultValue() *Config {
	return &Config{
		Server: &ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		Datasource: &DatasourceConfig{
			Mongo: &MongoDatasourceConfig{
				Host:   "localhost",
				Port:   27017,
				Schema: "dictup",
			},
		},
		Logger: &LoggerConfig{Level: "info"},
	}
}
