package service

import (
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/ujum/dictap/internal/client"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/pkg/logger"
)

type Services struct {
	UserService  UserService
	TokenService TokenService
	JwtVerifier  *jwt.Verifier
}

type Deps struct {
	Logger   logger.Logger
	Clients  *client.Clients
	Repos    *repo.Repositories
	Services *Services
}

func NewServices(cfg *config.Config, appLogger logger.Logger, repos *repo.Repositories) (*Services, error) {
	userService := newUserService(repos, appLogger)
	verifier, err := NewJwtVerifier(cfg)
	if err != nil {
		return nil, err
	}
	signer, err := NewJwtSigner(cfg)
	if err != nil {
		return nil, err
	}

	jwtTokenService := NewJwtTokenService(cfg, appLogger, verifier, signer, userService)
	return &Services{
		UserService:  userService,
		TokenService: jwtTokenService,
		JwtVerifier:  verifier,
	}, nil
}

func NewDeps(appLogger logger.Logger, clients *client.Clients, repos *repo.Repositories, services *Services) *Deps {
	return &Deps{
		Logger:   appLogger,
		Repos:    repos,
		Clients:  clients,
		Services: services,
	}
}
