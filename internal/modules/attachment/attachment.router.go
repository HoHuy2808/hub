package attachment

import (
	"hub/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AttachmentRouter(r *gin.Engine, db *gorm.DB) {
	attachRepo := NewAttachmentRepositoryImp(db)
	attachService := NewAttachmentServiceImp(attachRepo)
	attachController := NewAttachmentController(attachService)

	attachmentRoutes := r.Group("/posts")
	attachmentRoutes.Use(middleware.AuthMiddleware())
	{
		attachmentRoutes.GET("/:id/attachments", attachController.GetAll)
		attachmentRoutes.POST("/:id/attachments", attachController.AddAttachment)
		attachmentRoutes.DELETE("/attachments/:id", attachController.DeleteAttachment)
	}

	// Route cho user (Tách riêng ra để không bị đụng độ với /posts/:post_id)
	userAttachmentRoutes := r.Group("/users")
	userAttachmentRoutes.Use(middleware.AuthMiddleware())
	{
		userAttachmentRoutes.GET("/me/attachments", attachController.GetAllByUser)
	}

}
