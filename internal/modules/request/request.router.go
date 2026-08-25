package request

import (
	"hub/internal/modules/contact"
	"hub/internal/modules/notification"
	"hub/internal/websocket"
	"hub/pkg/middleware"
	"hub/pkg/transaction"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequestRouter(r *gin.Engine, db *gorm.DB, hub *websocket.Hub) {
	txManager := transaction.NewTransactionManager(db)
	notificationRepo := notification.NewNotificationRepositoryImp(db)
	contactRepo := contact.NewContactRepositoryImp(db)

	requestRepo := NewRequestRepositoryImp(db)
	requestService := NewRequestServiceImp(requestRepo, notificationRepo, contactRepo, txManager)
	requestController := NewRequestController(requestService, hub)

	requestRouter := r.Group("/requests")
	requestRouter.Use(middleware.AuthMiddleware())
	{
		requestRouter.GET("/", requestController.GetAll)
		requestRouter.POST("/", requestController.SendRequest)
		requestRouter.PATCH("/:id", requestController.UpdateRequest)
		requestRouter.DELETE("/:id", requestController.DeleteRequest)
	}
}
