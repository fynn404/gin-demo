package service

import (
	"fmt"
	"time"

	"github.com/fynn404/gin-demo/backend/internal/common/config"
	"github.com/fynn404/gin-demo/backend/internal/model"
	"github.com/fynn404/gin-demo/backend/internal/repository"
	"github.com/fynn404/gin-demo/backend/pkg"
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

// TodoService defines the methods available for managing todos.
type TodoService interface {
	ListUserTodos(c *gin.Context, req *ListUserTodosReq) (*ListUserTodosResp, error)
	CreateTodoItem(c *gin.Context, req *CreateTodoItemReq) (*CreateTodoItemResp, error)
	GetTodoItemDetail(c *gin.Context, req *GetTodoItemDetailReq) (*GetTodoItemDetailResp, error)
	UpdateTodoItem(c *gin.Context, req *UpdateTodoItemReq) (*UpdateTodoItemResp, error)
	DeleteTodoItem(c *gin.Context, req *DeleteTodoItemReq) (*DeleteTodoItemResp, error)
	ChangeTodoItemStatus(c *gin.Context, req *ChangeTodoItemStatusReq) (*ChangeTodoItemStatusResp, error)
}

// Todo represents the service for managing todos.
type Todo struct {
	todoRepo repository.TodoRepository
}

func (t *Todo) ListUserTodos(c *gin.Context, req *ListUserTodosReq) (*ListUserTodosResp, error) {
	// 将字符串类型的 UserID 转换为 uint64
	userID := pkg.ConvertStr2Int(req.UserID)

	// 使用事务来确保数据一致性
	var todos []model.TodoItem
	var total int64
	var err error

	err = config.Transaction(func(session *xorm.Session) error {
		// 获取用户的待办事项列表和总数
		todos, total, err = t.todoRepo.ListByUserID(uint64(userID), req.Page, req.Size)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 转换数据模型到响应结构
	todoItems := make([]*TodoItem, 0, len(todos))
	for _, todo := range todos {
		todoItems = append(todoItems, &TodoItem{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
			Priority:    todo.Priority,
			DueDate:     todo.DueDate,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		})
	}

	// 返回响应
	return &ListUserTodosResp{
		Todos: todoItems,
		Total: total,
	}, nil
}

// CreateTodoItem 创建待办事项
func (t *Todo) CreateTodoItem(c *gin.Context, req *CreateTodoItemReq) (*CreateTodoItemResp, error) {
	userID := pkg.ConvertStr2Int(req.UserID)
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}
	var (
		priority = "medium"
		dueDate  = time.Now().Add(24 * time.Hour)
	)
	var todo *model.TodoItem
	err := config.Transaction(func(session *xorm.Session) error {
		// 检查用户是否存在
		exists, err := session.Table("users_tab").ID(userID).Exist()
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("user not found")
		}
		if req.priority != "" {
			priority = req.priority
		}
		if req.DueDate != time.Unix(0, 0) {
			dueDate = req.DueDate
		}

		// 创建新的待办事项
		todo = &model.TodoItem{
			UserId:      uint64(userID),
			Title:       req.Title,
			Description: req.Description,
			Priority:    priority, // 默认优先级
			DueDate:     dueDate,  // 默认截止时间为24小时后
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := t.todoRepo.Create(todo); err != nil {
			return fmt.Errorf("failed to create todo item: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateTodoItemResp{
		Todo: &TodoItem{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
			Priority:    todo.Priority,
			DueDate:     todo.DueDate,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}

// GetTodoItemDetail 获取待办事项详情
func (t *Todo) GetTodoItemDetail(c *gin.Context, req *GetTodoItemDetailReq) (*GetTodoItemDetailResp, error) {
	todoID := pkg.ConvertStr2Int(req.TodoID)
	if todoID <= 0 {
		return nil, fmt.Errorf("invalid todo ID")
	}

	var todo *model.TodoItem
	err := config.Transaction(func(session *xorm.Session) error {
		var err error
		todo, err = t.todoRepo.GetByID(uint64(todoID))
		if err != nil {
			return fmt.Errorf("failed to get todo item: %v", err)
		}
		if todo == nil {
			return fmt.Errorf("todo item not found")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &GetTodoItemDetailResp{
		Todo: &TodoItem{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
			Priority:    todo.Priority,
			DueDate:     todo.DueDate,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}

// UpdateTodoItem 更新待办事项
func (t *Todo) UpdateTodoItem(c *gin.Context, req *UpdateTodoItemReq) (*UpdateTodoItemResp, error) {
	todoID := pkg.ConvertStr2Int(req.TodoID)
	if todoID <= 0 {
		return nil, fmt.Errorf("invalid todo ID")
	}

	var todo *model.TodoItem
	err := config.Transaction(func(session *xorm.Session) error {
		// 先获取现有的待办事项
		var err error
		todo, err = t.todoRepo.GetByID(uint64(todoID))
		if err != nil {
			return fmt.Errorf("failed to get todo item: %v", err)
		}
		if todo == nil {
			return fmt.Errorf("todo item not found")
		}

		// 更新字段
		if req.Title != "" {
			todo.Title = req.Title
		}
		if req.Description != "" {
			todo.Description = req.Description
		}
		if req.Priority != "" {
			todo.Priority = req.Priority
		}
		todo.UpdatedAt = time.Now()

		// 保存更新
		if err := t.todoRepo.Update(todo); err != nil {
			return fmt.Errorf("failed to update todo item: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &UpdateTodoItemResp{
		Todo: &TodoItem{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
			Priority:    todo.Priority,
			DueDate:     todo.DueDate,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}

// DeleteTodoItem 删除待办事项
func (t *Todo) DeleteTodoItem(c *gin.Context, req *DeleteTodoItemReq) (*DeleteTodoItemResp, error) {
	todoID := pkg.ConvertStr2Int(req.TodoID)
	if todoID <= 0 {
		return nil, fmt.Errorf("invalid todo ID")
	}

	err := config.Transaction(func(session *xorm.Session) error {
		// 检查待办事项是否存在
		exists, err := session.Table("todo_items_tab").ID(todoID).Exist()
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("todo item not found")
		}

		// 删除待办事项
		if err := t.todoRepo.Delete(uint64(todoID)); err != nil {
			return fmt.Errorf("failed to delete todo item: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteTodoItemResp{
		Success: true,
	}, nil
}

// ChangeTodoItemStatus 更改待办事项状态
func (t *Todo) ChangeTodoItemStatus(c *gin.Context, req *ChangeTodoItemStatusReq) (*ChangeTodoItemStatusResp, error) {
	todoID := pkg.ConvertStr2Int(req.TodoID)
	if todoID <= 0 {
		return nil, fmt.Errorf("invalid todo ID")
	}

	var todo *model.TodoItem
	err := config.Transaction(func(session *xorm.Session) error {
		// 获取现有的待办事项
		var err error
		todo, err = t.todoRepo.GetByID(uint64(todoID))
		if err != nil {
			return fmt.Errorf("failed to get todo item: %v", err)
		}
		if todo == nil {
			return fmt.Errorf("todo item not found")
		}

		// 更新状态
		todo.Completed = *req.Completed
		todo.UpdatedAt = time.Now()
		fmt.Println(todo)
		// 保存更新
		if err := t.todoRepo.Update(todo); err != nil {
			return fmt.Errorf("failed to update todo status: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ChangeTodoItemStatusResp{
		Todo: TodoItem{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
			Priority:    todo.Priority,
			DueDate:     todo.DueDate,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		},
	}, nil
}

func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &Todo{todoRepo: todoRepo}
}

// ListUserTodosReq represents the request structure for listing user todos.
type ListUserTodosReq struct {
	UserID string `json:"userId"`
	Page   int    `json:"page" binding:"required"`
	Size   int    `json:"size" binding:"required"`
}

// ListUserTodosResp represents the response structure for listing user todos.
type ListUserTodosResp struct {
	Todos []*TodoItem `json:"todos"`
	Total int64       `json:"total"`
}

// CreateTodoItemReq represents the request structure for creating a todo item.
type CreateTodoItemReq struct {
	UserID      string    `json:"user_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
}

// CreateTodoItemResp represents the response structure for creating a todo item.
type CreateTodoItemResp struct {
	Todo *TodoItem `json:"todo"`
}

// GetTodoItemDetailReq represents the request structure for getting a todo item detail.
type GetTodoItemDetailReq struct {
	TodoID string `json:"todoId" binding:"required"`
}

// GetTodoItemDetailResp represents the response structure for getting a todo item detail.
type GetTodoItemDetailResp struct {
	Todo *TodoItem `json:"todo"`
}

// UpdateTodoItemReq represents the request structure for updating a todo item.
type UpdateTodoItemReq struct {
	TodoID      string `json:"todoId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// UpdateTodoItemResp represents the response structure for updating a todo item.
type UpdateTodoItemResp struct {
	Todo *TodoItem `json:"todo"`
}

// DeleteTodoItemReq represents the request structure for deleting a todo item.
type DeleteTodoItemReq struct {
	TodoID string `json:"todoId" binding:"required"`
}

// DeleteTodoItemResp represents the response structure for deleting a todo item.
type DeleteTodoItemResp struct {
	Success bool `json:"success"`
}

// ChangeTodoItemStatusReq represents the request structure for changing a todo item status.
type ChangeTodoItemStatusReq struct {
	TodoID    string `json:"todoId"`
	Completed *bool  `json:"completed" binding:"required"`
}

// ChangeTodoItemStatusResp represents the response structure for changing a todo item status.
type ChangeTodoItemStatusResp struct {
	Todo TodoItem `json:"todo"`
}

// TodoItem represents a single todo item.
type TodoItem struct {
	Id          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"dueDate"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
