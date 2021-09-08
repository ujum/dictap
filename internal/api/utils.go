package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/kataras/iris/v12"
	irisContext "github.com/kataras/iris/v12/context"
)

func GetCurrentUserUID(ctx iris.Context) (string, error) {
	user, err := GetCurrentUser(ctx)
	if err != nil {
		return "", err
	}
	id, err := user.GetID()
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetCurrentUser(ctx iris.Context) (irisContext.User, error) {
	if ctx.User() == nil {
		return nil, errors.New("user is nil")
	}
	return ctx.User(), nil
}

func RequestContext(ctx iris.Context) context.Context {
	return ctx.Request().Context()
}

func GenerateRandomString() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(bytes)
	return state, nil
}
