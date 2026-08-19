package post

import (
	"net/http"

	"hub/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostController struct {
	postService PostService
}

func NewPostController(postService PostService) *PostController {
	return &PostController{postService: postService}
}

// GetAll 		godoc
// @Security	BearerAuth
// @Summary		Get all posts
// @Description Get all posts with pagination
// @Tags        Post
// @Produce     json
// @Param       page query int false "Page number"
// @Param       limit query int false "Number of items per page"
// @Param       search query string false "Search term"
// @Param       created_at query string false "Sort by"
// @Param       desc query string false "Sort order"
// @Success     200 {object} PaginatedResponse "Returns the list of posts"
// @Router      /posts/ [get]
func (p *PostController) GetAll(ctx *gin.Context) {
	var query QueryParams
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	posts, total, err := p.postService.GetAll(query)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Lấy danh sách bài viết thành công", gin.H{
		"items": posts,
		"total": total,
	})
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
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, _ := ctx.Get("userId")
	createPostRequest.UserId = userId.(uuid.UUID)

	result, err := p.postService.CreatePost(&createPostRequest)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Tạo bài viết thành công", result)
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
		response.Error(ctx, http.StatusBadRequest, "ID bài viết không hợp lệ")
		return
	}
	updatePostRequest := UpdatePostRequest{}
	if err := ctx.ShouldBindJSON(&updatePostRequest); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	updatePostRequest.Id = postId
	result, err := p.postService.UpdatePost(&updatePostRequest)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Cập nhật bài viết thành công", result)
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
		response.Error(ctx, http.StatusBadRequest, "ID bài viết không hợp lệ")
		return
	}
	userId, _ := ctx.Get("userId")
	result, err := p.postService.DeletePost(postId, userId.(uuid.UUID))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if !result {
		response.Error(ctx, http.StatusForbidden, "Invalid permission or post not found")
		return
	}

	response.Success(ctx, http.StatusOK, "Xóa bài viết thành công", nil)
}
