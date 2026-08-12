package middleware

import (
	"net/http"
	"strings"

	"hub/pkg/authenticate"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"error":  "Authorization header is required",
			})
			ctx.Abort()
			return
		}
		var tokenString string
		if strings.Contains(authHeader, " ") {
			tokenString = strings.Split(authHeader, " ")[1]
		} else {
			tokenString = authHeader
		}

		claims, err := authenticate.ValidateToken(tokenString)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"error":  "Invalid token: " + err.Error(),
			})
			return
		}

		ctx.Set("userId", claims.UserId)
		ctx.Next()
	}
}
