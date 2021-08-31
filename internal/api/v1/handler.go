package v1

import (
	"fmt"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	_ "github.com/ujum/dictap/api"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/service"
	"github.com/ujum/dictap/pkg/logger"
)

type Handler struct {
	logger   logger.Logger
	config   *config.ServerConfig
	services *service.Services
}

func NewHandler(log logger.Logger, cfg *config.ServerConfig, services *service.Services) *Handler {
	return &Handler{
		logger:   log,
		services: services,
		config:   cfg,
	}
}

func (handler *Handler) RegisterRoutes(app *iris.Application) {
	handler.routeSwagger(app)
	app.Use(requestid.New())
	handler.routeV1(app)
}

func (handler *Handler) routeV1(app *iris.Application) {
	app.Post("/auth", handler.auth)
	app.Post("/refresh", handler.refresh)

	v1Group := app.Party("/api/v1", handler.services.TokenService.VerifyHandler())
	{
		userGroup := v1Group.Party("/users")
		{
			userGroup.Get("/", handler.getAllUsers)
			userGroup.Get("/{uid}", handler.userInfo)
			userGroup.Post("/", handler.createUser)
			userGroup.Put("/{uid}", handler.updateUser)
			userGroup.Delete("/{uid}", handler.deleteUser)
		}
	}
}

func (handler *Handler) routeSwagger(app *iris.Application) {
	protocol := "http"
	hostPort := fmt.Sprintf("%s:%d", handler.config.Host, handler.config.Port)
	url := protocol + "://" + hostPort + "/swagger/doc.json"

	swaggerUI := swagger.CustomWrapHandler(
		&swagger.Config{
			URL:         url,
			DeepLinking: true,
		},
		swaggerFiles.Handler)

	app.Get("/swagger", swaggerUI)
	app.Get("/swagger/{any:path}", swaggerUI)
}
