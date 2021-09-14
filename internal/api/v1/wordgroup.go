package v1

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/api"
	"github.com/ujum/dictap/internal/dto"
	derr "github.com/ujum/dictap/internal/error"
	"net/http"
)

const wordGroupNotFoundMsg = "word group not found"

// createWordGroup godoc
// @Summary Create word group
// @Tags WordGroups
// @Description Create new word group
// @Accept  json
// @Produce  json
// @Param word body dto.WordGroupCreate true "Word Group"
// @Success 201
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/wordgroups [post]
func (handler *Handler) createWordGroup(ctx iris.Context) {
	wordGroup := &dto.WordGroupCreate{}
	if err := ctx.ReadJSON(wordGroup); err != nil {
		badRequestResponse(ctx, err)
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
		serverErrorResponse(ctx, err)
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
// @Failure 404 {object} errResponse
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/wordgroups/{gid} [get]
func (handler *Handler) getWordGroup(ctx iris.Context) {
	gid := ctx.Params().Get("gid")
	if gid == "" {
		badRequestResponse(ctx, errors.New("param gid not provided"))
		return
	}
	userUID, err := api.GetCurrentUserUID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroup, err := handler.services.WordGroupService.GetByIDAndUser(api.RequestContext(ctx), gid, userUID)
	if err := err; err != nil {
		if err == derr.ErrNotFound {
			notFoundResponse(ctx, wordGroupNotFoundMsg)
			return
		}
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

// getWordGroupsByLang godoc
// @Summary List word groups for language
// @Tags WordGroups
// @Description Get all word groups for language
// @Param from_iso path string true "from language iso code"
// @Param to_iso path string true "to language iso code"
// @Produce  json
// @Success 200 {array} dto.WordGroup
// @Failure 400 {object} errResponse
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/wordgroups/langs/{from_iso}/{to_iso} [get]
func (handler *Handler) getWordGroupsByLang(ctx iris.Context) {
	fromISO := ctx.Params().Get("from_iso")
	if fromISO == "" {
		badRequestResponse(ctx, errors.New("param from_iso not provided"))
		return
	}
	toISO := ctx.Params().Get("to_iso")
	if toISO == "" {
		badRequestResponse(ctx, errors.New("param to_iso not provided"))
		return
	}
	langBinding := &dto.LangBinding{LangFromISO: fromISO, LangToISO: toISO}

	userUID, err := api.GetCurrentUserUID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroups, err := handler.services.WordGroupService.GetAllByLangAndUser(api.RequestContext(ctx), langBinding, userUID)
	if err := err; err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	wordGroupsDTO := &[]*dto.WordGroup{}
	if err = copier.Copy(wordGroupsDTO, wordGroups); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithJSON(http.StatusOK, wordGroupsDTO)
}

// getDefaultWordGroupByLang godoc
// @Summary Get default word group for language
// @Tags WordGroups
// @Description Get word group for language
// @Param from_iso path string true "from lang iso code"
// @Param to_iso path string true "to lang iso code"
// @Produce  json
// @Success 200 {object} dto.WordGroup
// @Failure 400 {object} errResponse
// @Failure 404
// @Failure 500 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/wordgroups/langs/{from_iso}/{to_iso}/default [get]
func (handler *Handler) getDefaultWordGroupByLang(ctx iris.Context) {
	fromISO := ctx.Params().Get("from_iso")
	if fromISO == "" {
		badRequestResponse(ctx, errors.New("param from_iso not provided"))
		return
	}
	toISO := ctx.Params().Get("to_iso")
	if toISO == "" {
		badRequestResponse(ctx, errors.New("param to_iso not provided"))
		return
	}
	langBinding := &dto.LangBinding{LangFromISO: fromISO, LangToISO: toISO}

	userUID, err := api.GetCurrentUserUID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	wordGroup, err := handler.services.WordGroupService.GetDefault(api.RequestContext(ctx), langBinding, userUID)
	if err := err; err != nil {
		if err == derr.ErrNotFound {
			notFoundResponse(ctx, wordGroupNotFoundMsg)
			return
		}
		serverErrorResponse(ctx, err)
		return
	}
	wordGroupDTO := &dto.WordGroup{}
	if err = copier.Copy(wordGroupDTO, wordGroup); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithJSON(http.StatusOK, wordGroupDTO)
}
