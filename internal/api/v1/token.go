package v1

import (
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/api"
	"github.com/ujum/dictap/internal/dto"
	"net/http"
)

// auth godoc
// @Summary Create new token pair
// @Tags Token
// @Description Generate access and refresh token
// @Param tokenRequest body dto.UserCredentials true "User credentials"
// @Produce  json
// @Success 200 {object} dto.TokenDTO
// @Failure 400 {object} errResponse
// @Router /auth [post]
func (handler *Handler) auth(ctx iris.Context) {
	credentials := &dto.UserCredentials{}
	err := ctx.ReadJSON(credentials)
	tokens, err := handler.services.TokenService.Generate(api.RequestContext(ctx), credentials)
	if err != nil {
		ctx.StopWithJSON(http.StatusUnauthorized, &errResponse{Message: err.Error()})
		return
	}
	ctx.StopWithJSON(http.StatusOK, tokens)
}

// refresh godoc
// @Summary Generate new token pair by refresh token
// @Tags Token
// @Description Generate access and refresh token by refresh token
// @Param refresh_token query string true "refresh token"
// @Produce  json
// @Success 200 {object} dto.TokenDTO
// @Failure 400 {object} errResponse
// @Router /refresh [post]
func (handler *Handler) refresh(ctx iris.Context) {
	refreshToken := []byte(ctx.URLParam("refresh_token"))
	if len(refreshToken) == 0 {
		ctx.StopWithStatus(http.StatusBadRequest)
		return
	}
	tokens, err := handler.services.TokenService.Refresh(api.RequestContext(ctx), refreshToken)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	ctx.StopWithJSON(http.StatusOK, tokens)
}
