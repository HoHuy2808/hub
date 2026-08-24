package reaction

import (
	"hub/internal/modules/notification"
	"hub/internal/modules/post"
	"hub/internal/websocket"
	"hub/pkg/middleware"
	"hub/pkg/transaction"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReactRouter(r *gin.Engine, db *gorm.DB, hub *websocket.Hub) {
	postRepo := post.NewPostRepositoryImp(db)
	notificationRepo := notification.NewNotificationRepositoryImp(db)
	txManager := transaction.NewTransactionManager(db)

	reactionRepo := NewReactionRepositoryImp(db)
	reactionService := NewReactionServiceImp(reactionRepo, postRepo, notificationRepo, txManager)
	reactionController := NewReactionController(reactionService, hub)

	reactionRouter := r.Group("/reactions")
	reactionRouter.Use(middleware.AuthMiddleware())
	{
		reactionRouter.GET("/:postId", reactionController.GetAll)
		reactionRouter.POST("/:postId", reactionController.React)
		reactionRouter.PATCH("/:postId", reactionController.UpdateReact)
		reactionRouter.DELETE("/:postId", reactionController.Unreact)
	}
}
