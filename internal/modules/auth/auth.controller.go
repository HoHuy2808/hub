package auth

import (
	// "hub/internal/modules/auth/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService UserService
}

func NewUserController(service UserService) *UserController {
	return &UserController{userService: service}
}

// Register 	godoc
// @Summary		Register User
// @Description	Register a new user
// @Tags		Auth
// @Produce		json
// @Param		user body RegisterRequest true "User object"
// @Success		200 {object} RegisterResponse
// @Router		/auth/register [post]
func (u *UserController) Register(ctx *gin.Context) {

	registerInput := RegisterRequest{}

	if err := ctx.ShouldBindJSON(&registerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Invalid Input",
			"error":   err.Error(),
		})
		return
	}

	result, _ := u.userService.Register(&registerInput)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}
