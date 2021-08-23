package config

func defaultValue() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		}}
}
