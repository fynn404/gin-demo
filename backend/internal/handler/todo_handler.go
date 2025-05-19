package handler

import (
	"github.com/fynn404/gin-demo/backend/internal/service"
	"github.com/fynn404/gin-demo/backend/pkg"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// ListUserTodos 获取用户的待办事项列表
func (h *TodoHandler) ListUserTodos(c *gin.Context) {
	// 从请求中获取参数
	req := &service.ListUserTodosReq{
		Page: int(pkg.ConvertStr2Int(c.DefaultQuery("page", "1"))),
		Size: int(pkg.ConvertStr2Int(c.DefaultQuery("size", "10"))),
	}

	// 从上下文获取用户ID（假设已经通过认证中间件设置）
	userID := c.GetString("user_id")
	if userID == "" {
		Unauthorized(c, "User not authenticated")
		return
	}
	req.UserID = userID

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	// 调用服务层方法
	resp, err := h.todoService.ListUserTodos(c, req)
	if err != nil {
		InternalServerError(c, "Failed to get todo list", err)
		return
	}

	Success(c, resp)
}

// CreateTodoItem 创建待办事项
func (h *TodoHandler) CreateTodoItem(c *gin.Context) {
	var req service.CreateTodoItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body")
		return
	}

	// 从上下文获取用户ID（假设已经通过认证中间件设置）
	userID := c.GetString("user_id")
	if userID == "" {
		Unauthorized(c, "User not authenticated")
		return
	}
	req.UserID = userID

	resp, err := h.todoService.CreateTodoItem(c, &req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			NotFound(c, "User not found")
		case "invalid user ID":
			BadRequest(c, "Invalid user ID")
		default:
			InternalServerError(c, "Failed to create todo item", err)
		}
		return
	}

	Success(c, resp)
}

// GetTodoItemDetail 获取待办事项详情
func (h *TodoHandler) GetTodoItemDetail(c *gin.Context) {
	todoID := c.Param("id")
	if todoID == "" {
		BadRequest(c, "Todo ID is required")
		return
	}

	req := &service.GetTodoItemDetailReq{
		TodoID: todoID,
	}

	resp, err := h.todoService.GetTodoItemDetail(c, req)
	if err != nil {
		switch err.Error() {
		case "todo item not found":
			NotFound(c, "Todo item not found")
		case "invalid todo ID":
			BadRequest(c, "Invalid todo ID")
		default:
			InternalServerError(c, "Failed to get todo item", err)
		}
		return
	}

	Success(c, resp)
}

// UpdateTodoItem 更新待办事项
func (h *TodoHandler) UpdateTodoItem(c *gin.Context) {
	var req service.UpdateTodoItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body")
		return
	}

	// 从路径参数获取待办事项ID
	req.TodoID = c.Param("id")
	if req.TodoID == "" {
		BadRequest(c, "Todo ID is required")
		return
	}

	resp, err := h.todoService.UpdateTodoItem(c, &req)
	if err != nil {
		switch err.Error() {
		case "todo item not found":
			NotFound(c, "Todo item not found")
		case "invalid todo ID":
			BadRequest(c, "Invalid todo ID")
		default:
			InternalServerError(c, "Failed to update todo item", err)
		}
		return
	}

	Success(c, resp)
}

// DeleteTodoItem 删除待办事项
func (h *TodoHandler) DeleteTodoItem(c *gin.Context) {
	todoID := c.Param("id")
	if todoID == "" {
		BadRequest(c, "Todo ID is required")
		return
	}

	req := &service.DeleteTodoItemReq{
		TodoID: todoID,
	}

	resp, err := h.todoService.DeleteTodoItem(c, req)
	if err != nil {
		switch err.Error() {
		case "todo item not found":
			NotFound(c, "Todo item not found")
		case "invalid todo ID":
			BadRequest(c, "Invalid todo ID")
		default:
			InternalServerError(c, "Failed to delete todo item", err)
		}
		return
	}

	Success(c, resp)
}

// ChangeTodoItemStatus 更改待办事项状态
func (h *TodoHandler) ChangeTodoItemStatus(c *gin.Context) {
	var req service.ChangeTodoItemStatusReq

	if err := c.ShouldBindJSON(&req); err != nil {

		BadRequest(c, "Invalid request body")
		return
	}

	// 从路径参数获取待办事项ID
	req.TodoID = c.Param("id")
	if req.TodoID == "" {
		BadRequest(c, "Todo ID is required")
		return
	}

	resp, err := h.todoService.ChangeTodoItemStatus(c, &req)
	if err != nil {
		switch err.Error() {
		case "todo item not found":
			NotFound(c, "Todo item not found")
		case "invalid todo ID":
			BadRequest(c, "Invalid todo ID")
		default:
			InternalServerError(c, "Failed to change todo item status", err)
		}
		return
	}

	Success(c, resp)
}
