package repository

import (
	"errors"

	"github.com/fynn404/gin-demo/backend/internal/model"
	"xorm.io/xorm"
)

// TodoRepository 待办事项仓储接口
type TodoRepository interface {
	Create(todo *model.TodoItem) error
	GetByID(id uint64) (*model.TodoItem, error)
	Update(todo *model.TodoItem) error
	Delete(id uint64) error
	List(page, pageSize int) ([]model.TodoItem, int64, error)
	ListByUserID(userID uint64, page, pageSize int) ([]model.TodoItem, int64, error)
	GetTodoWithUser(id uint64) (*model.TodoItem, error)
}

// todoRepository 待办事项仓储实现
type todoRepository struct {
	engine *xorm.Engine
}

// NewTodoRepository 创建待办事项仓储实例
func NewTodoRepository(engine *xorm.Engine) TodoRepository {
	return &todoRepository{engine: engine}
}

// Create 创建待办事项
func (r *todoRepository) Create(todo *model.TodoItem) error {
	_, err := r.engine.Insert(todo)
	return err
}

// GetByID 根据ID获取待办事项
func (r *todoRepository) GetByID(id uint64) (*model.TodoItem, error) {
	todo := &model.TodoItem{}
	has, err := r.engine.ID(id).Get(todo)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("todo item not found")
	}
	return todo, nil
}

// Update 更新待办事项
func (r *todoRepository) Update(todo *model.TodoItem) error {
	_, err := r.engine.ID(todo.Id).Cols("completed", "updated_at", "title", "description", "priority", "due_date").Update(todo)
	return err
}

// Delete 删除待办事项
func (r *todoRepository) Delete(id uint64) error {
	_, err := r.engine.ID(id).Delete(&model.TodoItem{})
	return err
}

// List 分页获取待办事项列表
func (r *todoRepository) List(page, pageSize int) ([]model.TodoItem, int64, error) {
	start := (page - 1) * pageSize
	var todos []model.TodoItem

	total, err := r.engine.Limit(pageSize, start).FindAndCount(&todos)
	if err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

// ListByUserID 获取指定用户的待办事项列表
func (r *todoRepository) ListByUserID(userID uint64, page, pageSize int) ([]model.TodoItem, int64, error) {
	start := (page - 1) * pageSize
	var todos []model.TodoItem

	total, err := r.engine.Where("user_id = ?", userID).Limit(pageSize, start).FindAndCount(&todos)
	if err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

// GetTodoWithUser 获取待办事项及其关联的用户信息
func (r *todoRepository) GetTodoWithUser(id uint64) (*model.TodoItem, error) {
	todo := &model.TodoItem{}
	has, err := r.engine.ID(id).Get(todo)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("todo item not found")
	}

	// 获取关联的用户信息
	if todo.UserId > 0 {
		todo.User = &model.User{}
		if _, err := r.engine.ID(todo.UserId).Get(todo.User); err != nil {
			return nil, err
		}
	}

	return todo, nil
}
