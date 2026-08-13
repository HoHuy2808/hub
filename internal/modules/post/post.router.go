package post

import (
	"hub/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostRouter(r *gin.Engine, db *gorm.DB) {
	postRepo := NewPostRepositoryImp(db)
	postService := NewPostServiceImp(postRepo)
	postController := NewPostController(postService)

	postRoutes := r.Group("/posts")
	postRoutes.Use(middleware.AuthMiddleware())
	{
		postRoutes.GET("/", postController.GetAll)
		postRoutes.POST("/", postController.CreatePost)
		postRoutes.PATCH("/:id", postController.UpdatePost)
		postRoutes.DELETE("/:id", postController.DeletePost)
	}
}
