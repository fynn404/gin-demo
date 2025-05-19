// backend/internal/handler/auth_handler.go
package handler

import (
	"github.com/fynn404/gin-demo/backend/internal/service"
	"strings"
	"time"

	"github.com/fynn404/gin-demo/backend/internal/common/config"
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
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从请求头获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		Error(c, 401, "No token provided", nil)
		return
	}

	// 提取 token（去掉 "Bearer " 前缀）
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 将 token 加入 Redis 黑名单
	// 设置过期时间为 token 的剩余有效期，这里假设为24小时
	err := config.AddTokenToBlacklist(token, 24*time.Hour)
	if err != nil {
		InternalServerError(c, "Failed to logout", err)
		return
	}
	Success(c, nil)

}
