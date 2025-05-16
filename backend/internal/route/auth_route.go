package route

import (
	"github.com/fynn404/gin-demo/backend/internal/handler"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler *handler.AuthHandler
}

func NewRouter(authHandler *handler.AuthHandler) *Router {
	return &Router{
		authHandler: authHandler,
	}
}

func (r *Router) InitAuthRouter(Router *gin.RouterGroup) {
	authRouter := Router.Group("auth")
	{
		authRouter.POST("/login", r.authHandler.Login)
		authRouter.POST("/logout", r.authHandler.Logout)
	}
}
