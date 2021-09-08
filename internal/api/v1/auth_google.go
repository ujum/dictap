package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/api"
	"github.com/ujum/dictap/internal/domain"
	"github.com/ujum/dictap/internal/dto"
	"io/ioutil"
	"net/http"
)

const (
	authState = "state"
)

type googleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

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
		fmt.Println(err.Error())
		ctx.Redirect("/", http.StatusTemporaryRedirect)
		return
	}
	googleUserIngo := &googleUser{}
	err = json.Unmarshal(content, googleUserIngo)

	userCreate := &dto.UserCreate{
		User: dto.User{
			Email: googleUserIngo.Email,
			Name:  googleUserIngo.Name,
		},
		OAuthUser: true,
	}

	requestContext := api.RequestContext(ctx)
	_, err = handler.services.UserService.Create(requestContext, userCreate)
	if err != nil && err != domain.ErrUserAlreadyExists {
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
