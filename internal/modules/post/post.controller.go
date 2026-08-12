package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostController struct {
	postService PostService
}

func NewPostController(postService PostService) *PostController {
	return &PostController{postService: postService}
}

// CreatePost	godoc
// @Security	BearerAuth
// @Summary		Create a new post
// @Description Create a new post with the given content and attachments
// @Tags        Post
// @Produce     json
// @Param       content body CreatePostRequest true "Post content"
// @Success     200 {object} CreatePostResponse
// @Router      /posts/ [post]
func (p *PostController) CreatePost(ctx *gin.Context) {

	createPostRequest := CreatePostRequest{}
	if err := ctx.ShouldBindJSON(&createPostRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	userId, _ := ctx.Get("userId")
	createPostRequest.UserId = userId.(uuid.UUID)

	result, err := p.postService.CreatePost(&createPostRequest)
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

// UpdatePost	godoc
// @Security	BearerAuth
// @Summary		Update an existing post
// @Description Update the content, attachments or privacy of an existing post
// @Tags        Post
// @Produce     json
// @Param 		id path string true "Post ID"
// @Param       content body UpdatePostRequest true "Post data to update"
// @Success     200 {object} map[string]interface{} "Returns the updated post data"
// @Router      /posts/{id} [patch]
func (p *PostController) UpdatePost(ctx *gin.Context) {
	postIdParam := ctx.Param("id")
	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID bài viết không hợp lệ"})
		return
	}
	updatePostRequest := UpdatePostRequest{}
	if err := ctx.ShouldBindJSON(&updatePostRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}
	updatePostRequest.Id = postId
	result, err := p.postService.UpdatePost(&updatePostRequest)
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

// DeletePost	godoc
// @Security	BearerAuth
// @Summary		Delete an existing post
// @Description Delete an existing post
// @Tags        Post
// @Produce     json
// @Param 		id path string true "Post ID"
// @Success     200 {object} map[string]interface{} "Returns the deleted post data"
// @Router      /posts/{id} [delete]
func (p *PostController) DeletePost(ctx *gin.Context) {
	postIdParam := ctx.Param("id")
	postId, err := uuid.Parse(postIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID bài viết không hợp lệ"})
		return
	}
	userId, _ := ctx.Get("userId")
	result, err := p.postService.DeletePost(postId, userId.(uuid.UUID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	if !result {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status": "failed",
			"error":  "Invalid permission or post not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete post successfully",
	})
}
