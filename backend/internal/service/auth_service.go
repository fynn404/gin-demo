package service

import (
	"errors"
	"fmt"
	"github.com/fynn404/gin-demo/backend/internal/repository"
	"strconv"
	"time"

	"github.com/fynn404/gin-demo/backend/internal/model"
	"github.com/fynn404/gin-demo/backend/pkg"
	"golang.org/x/crypto/bcrypt"
)

// 错误定义
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserIdInvalid      = errors.New("user ID invalid")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type RegisterResp struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type AuthService interface {
	Login(req *LoginReq) (*LoginResp, error)
	Register(req *RegisterReq) (*RegisterResp, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Login 处理用户登录
func (s *authService) Login(req *LoginReq) (*LoginResp, error) {
	// 查找用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 验证密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成 token
	token, err := pkg.GenerateToken(&pkg.GenerateTokenReq{
		UserID:   fmt.Sprintf(strconv.FormatUint(user.Id, 10)),
		Role:     user.Role,
		UserName: user.Username,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResp{
		UserID:      fmt.Sprint(user.Id),
		Username:    user.Username,
		AccessToken: token,
	}, nil
}

// Register 处理用户注册
func (s *authService) Register(req *RegisterReq) (*RegisterResp, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, ErrUserExists
	}

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		Nickname:     req.Nickname,
		Email:        req.Email,
		Status:       "enabled",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 保存用户
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &RegisterResp{
		UserID:   fmt.Sprint(user.Id),
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
	}, nil
}
