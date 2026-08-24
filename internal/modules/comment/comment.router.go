package comment

import (
	"hub/internal/modules/notification"
	"hub/internal/modules/post"
	"hub/pkg/middleware"

	"hub/internal/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CommentRouter(r *gin.Engine, db *gorm.DB, hub *websocket.Hub) {
	commentRepo := NewCommentRepositoryImp(db)
	commentService := NewCommentServiceImp(commentRepo, post.NewPostRepositoryImp(db), notification.NewNotificationRepositoryImp(db))
	commentController := NewCommentController(commentService, hub)

	commentRouter := r.Group("/comments")
	commentRouter.Use(middleware.AuthMiddleware())
	{
		commentRouter.GET("/:postId", commentController.GetAll)
		commentRouter.POST("/:postId", commentController.CreateComment)
		commentRouter.PATCH("/:id", commentController.UpdateComment)
		commentRouter.DELETE("/:id", commentController.DeleteComment)
	}
}
