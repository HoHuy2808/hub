package attachment

import (
	"net/http"

	"hub/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentController struct {
	attachmentService AttachmentService
}

func NewAttachmentController(attachmentService AttachmentService) *AttachmentController {
	return &AttachmentController{attachmentService: attachmentService}
}

// GetAll godoc
// @Security		BearerAuth
// @Summary			Get all attachments
// @Description 	Get all attachments
// @Tags        	Attachment
// @Produce     	json
// @Param 			post_id path string true "Post ID"
// @Param			page query int false "Page number" default(1)
// @Param			limit query int false "Number of items per page" default(10)
// @Success     	200 {object} AddAttachmentResponse
// @Router      	/posts/{post_id}/attachments [get]
func (a *AttachmentController) GetAll(ctx *gin.Context) {
	postIdParam := ctx.Param("post_id")

	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var queryParams QueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	queryParams.PostId = postId

	attachments, total, err := a.attachmentService.GetAll(ctx.Request.Context(), queryParams)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Lấy danh sách đính kèm thành công", gin.H{
		"items": attachments,
		"total": total,
	})
}

// GetAllByUser godoc
// @Security		BearerAuth
// @Summary			Get all attachments by user
// @Description 	Get all attachments for a specific user
// @Tags        	Attachment
// @Produce     	json
// @Param			page query int false "Page number" default(1)
// @Param			limit query int false "Number of items per page" default(10)
// @Success     	200 {object} AddAttachmentResponse
// @Router      	/users/me/attachments [get]
func (a *AttachmentController) GetAllByUser(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	UserId, ok := userId.(uuid.UUID)
	if !ok {
		response.Error(ctx, http.StatusInternalServerError, "Invalid user ID format in token")
		return
	}

	var queryParams QueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	queryParams.UserId = UserId

	attachments, total, err := a.attachmentService.GetAll(ctx.Request.Context(), queryParams)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Lấy danh sách đính kèm của người dùng thành công", gin.H{
		"items": attachments,
		"total": total,
	})
}

// AddAttachment	godoc
// @Security		BearerAuth
// @Summary			Add attachment to a post
// @Description 	Add attachment to a post
// @Tags        	Attachment
// @Produce     	json
// @Param 			post_id path string true "Post ID"
// @Param       	content body AddAttachmentRequest true "Attachment data"
// @Success     	200 {object} AddAttachmentResponse
// @Router      	/posts/{post_id}/attachments [post]
func (a *AttachmentController) AddAttachment(ctx *gin.Context) {
	postIdParam := ctx.Param("post_id")

	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var addAttachmentRequest AddAttachmentRequest
	if err := ctx.ShouldBindJSON(&addAttachmentRequest); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	addAttachmentRequest.PostId = postId

	result, err := a.attachmentService.AddAttachment(ctx.Request.Context(), &addAttachmentRequest)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Thêm đính kèm thành công", result)
}

// DeleteAttachment godoc
// @Security		BearerAuth
// @Summary			Delete attachment
// @Description 	Delete attachment
// @Tags        	Attachment
// @Produce     	json
// @Param 			id path string true "Attachment ID"
// @Success     	200 {object} AddAttachmentResponse
// @Router      	/posts/attachments/{id} [delete]
func (a *AttachmentController) DeleteAttachment(ctx *gin.Context) {
	attachmentIdParam := ctx.Param("id")

	attachmentId, err := uuid.Parse(attachmentIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = a.attachmentService.DeleteAttachment(ctx.Request.Context(), attachmentId)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Xóa đính kèm thành công", nil)
}
