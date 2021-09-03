package api

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	irisContext "github.com/kataras/iris/v12/context"
)

func GetCurrentUserID(ctx iris.Context) (string, error) {
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