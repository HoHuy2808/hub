package response

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Success trả về một response thành công
func Success(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, Response{
		Code:    code,
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// Error trả về một response lỗi
func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, Response{
		Code:    code,
		Status:  "failed",
		Message: message,
		Data:    nil,
	})
}
