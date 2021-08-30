package repo

import (
	"github.com/ujum/dictap/internal/client"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/pkg/logger"
)

type Repositories struct {
	UserRepo UserRepo
}

func New(cfg *config.Config, log logger.Logger, clients *client.Clients) *Repositories {
	mongoDatabase := clients.Mongo.Client.Database(cfg.Datasource.Mongo.Schema)
	return &Repositories{
		UserRepo: NewUserRepoMongo(cfg, log, mongoDatabase.Collection("users")),
	}
}
