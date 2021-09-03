package service

import (
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/ujum/dictap/internal/client"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/repo"
	"github.com/ujum/dictap/pkg/logger"
)

type Services struct {
	UserService      UserService
	TokenService     TokenService
	JwtVerifier      *jwt.Verifier
	WordService      WordService
	WordGroupService WordGroupService
}

type Deps struct {
	Logger   logger.Logger
	Clients  *client.Clients
	Repos    *repo.Repositories
	Services *Services
}

func NewServices(cfg *config.Config, appLogger logger.Logger, repos *repo.Repositories) (*Services, error) {
	userService := newUserService(repos, appLogger)
	verifier, err := newJwtVerifier(cfg)
	if err != nil {
		return nil, err
	}
	signer, err := newJwtSigner(cfg)
	if err != nil {
		return nil, err
	}

	jwtTokenService := newJwtTokenService(cfg, appLogger, verifier, signer, userService)

	wordGroupService := newWordGroupService(repos, appLogger)
	wordService := newWordService(repos, appLogger)

	return &Services{
		UserService:      userService,
		TokenService:     jwtTokenService,
		JwtVerifier:      verifier,
		WordService:      wordService,
		WordGroupService: wordGroupService,
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
