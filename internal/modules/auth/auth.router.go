package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRouter(router *gin.Engine, db *gorm.DB) {
	// 1. Setup Dependency Injection
	repo := NewUserRepositoryImp(db)
	service := NewUserServiceImp(repo)
	controller := NewUserController(service)

	// 2. Setup Router
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", controller.Register)
	}
}
