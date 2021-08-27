package v1

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/requestid"
	"github.com/ujum/dictap/internal/service"
	"github.com/ujum/dictap/pkg/logger"
)

type Handler struct {
	logger   logger.Logger
	services *service.Services
}

func NewHandler(log logger.Logger, services *service.Services) *Handler {
	return &Handler{
		logger:   log,
		services: services,
	}
}

func (handler *Handler) RegisterRoutes(app *iris.Application) {
	app.Use(requestid.New())
	handler.routeV1(app)
}

func (handler *Handler) routeV1(app *iris.Application) {
	v1Group := app.Party("/api/v1")
	{
		userGroup := v1Group.Party("/user")
		userGroup.Get("/{uid}", handler.userInfo)
	}
}
