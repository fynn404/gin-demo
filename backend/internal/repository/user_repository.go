package repository

import (
	"errors"

	"github.com/fynn404/gin-demo/backend/internal/model"
	"xorm.io/xorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint64) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint64) error
	List(page, pageSize int) ([]model.User, int64, error)
	GetUserWithTodos(id uint64) (*model.User, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	engine *xorm.Engine
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(engine *xorm.Engine) UserRepository {
	return &userRepository{engine: engine}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	_, err := r.engine.Insert(user)
	return err
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint64) (*model.User, error) {
	user := &model.User{}
	has, err := r.engine.ID(id).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	has, err := r.engine.Where("username = ?", username).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	_, err := r.engine.ID(user.Id).Update(user)
	return err
}

// Delete 删除用户
func (r *userRepository) Delete(id uint64) error {
	_, err := r.engine.ID(id).Delete(&model.User{})
	return err
}

// List 分页获取用户列表
func (r *userRepository) List(page, pageSize int) ([]model.User, int64, error) {
	start := (page - 1) * pageSize
	var users []model.User

	total, err := r.engine.Limit(pageSize, start).FindAndCount(&users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserWithTodos 获取用户及其所有待办事项
func (r *userRepository) GetUserWithTodos(id uint64) (*model.User, error) {
	user := &model.User{}
	has, err := r.engine.ID(id).Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}

	// 获取用户的所有待办事项
	if err := r.engine.Where("user_id = ?", id).Find(&user.TodoItems); err != nil {
		return nil, err
	}

	return user, nil
}
