// backend/internal/handler/auth_handler.go
package handler

import (
	"net/http"

	"github.com/fynn404/gin-demo/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 处理用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body")
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			NotFound(c, "User not found")
		case service.ErrInvalidCredentials:
			Unauthorized(c, "Invalid credentials")
		default:
			InternalServerError(c, "Failed to login", err)
		}
		return
	}

	Success(c, resp)
}

// Register 处理用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body")
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		switch err {
		case service.ErrUserExists:
			Error(c, 409, "User already exists", err)
		default:
			InternalServerError(c, "Failed to register user", err)
		}
		return
	}

	Success(c, resp)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the token cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Successfully logged out",
	})
}
