package v1

import "github.com/kataras/iris/v12"

// GET /api/v1/user/{uid}
func (handler *Handler) userInfo(ctx iris.Context) {
	uid := ctx.Params().Get("uid")
	ctx.JSON(iris.Map{"request_id": ctx.GetID(), "user": handler.services.UserService.GetByUid(uid)})
}
