package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func defaultValue() *Config {
	return &Config{
		Server: &ServerConfig{
			Host: "",
			Port: 8080,
			Security: &SecurityConfig{
				ApiKeyAuth: &ApiKeyAuthConfig{
					AccessTokenMaxAgeMin:  30,
					RefreshTokenMaxAgeMin: 60,
				},
				GoogleOAuth2: &OAuthConfig{
					Config: oauth2.Config{
						ClientID:     "",
						ClientSecret: "",
						Endpoint:     google.Endpoint,
						RedirectURL:  "http://localhost:8080/redirect/google",
						Scopes:       []string{},
					},
					UserInfoURL:        "https://www.googleapis.com/oauth2/v2/userinfo?access_token=",
					RedirectOnErrorURL: "/",
				},
				CORS: &CORSConfig{
					AllowOrigin:  "*",
					AllowMethods: "*",
					AllowHeaders: "*",
				},
			},
		},
		Datasource: &DatasourceConfig{
			Mongo: &MongoDatasourceConfig{
				Host:     "localhost",
				Port:     27017,
				Database: "dictup",
			},
		},
		Logger: &LoggerConfig{Level: "info"},
	}
}
