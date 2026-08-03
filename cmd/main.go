package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Khởi tạo router mặc định của Gin
	r := gin.Default()

	// Khai báo đường dẫn GET /ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Chạy server ở cổng 8080 (mặc định)
	r.Run(":8080")
}
