package contact

import (
	"hub/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ContactRouter(r *gin.Engine, db *gorm.DB) {
	contactRepo := NewContactRepositoryImp(db)
	contactService := NewContactServiceImp(contactRepo)
	contactController := NewContactController(contactService)

	contactRouter := r.Group("/contacts")
	contactRouter.Use(middleware.AuthMiddleware())
	{
		contactRouter.GET("/", contactController.GetAll)
		contactRouter.PATCH("/:id", contactController.UpdateContact)
		contactRouter.DELETE("/:id", contactController.DeleteContact)
	}
}
