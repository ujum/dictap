package v1

import (
	"github.com/kataras/iris/v12"
	"net/http"
)

type errResponse struct {
	Message string `json:"message"`
}

type createResponse struct {
	ID string `json:"id"`
}

func createdResponse(ctx iris.Context, id string) {
	ctx.StopWithJSON(http.StatusAccepted, &createResponse{ID: id})
}

func serverErrorResponse(ctx iris.Context, err error) {
	ctx.StopWithJSON(http.StatusInternalServerError, &errResponse{Message: err.Error()})
}

func badRequestResponse(ctx iris.Context, err error) {
	ctx.StopWithJSON(http.StatusBadRequest, &errResponse{Message: err.Error()})
}
