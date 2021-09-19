package v1

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	_ "github.com/ujum/dictap/api"
	"github.com/ujum/dictap/internal/config"
	"github.com/ujum/dictap/internal/service"
	"github.com/ujum/dictap/pkg/logger"
	"net/http"
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
	app.UseRouter(handler.corsMiddleware)
	app.Validator = validator.New()
	handler.routeV1(app)
}

func (handler *Handler) routeV1(app *iris.Application) {
	app.Post("/auth", handler.auth)
	app.Post("/refresh", handler.refresh)
	app.Get("/auth/google", handler.googleLogin)
	app.Get("/auth/google/callback", handler.googleCallback)

	tokenVerifyHandler := handler.services.JwtVerifier.Verify(func() interface{} {
		return new(context.SimpleUser)
	})

	v1Group := app.Party("/api/v1", tokenVerifyHandler)
	{
		userGroup := v1Group.Party("/users")
		{
			userGroup.Get("/", handler.getAllUsers)
			userGroup.Get("/{uid}", handler.userInfo)
			userGroup.Post("/", handler.createUser)
			userGroup.Patch("/{uid}", handler.updateUser)
			userGroup.Delete("/{uid}", handler.deleteUser)
			userGroup.Put("/{uid}/pass", handler.changeUserPass)
		}
		wordGroup := v1Group.Party("/words")
		{
			wordGroup.Get("/", handler.getAllUsers)
			wordGroup.Get("/groups/{gid}", handler.wordsByGroup)
			wordGroup.Post("/", handler.createWord)
			wordGroup.Post("/{name}/groups/{gid}", handler.addWordToGroup)
			wordGroup.Post("/{name}/groups", handler.moveWordToGroup)
			wordGroup.Delete("/{name}/groups/{gid}", handler.removeWordFromGroup)
		}
		wordGroupGroup := v1Group.Party("/wordgroups")
		{
			wordGroupGroup.Get("/langs/{from_iso}/{to_iso}", handler.getWordGroupsByLang)
			wordGroupGroup.Post("/", handler.createWordGroup)
			wordGroupGroup.Get("/{gid}", handler.getWordGroup)
			wordGroupGroup.Get("/langs/{from_iso}/{to_iso}/default", handler.getDefaultWordGroupByLang)
		}
	}
}

func (handler *Handler) routeSwagger(app *iris.Application) {
	const protocol = "http"
	var host string
	if handler.config.Host == "" {
		host = "localhost"
	}
	hostPort := fmt.Sprintf("%s:%d", host, handler.config.Port)
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

func (handler *Handler) corsMiddleware(ctx *context.Context) {
	ctx.Header("Access-Control-Allow-Origin", handler.config.Security.CORS.AllowOrigin)
	ctx.Header("Access-Control-Allow-Methods", handler.config.Security.CORS.AllowMethods)
	ctx.Header("Access-Control-Allow-Headers", handler.config.Security.CORS.AllowHeaders)
	if ctx.Request().Method != "OPTIONS" {
		ctx.Next()
	} else {
		ctx.StopWithStatus(http.StatusOK)
	}
}
