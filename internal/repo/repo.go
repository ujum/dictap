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
	return &Repositories{
		UserRepo: NewUserRepoMongo(cfg, log, clients.Mongo),
	}
}
