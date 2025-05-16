// backend/internal/handler/handler.go
package handler

import (
	"github.com/fynn404/gin-demo/backend/internal/repository"
	"github.com/fynn404/gin-demo/backend/internal/service"
)

// Handler 包含所有HTTP处理器
type Handler struct {
	Auth *AuthHandler
	User *UserHandler
	Todo *TodoHandler
}

// HandlerConfig 处理器配置
type HandlerConfig struct {
	UserRepo repository.UserRepository
	TodoRepo repository.TodoRepository
}

// NewHandler 创建一个新的Handler实例
func NewHandler(cfg *HandlerConfig) *Handler {
	// 初始化 services
	authService := service.NewAuthService(cfg.UserRepo)
	userService := service.NewUserService(cfg.UserRepo)
	todoService := service.NewTodoService(cfg.TodoRepo)

	// 初始化 handlers
	return &Handler{
		Auth: NewAuthHandler(authService),
		User: NewUserHandler(userService),
		Todo: NewTodoHandler(todoService),
	}
}
