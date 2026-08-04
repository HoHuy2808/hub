package main

import (
	"fmt"
	"net/http"

	"hub/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "hub/docs" // Import thư mục docs sinh ra bởi swag

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Hub API
// @version         1.0
// @description     Đây là API cho dự án Backend Hub.
// @host            localhost:8080
// @BasePath        /
func main() {
	// Load biến môi trường từ file .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ Không tìm thấy file .env, sẽ sử dụng biến môi trường hệ thống")
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

	// Khai báo đường dẫn GET /ping
	// @Summary      Kiểm tra API
	// @Description  Trả về một thông báo pong để xác nhận server đang chạy
	// @Tags         Health
	// @Produce      json
	// @Success      200  {object}  map[string]string
	// @Router       /ping [get]
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Cấu hình Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Chạy server ở cổng 8080 (mặc định)
	r.Run(":8080")
}
