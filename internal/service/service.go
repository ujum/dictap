package service

import (
	"github.com/ujum/dictap/internal/client"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/pkg/logger"
)

type Services struct {
	UserService UserService
}

type Deps struct {
	Logger   logger.Logger
	Services *Services
	Clients  *client.Clients
}

func New(appLogger logger.Logger, repos *repo.Repositories, clients *client.Clients) *Deps {
	srvs := &Services{
		UserService: &UserServiceImpl{
			userRepo: repos.UserRepo,
			logger:   appLogger,
		},
	}
	return &Deps{
		Logger:   appLogger,
		Services: srvs,
		Clients:  clients,
	}
}
