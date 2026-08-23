package main

import (
	"fmt"
	"net/http"
	"time"

	"hub/internal/database"
	"hub/internal/modules/attachment"
	"hub/internal/modules/auth"
	"hub/internal/modules/comment"
	"hub/internal/modules/post"
	"hub/internal/modules/reaction"

	"hub/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "hub/docs" // Import thư mục docs sinh ra bởi swag

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           			Hub API
// @version         			1.0
// @description     			Đây là API cho dự án Backend Hub.
// @host            			localhost:2808
// @BasePath        			/
// @securityDefinitions.apikey 	BearerAuth
// @in 							header
// @name 						Authorization
func main() {
	// Load biến môi trường từ file .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ Không tìm thấy file .env")
	}

	// Kiểm tra kết nối database
	_, err := database.ConnectToPostgreSQL()
	if err != nil {
		fmt.Println("❌ Lỗi kết nối database:", err)
	} else {
		fmt.Println("✅ Kết nối PostgreSQL thành công bằng GORM!")
	}

	// Khởi tạo router mặc định của Gin
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Khởi tạo Hub
	hub := websocket.NewHub()
	// Chạy goroutine quản lý hub
	go hub.Run()

	// Khai báo đường dẫn kết nối WebSocket
	r.GET("/ws", websocket.ServeWs(hub))

	// API Health Check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"message":   "Backend server is running healthy 🚀",
			"timestamp": time.Now().Unix(),
		})
	})

	// Cấu hình Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routers
	auth.AuthRouter(r, database.DB)
	post.PostRouter(r, database.DB)
	attachment.AttachmentRouter(r, database.DB)
	comment.CommentRouter(r, database.DB, hub)
	reaction.ReactRouter(r, database.DB, hub)

	// Chạy server ở cổng 2808
	r.Run(":2808")
}
