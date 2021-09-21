package v1

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/api"
	"github.com/ujum/dictap/internal/dto"
	derr "github.com/ujum/dictap/internal/error"
	"net/http"
)

const userNotFoundMsg = "user not found"

// userInfo godoc
// @Summary User info
// @Tags Users
// @Description Get user info
// @Produce  json
// @Param uid path string true "user uid"
// @Success 200 {object} dto.User
// @Failure 400 {object} errResponse
// @Failure 404 {object} errResponse
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users/{uid} [get]
func (handler *Handler) userInfo(ctx iris.Context) {
	uid := ctx.Params().Get("uid")
	if uid == "current" {
		currentUID, err := api.GetCurrentUserUID(ctx)
		if err != nil {
			badRequestResponse(ctx, err)
			return
		}
		uid = currentUID
	}
	user, err := handler.services.UserService.GetByUID(api.RequestContext(ctx), uid)
	if err != nil {
		notFoundResponse(ctx, userNotFoundMsg)
		return
	}
	userDTO := &dto.User{}
	if err = copier.Copy(userDTO, user); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithJSON(http.StatusOK, userDTO)
}

// getAllUsers godoc
// @Summary List users
// @Tags Users
// @Description Get all users
// @Produce  json
// @Success 200 {array} dto.User
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users [get]
func (handler *Handler) getAllUsers(ctx iris.Context) {
	users, err := handler.services.UserService.GetAll(api.RequestContext(ctx))
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	usersDTO := &[]*dto.User{}
	err = copier.Copy(usersDTO, users)
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithJSON(http.StatusOK, usersDTO)
}

// createUser godoc
// @Summary Create user
// @Tags Users
// @Description Create new user
// @Accept  json
// @Produce  json
// @Param user body dto.UserCreate true "Create user"
// @Success 201
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users [post]
func (handler *Handler) createUser(ctx iris.Context) {
	user := &dto.UserCreate{}
	if err := ctx.ReadJSON(user); err != nil {
		badRequestResponse(ctx, err)
		return
	}
	requestContext := api.RequestContext(ctx)
	uid, err := handler.services.UserService.Create(requestContext, user)
	if err := err; err != nil {
		if err == derr.ErrAlreadyExists {
			badRequestResponse(ctx, err)
			return
		}
		serverErrorResponse(ctx, err)
		return
	}
	handler.createDefaultGroup(ctx, user, uid, requestContext)
	createdResponse(ctx, uid)
}

func (handler *Handler) createDefaultGroup(ctx iris.Context, user *dto.UserCreate, uid string, requestContext context.Context) {
	if len(user.LangBinding) > 0 {
		wordGroupDTO := dto.WordGroupCreate{
			Name:        "Default",
			UserUID:     uid,
			LangBinding: user.LangBinding[0],
			Default:     true,
		}
		_, err := handler.services.WordGroupService.Create(requestContext, &wordGroupDTO)
		if err != nil {
			serverErrorResponse(ctx, err)
		}
	}
}

// updateUser godoc
// @Summary Update user
// @Tags Users
// @Description Update user
// @Produce  json
// @Param uid path string true "update by uid"
// @Param user body dto.UserUpdate true "update user"
// @Success 200
// @Failure 400 {object} errResponse
// @Failure 404 {object} errResponse
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users/{uid} [patch]
func (handler *Handler) updateUser(ctx iris.Context) {
	uid := ctx.Params().Get("uid")
	if uid == "current" {
		currentUID, err := api.GetCurrentUserUID(ctx)
		if err != nil {
			badRequestResponse(ctx, err)
			return
		}
		uid = currentUID
	}

	user := &dto.UserUpdate{}
	if err := ctx.ReadJSON(user); err != nil {
		badRequestResponse(ctx, err)
		return
	}
	user.UID = uid
	err := handler.services.UserService.Update(api.RequestContext(ctx), user)
	if err != nil {
		if err == derr.ErrNotFound {
			notFoundResponse(ctx, userNotFoundMsg)
			return
		}
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithStatus(http.StatusOK)
}

// deleteUser godoc
// @Summary Delete user
// @Tags Users
// @Description Delete user
// @Produce  json
// @Param uid path string true "delete by uid"
// @Success 200
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users/{uid} [delete]
func (handler *Handler) deleteUser(ctx iris.Context) {
	uid := ctx.Params().Get("uid")
	err := handler.services.UserService.DeleteByUid(api.RequestContext(ctx), uid)
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithStatus(http.StatusOK)
}

// changeUserPass godoc
// @Summary Change user password
// @Tags Users
// @Description Change user password
// @Produce  json
// @Param user body dto.ChangeUserPassword true "Change user password dto"
// @Success 200
// @Failure 400 {object} errResponse
// @Failure 404 {object} errResponse
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users/pass [put]
func (handler *Handler) changeUserPass(ctx iris.Context) {
	uid, err := api.GetCurrentUserUID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	change := &dto.ChangeUserPassword{}
	if err := ctx.ReadJSON(change); err != nil {
		badRequestResponse(ctx, err)
		return
	}
	err = handler.services.UserService.ChangePassword(api.RequestContext(ctx), uid, change)
	if err != nil {
		if err == derr.ErrNotFound {
			notFoundResponse(ctx, userNotFoundMsg)
			return
		}
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithStatus(http.StatusOK)
}
