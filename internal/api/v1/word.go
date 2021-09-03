package v1

import (
	"github.com/jinzhu/copier"
	"github.com/kataras/iris/v12"
	"github.com/ujum/dictap/internal/api"
	"github.com/ujum/dictap/internal/dto"
	"net/http"
)

// wordsByGroup godoc
// @Summary List words by group
// @Tags Words
// @Description Get words by group
// @Param gid path string true "group id"
// @Produce  json
// @Success 200 {array} dto.Word
// @Failure 404 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/words/groups/{gid} [get]
func (handler *Handler) wordsByGroup(ctx iris.Context) {
	groupID := ctx.Params().Get("gid")
	userID, err := api.GetCurrentUserID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	requestContext := api.RequestContext(ctx)

	// check if group belongs to current user
	wg, err := handler.services.WordGroupService.GetByIDAndUser(requestContext, groupID, userID)
	if err != nil && wg == nil {
		badRequestResponse(ctx, err)
		return
	}

	words, err := handler.services.WordService.GetByGroup(requestContext, groupID)
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	wordsDTO := &[]*dto.Word{}
	err = copier.Copy(wordsDTO, words)
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}

	ctx.StopWithJSON(http.StatusOK, wordsDTO)
}

// createWord godoc
// @Summary Create word
// @Tags Words
// @Description Create new word
// @Accept  json
// @Produce  json
// @Param word body dto.WordCreate true "Word"
// @Success 200
// @Failure 400 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/words [post]
func (handler *Handler) createWord(ctx iris.Context) {
	word := &dto.WordCreate{}
	if err := ctx.ReadJSON(word); err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	userID, err := api.GetCurrentUserID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}

	wordGroup, err := handler.services.WordGroupService.GetByIDAndUser(api.RequestContext(ctx), word.GroupId, userID)
	if err != nil || wordGroup == nil {
		badRequestResponse(ctx, err)
		return
	}

	wordID, err := handler.services.WordService.Create(api.RequestContext(ctx), word)
	if err := err; err != nil {
		badRequestResponse(ctx, err)
		return
	}
	createdResponse(ctx, wordID)
}

// addWordToGroup godoc
// @Summary Add word to group
// @Tags Words
// @Description Add word to group
// @Param name path string true "word name"
// @Param gid path string true "group id"
// @Produce  json
// @Success 200
// @Failure 404 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/words/{name}/groups/{gid} [post]
func (handler *Handler) addWordToGroup(ctx iris.Context) {
	groupID := ctx.Params().Get("gid")
	wordName := ctx.Params().Get("name")

	userID, err := api.GetCurrentUserID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	requestContext := api.RequestContext(ctx)
	_, err = handler.services.WordGroupService.GetByIDAndUser(requestContext, groupID, userID)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	err = handler.services.WordService.AddToGroup(requestContext, wordName, groupID)
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithStatus(http.StatusAccepted)
}

// removeWordFromGroup godoc
// @Summary Remove word from group
// @Tags Words
// @Description Remove word from group
// @Param name path string true "word name"
// @Param gid path string true "group id"
// @Produce  json
// @Success 200
// @Failure 404 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/words/{name}/groups/{gid} [delete]
func (handler *Handler) removeWordFromGroup(ctx iris.Context) {
	groupID := ctx.Params().Get("gid")
	wordName := ctx.Params().Get("name")

	userID, err := api.GetCurrentUserID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	requestContext := api.RequestContext(ctx)
	_, err = handler.services.WordGroupService.GetByIDAndUser(requestContext, groupID, userID)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	err = handler.services.WordService.RemoveFromGroup(requestContext, wordName, groupID)
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithStatus(http.StatusAccepted)
}

// moveWordToGroup godoc
// @Summary Move word to group
// @Tags Words
// @Description Move word to group
// @Param name path string true "word name"
// @Param move body dto.WordGroupMovement true "Word Group Movement"
// @Produce  json
// @Success 200
// @Failure 404 {object} errResponse
// @Security ApiKeyAuth
// @Router /api/v1/words/{name}/groups [post]
func (handler *Handler) moveWordToGroup(ctx iris.Context) {
	wordName := ctx.Params().Get("name")

	wgMove := &dto.WordGroupMovement{}
	if err := ctx.ReadJSON(wgMove); err != nil {
		serverErrorResponse(ctx, err)
		return
	}

	userID, err := api.GetCurrentUserID(ctx)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}
	requestContext := api.RequestContext(ctx)
	_, err = handler.services.WordGroupService.GetByIDAndUser(requestContext, wgMove.FromGroupId, userID)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}

	_, err = handler.services.WordGroupService.GetByIDAndUser(requestContext, wgMove.ToGroupId, userID)
	if err != nil {
		badRequestResponse(ctx, err)
		return
	}

	err = handler.services.WordService.MoveToGroup(requestContext, wordName, wgMove.FromGroupId, wgMove.ToGroupId)
	if err != nil {
		serverErrorResponse(ctx, err)
		return
	}
	ctx.StopWithStatus(http.StatusAccepted)
}
