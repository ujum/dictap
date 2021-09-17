package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/api"
	"github.com/ujum/dictap/internal/dto"
	derr "github.com/ujum/dictap/internal/error"
	"io/ioutil"
	"net/http"
)

const authState = "state"

type googleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// googleLogin godoc
// @Summary Google Login
// @Tags Token
// @Description Sign up/in with Google
// @Produce  json
// @Success 307
// @Failure 500 {object} errResponse
// @Router /auth/google [get]
func (handler *Handler) googleLogin(ctx iris.Context) {
	state, err := api.GenerateRandomString()
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	url := handler.config.Security.GoogleOAuth2.AuthCodeURL(state)
	ctx.SetCookie(&http.Cookie{
		Name:     authState,
		Value:    state,
		Path:     "/",
		MaxAge:   10,
		HttpOnly: true,
	})
	ctx.Redirect(url, http.StatusTemporaryRedirect)
}

func (handler *Handler) googleCallback(ctx iris.Context) {
	stateFromCookie := ctx.GetCookie(authState)
	stateFormGoogle := ctx.FormValue(authState)
	content, err := handler.getGoogleUserInfo(ctx.FormValue("code"))
	if stateFromCookie != stateFormGoogle {
		badRequestResponse(ctx, errors.New("invalid oauth state"))
		return
	}
	if err != nil {
		handler.logger.Errorf("cant get google user info: %v", err)
		ctx.Redirect(handler.config.Security.GoogleOAuth2.RedirectOnErrorURL, http.StatusTemporaryRedirect)
		return
	}
	googleUserInfo := &googleUser{}
	err = json.Unmarshal(content, googleUserInfo)

	userCreate := &dto.UserCreate{
		Email:     googleUserInfo.Email,
		Name:      googleUserInfo.Name,
		OAuthUser: true,
	}

	requestContext := api.RequestContext(ctx)
	_, err = handler.services.UserService.Create(requestContext, userCreate)
	if err != nil && err != derr.ErrAlreadyExists {
		badRequestResponse(ctx, err)
		return
	}
	user, err := handler.services.UserService.GetByEmail(requestContext, userCreate.Email)
	if err != nil {
		ctx.StopWithJSON(http.StatusNotFound, errResponse{Message: err.Error()})
		return
	}

	tokens, err := handler.services.TokenService.GenerateForUser(user)
	if err != nil {
		ctx.StopWithJSON(http.StatusUnauthorized, &errResponse{Message: err.Error()})
		return
	}

	if err = handler.services.UserService.FlagUserAsOAuth(requestContext, user); err != nil {
		handler.logger.Warnf("can't flag user %s as oauth: %v", user.Email, err)
	}

	ctx.StopWithJSON(http.StatusOK, tokens)

}
func (handler *Handler) getGoogleUserInfo(code string) ([]byte, error) {
	token, err := handler.config.Security.GoogleOAuth2.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get(handler.config.Security.GoogleOAuth2.UserInfoURL + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, nil
}
