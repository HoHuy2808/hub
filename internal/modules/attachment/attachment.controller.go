package attachment

import (
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	var queryParams QueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	queryParams.PostId = postId

	attachments, total, err := a.attachmentService.GetAll(queryParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  attachments,
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
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error":  "Unauthorized",
		})
		return
	}

	UserId, ok := userId.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "Invalid user ID format in token",
		})
		return
	}

	var queryParams QueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	queryParams.UserId = UserId

	attachments, total, err := a.attachmentService.GetAll(queryParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  attachments,
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	var addAttachmentRequest AddAttachmentRequest
	if err := ctx.ShouldBindJSON(&addAttachmentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	addAttachmentRequest.PostId = postId

	result, err := a.attachmentService.AddAttachment(&addAttachmentRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	err = a.attachmentService.DeleteAttachment(attachmentId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete attachment successfully",
	})
}
