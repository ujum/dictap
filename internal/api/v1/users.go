package v1

import (
	"github.com/jinzhu/copier"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/dto"
	"net/http"
)

// userInfo godoc
// @Summary User info
// @Tags Users
// @Description Get user info
// @Produce  json
// @Param uid path string true "search by uid"
// @Success 200 {object} dto.User
// @Failure 404 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users/{uid} [get]
func (handler *Handler) userInfo(ctx iris.Context) {
	uid := ctx.Params().Get("uid")
	user, err := handler.services.UserService.GetByUid(ctx.Request().Context(), uid)
	if err != nil {
		ctx.StopWithJSON(http.StatusNotFound, errResponse{Message: err.Error()})
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
// @Description get all users
// @Produce  json
// @Success 200 {array} dto.User
// @Failure 404 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users [get]
func (handler *Handler) getAllUsers(ctx iris.Context) {
	users, err := handler.services.UserService.GetAll(ctx.Request().Context())
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
// @Success 200
// @Failure 400 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users [post]
func (handler *Handler) createUser(ctx iris.Context) {
	user := &dto.UserCreate{}
	if err := ctx.ReadJSON(user); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	if err := handler.services.UserService.Create(ctx.Request().Context(), user); err != nil {
		ctx.StopWithJSON(http.StatusBadRequest, &errResponse{Message: err.Error()})
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
// @Failure 404 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users/{uid} [put]
func (handler *Handler) updateUser(ctx iris.Context) {
	uid := ctx.Params().Get("uid")

	user := &dto.UserUpdate{}
	if err := ctx.ReadJSON(user); err != nil {
		serverErrorResponse(ctx, err)
		return
	}

	user.Uid = uid
	err := handler.services.UserService.Update(ctx.Request().Context(), user)
	if err != nil {
		ctx.StopWithJSON(http.StatusNotFound, errResponse{Message: err.Error()})
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
// @Failure 404 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/users/{uid} [delete]
func (handler *Handler) deleteUser(ctx iris.Context) {
	uid := ctx.Params().Get("uid")
	err := handler.services.UserService.DeleteByUid(ctx.Request().Context(), uid)
	if err != nil {
		ctx.StopWithJSON(http.StatusNotFound, errResponse{Message: err.Error()})
		return
	}
	ctx.StopWithStatus(http.StatusOK)
}
