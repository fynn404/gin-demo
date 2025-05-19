package v1

import (
	controllers "github.com/fynn404/gin-demo/backend/internal/handler"
	"github.com/fynn404/gin-demo/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	handler *controllers.Handler
}

func NewUserRoutes(h *controllers.Handler) *UserRoutes {
	return &UserRoutes{handler: h}
}

func (r *UserRoutes) Register(group *gin.RouterGroup) {
	auth := group.Group("/users", middleware.AuthMiddleware())
	{
		auth.GET("/profile", r.handler.User.GetUserProfile)
		auth.PUT("/profile", r.handler.User.UpdateUserProfile)
		auth.POST("/avatar", r.handler.User.UploadAvatar)
	}
}
