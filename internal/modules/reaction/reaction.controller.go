package reaction

import (
	"net/http"

	"hub/internal/websocket"
	"hub/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReactionController struct {
	reactionService ReactionService
	hub             *websocket.Hub
}

func NewReactionController(reactionService ReactionService, hub *websocket.Hub) *ReactionController {
	return &ReactionController{reactionService: reactionService, hub: hub}
}

// GetAllReactions		godoc
// @Security 		BearerAuth
// @Summary			Get all reactions
// @Description 	Get all reactions
// @Tags        	Reaction
// @Produce     	json
// @Param       	postId path string true "Post ID"
// @Param       	limit query int false "Number of items per page"
// @Param       	offset query int false "Offset"
// @Param			type query string false "Type of reaction"
// @Router      	/reactions/{postId} [get]
func (r *ReactionController) GetAll(ctx *gin.Context) {
	postIdParam := ctx.Param("postId")
	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Post ID format")
		return
	}

	queryParams := &QueryParams{}
	if err := ctx.ShouldBindQuery(queryParams); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	queryParams.PostId = postId

	reactions, total, err := r.reactionService.GetAll(ctx.Request.Context(), queryParams)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Get all reactions successfully", gin.H{
		"reactions": reactions,
		"total":     total,
	})
}

// React 			godoc
// @Security 		BearerAuth
// @Summary			React to a post
// @Description 	React to a post
// @Tags        	Reaction
// @Produce     	json
// @Param       	postId path string true "Post ID"
// @Param       	reactRequest body ReactRequest true "Reaction request"
// @Success     	200 {object} ReactResponse
// @Router      	/reactions/{postId} [post]
func (r *ReactionController) React(ctx *gin.Context) {
	postIdParam := ctx.Param("postId")
	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Post ID format")
		return
	}

	reactRequest := &ReactRequest{}
	if err := ctx.ShouldBindJSON(reactRequest); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	reactRequest.PostId = postId
	reactRequest.ReactorId = userId.(uuid.UUID)

	result, err := r.reactionService.React(ctx.Request.Context(), reactRequest)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	msg := &websocket.NotificationMessage{
		TargetUserID: result.PostOwnerId.String(),
		Data:         []byte(result.Type + " vào bài viết của bạn"),
	}
	r.hub.SendToUser <- msg

	response.Success(ctx, http.StatusOK, "Gửi cảm xúc thành công", result)
}

// UpdateReact		godoc
// @Security 		BearerAuth
// @Summary			Update a reaction
// @Description 	Update user's reaction to a post
// @Tags        	Reaction
// @Produce     	json
// @Param       	postId path string true "Post ID"
// @Param       	updateReactRequest body UpdateReactRequest true "Update reaction request"
// @Success     	200 {object} UpdateReactResponse
// @Router      	/reactions/{postId} [patch]
func (r *ReactionController) UpdateReact(ctx *gin.Context) {
	postIdParam := ctx.Param("postId")
	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Post ID format")
		return
	}

	updateReactRequest := &UpdateReactRequest{}
	if err := ctx.ShouldBindJSON(updateReactRequest); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	updateReactRequest.PostId = postId
	updateReactRequest.ReactorId = userId.(uuid.UUID)

	result, err := r.reactionService.UpdateReact(ctx.Request.Context(), updateReactRequest)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Thay đổi cảm xúc thành công", result)
}

// Unreact			godoc
// @Security 		BearerAuth
// @Summary			Unreact to a post
// @Description 	Unreact to a post
// @Tags        	Reaction
// @Produce     	json
// @Param       	postId path string true "Post ID"
// @Success     	200 {object} ReactResponse
// @Router      	/reactions/{postId} [delete]
func (r *ReactionController) Unreact(ctx *gin.Context) {
	postIdParam := ctx.Param("postId")
	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid Post ID format")
		return
	}

	userId, _ := ctx.Get("userId")

	err = r.reactionService.Unreact(ctx.Request.Context(), postId, userId.(uuid.UUID))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Gỡ cảm xúc thành công", nil)
}
