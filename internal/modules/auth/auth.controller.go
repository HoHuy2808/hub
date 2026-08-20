package auth

import (
	// "hub/internal/modules/auth/service"
	"net/http"

	"hub/pkg/response"

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
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	result, err := u.userService.Register(&registerInput)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Đăng ký thành công", result)
}

// Login 		godoc
// @Summary		Login
// @Description	Login to the application
// @Tags		Auth
// @Produce		json
// @Param		user body LoginRequest true "User object"
// @Success		200 {object} LoginResponse
// @Router		/auth/login [post]
func (u *UserController) Login(ctx *gin.Context) {

	loginInput := LoginRequest{}

	if err := ctx.ShouldBindJSON(&loginInput); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	result, err := u.userService.Login(&loginInput)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "Đăng nhập thành công", result)
}
