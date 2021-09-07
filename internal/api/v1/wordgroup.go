package v1

import (
	"github.com/jinzhu/copier"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/api"
	"github.com/ujum/dictap/internal/dto"
	"net/http"
)

// createWordGroup godoc
// @Summary Create word group
// @Tags WordGroups
// @Description Create new word group
// @Accept  json
// @Produce  json
// @Param word body dto.WordGroupCreate true "Word Group"
// @Success 200
// @Failure 400 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/wordgroups [post]
func (handler *Handler) createWordGroup(ctx iris.Context) {

	wordGroup := &dto.WordGroupCreate{}
	if err := ctx.ReadJSON(wordGroup); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	userUID, err := api.GetCurrentUserUID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroup.UserUID = userUID
	wordGroupID, err := handler.services.WordGroupService.Create(api.RequestContext(ctx), wordGroup)
	if err := err; err != nil {
		badRequestResponse(ctx, err)
		return
	}
	createdResponse(ctx, wordGroupID)
}

// getWordGroup godoc
// @Summary Word group by id
// @Tags WordGroups
// @Description Get word group by id
// @Param gid path string true "group id"
// @Produce  json
// @Success 200 {object} dto.WordGroup
// @Failure 400 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/wordgroups/{gid} [get]
func (handler *Handler) getWordGroup(ctx iris.Context) {
	gid := ctx.Params().Get("gid")
	userUID, err := api.GetCurrentUserUID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroup, err := handler.services.WordGroupService.GetByIDAndUser(api.RequestContext(ctx), gid, userUID)
	if err := err; err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroupDTO := &dto.WordGroup{}
	if err = copier.Copy(wordGroupDTO, wordGroup); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithJSON(http.StatusOK, wordGroupDTO)
}

// getWordGroups godoc
// @Summary List word groups
// @Tags WordGroups
// @Description Get all word groups
// @Param lid path string true "language id"
// @Produce  json
// @Success 200 {array} dto.WordGroup
// @Failure 400 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/wordgroups/langs/{lid} [get]
func (handler *Handler) getWordGroups(ctx iris.Context) {
	lid := ctx.Params().Get("lid")

	userUID, err := api.GetCurrentUserUID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroups, err := handler.services.WordGroupService.GetAllByLangAndUser(api.RequestContext(ctx), lid, userUID)
	if err := err; err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroupsDTO := &[]*dto.WordGroup{}
	if err = copier.Copy(wordGroupsDTO, wordGroups); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithJSON(http.StatusOK, wordGroupsDTO)
}

// getDefaultWordGroup godoc
// @Summary Get default word group by language
// @Tags WordGroups
// @Description Get word group by language
// @Param lid path string true "lang id"
// @Produce  json
// @Success 200 {object} dto.WordGroup
// @Failure 400 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/wordgroups/default/{lid} [get]
func (handler *Handler) getDefaultWordGroup(ctx iris.Context) {
	lid := ctx.Params().Get("lid")
	userUID, err := api.GetCurrentUserUID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroup, err := handler.services.WordGroupService.GetDefault(api.RequestContext(ctx), lid, userUID)
	if err := err; err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroupDTO := &dto.WordGroup{}
	if err = copier.Copy(wordGroupDTO, wordGroup); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithJSON(http.StatusOK, wordGroupDTO)
}
