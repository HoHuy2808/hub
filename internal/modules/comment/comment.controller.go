package comment

import (
	"net/http"

	"hub/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentController struct {
	commentService CommentService
}

func NewCommentController(commentService CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

// GetAll 			godoc
// @Security 		BearerAuth
// @Summary			Get all comments
// @Description 	Get all comments of a specific post
// @Tags        	Comment
// @Produce     	json
// @Param       	postId path string true "Post ID"
// @Param       	limit query int false "Number of items per page"
// @Param       	offset query int false "Offset"
// @Param       	sortBy query string false "Sort by"
// @Param       	sortOrder query string false "Sort order"
// @Success     	200 {object} CreateCommentResponse
// @Router      	/comments/{postId} [get]
func (c *CommentController) GetAll(ctx *gin.Context) {
	params := QueryParams{}
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	postIdParam := ctx.Param("postId")
	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "ID bài viết không hợp lệ")
		return
	}
	params.PostId = postId
	comments, total, err := c.commentService.GetAll(params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Lấy danh sách bình luận thành công", gin.H{
		"items": comments,
		"total": total,
	})
}

// CreateComment 	godoc
// @Security		BearerAuth
// @Summary			Create a new comment
// @Description 	Create a new comment on a specific post
// @Tags        	Comment
// @Produce     	json
// @Param       	postId path string true "Post ID"
// @Param       	content body CreateCommentRequest true "Comment content"
// @Success     	200 {object} CreateCommentResponse
// @Router      	/comments/{postId} [post]
func (c *CommentController) CreateComment(ctx *gin.Context) {
	postIdParam := ctx.Param("postId")
	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "ID bài viết không hợp lệ")
		return
	}

	createCommentRequest := CreateCommentRequest{}
	if err := ctx.ShouldBindJSON(&createCommentRequest); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	createCommentRequest.PostId = postId
	createCommentRequest.UserId = userId.(uuid.UUID)

	result, err := c.commentService.CreateComment(&createCommentRequest)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Tạo bình luận thành công", result)
}

// UpdateComment 	godoc
// @Security 		BearerAuth
// @Summary			Update a comment
// @Description 	Update the content of an existing comment
// @Tags        	Comment
// @Produce     	json
// @Param       	id path string true "Comment ID"
// @Param       	content body UpdateCommentRequest true "New comment content"
// @Success     	200 {object} UpdateCommentResponse
// @Router      	/comments/{id} [patch]
func (c *CommentController) UpdateComment(ctx *gin.Context) {
	commentIdParam := ctx.Param("id")
	commentId, err := uuid.Parse(commentIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "ID comment không hợp lệ")
		return
	}
	updateCommentRequest := UpdateCommentRequest{}
	if err := ctx.ShouldBindJSON(&updateCommentRequest); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	updateCommentRequest.Id = commentId
	updateCommentRequest.UserId = userId.(uuid.UUID)

	result, err := c.commentService.UpdateComment(&updateCommentRequest)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Cập nhật bình luận thành công", result)

}

// DeleteComment 	godoc
// @Security 		BearerAuth
// @Summary			Delete a comment
// @Description 	Delete an existing comment by ID
// @Tags        	Comment
// @Produce     	json
// @Param       	id path string true "Comment ID"
// @Success     	200 {object} map[string]interface{}
// @Router      	/comments/{id} [delete]
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentIdParam := ctx.Param("id")
	commentId, err := uuid.Parse(commentIdParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "ID comment không hợp lệ")
		return
	}
	userId, _ := ctx.Get("userId")
	result, err := c.commentService.DeleteComment(commentId, userId.(uuid.UUID))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Xóa bình luận thành công", result)
}
