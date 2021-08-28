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
	Clients  *client.Clients
	Repos    *repo.Repositories
	Services *Services
}

func NewServices(appLogger logger.Logger, repos *repo.Repositories) *Services {
	return &Services{
		UserService: &UserServiceImpl{
			userRepo: repos.UserRepo,
			logger:   appLogger,
		},
	}
}

func NewDeps(appLogger logger.Logger, clients *client.Clients, repos *repo.Repositories, services *Services) *Deps {
	return &Deps{
		Logger:   appLogger,
		Repos:    repos,
		Clients:  clients,
		Services: services,
	}
}
