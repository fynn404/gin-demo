package v1

import (
	controllers "github.com/fynn404/gin-demo/backend/internal/handler"
	"github.com/fynn404/gin-demo/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type TodoRoutes struct {
	handler *controllers.Handler
}

func NewTodoRoutes(h *controllers.Handler) *TodoRoutes {
	return &TodoRoutes{handler: h}
}

func (t *TodoRoutes) Register(group *gin.RouterGroup) {
	todos := group.Group("/todos", middleware.AuthMiddleware())
	{
		todos.GET("", t.handler.Todo.ListUserTodos)                     // 获取待办事项列表
		todos.POST("", t.handler.Todo.CreateTodoItem)                   // 创建待办事项
		todos.GET("/:id", t.handler.Todo.GetTodoItemDetail)             // 获取待办事项详情
		todos.PUT("/:id", t.handler.Todo.UpdateTodoItem)                // 更新待办事项
		todos.DELETE("/:id", t.handler.Todo.DeleteTodoItem)             // 删除待办事项
		todos.PATCH("/:id/status", t.handler.Todo.ChangeTodoItemStatus) // 更改待办事项状态
	}
}
