package v1

import (
	controllers "github.com/fynn404/gin-demo/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

// RouteGroup 定义路由组接口
type RouteGroup interface {
	Register(group *gin.RouterGroup)
}

// Routes 路由分组结构体
type Routes struct {
	auth *AuthRoutes
	user *UserRoutes
}

// NewRoutes 创建路由实例
func NewRoutes(h *controllers.Handler) *Routes {
	return &Routes{
		auth: NewAuthRoutes(h),
		user: NewUserRoutes(h),
	}
}

// SetupRoutes 配置所有路由
func SetupRoutes(r *gin.Engine, h *controllers.Handler) {
	// API v1 版本分组
	v1 := r.Group("/api/v1")

	// 创建路由实例
	routes := NewRoutes(h)

	// 注册各个模块的路由
	routes.auth.Register(v1)
	routes.user.Register(v1)

}
