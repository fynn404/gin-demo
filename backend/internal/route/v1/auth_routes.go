package v1

import (
	controllers "github.com/fynn404/gin-demo/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

// AuthRoutes 认证相关路由
type AuthRoutes struct {
	handler *controllers.Handler
}

func NewAuthRoutes(h *controllers.Handler) *AuthRoutes {
	return &AuthRoutes{handler: h}
}

func (r *AuthRoutes) Register(group *gin.RouterGroup) {
	auth := group.Group("/auth")
	{
		auth.POST("/login", r.handler.Auth.Login)
		auth.POST("/register", r.handler.Auth.Register)
		//auth.POST("/logout", r.handler.Auth.)
	}
}
