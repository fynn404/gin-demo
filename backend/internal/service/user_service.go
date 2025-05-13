package service

import (
	"fmt"
	"github.com/fynn404/gin-demo/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

type GetUserProfileReq struct {
	UserID string `json:"user_id"`
}

type GetUserProfileResp struct {
	Id        string    `json:"user_id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	AvatarUrl string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserProfileReq struct {
	Id        string `json:"user_id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	PassWord  string `json:"passWord"`
	AvatarUrl string `json:"avatar_url"`
}

type UpdateUserProfileResp struct {
	Nickname          string `json:"nickname"`
	Email             string `json:"email"`
	AvatarUrl         string `json:"avatar_url"`
	HasChangePassword bool   `json:"has_change_password"`
}

type UserService interface {
	GetUserProfile(c *gin.Context, req *GetUserProfileReq) (*GetUserProfileResp, error)
	UpdateUserProfile(c *gin.Context, req *UpdateUserProfileReq) (*UpdateUserProfileResp, error)
}

type User struct {
	userRepo repository.UserRepository
}

func (u *User) GetUserProfile(c *gin.Context, req *GetUserProfileReq) (*GetUserProfileResp, error) {
	userID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		return nil, ErrUserIdInvalid
	}
	user, err := u.userRepo.GetByID(uint64(userID))
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &GetUserProfileResp{
		Id:        req.UserID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		AvatarUrl: user.AvatarUrl,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
func (u *User) UpdateUserProfile(c *gin.Context, req *UpdateUserProfileReq) (*UpdateUserProfileResp, error) {
	var (
		hasChangePassword = false
	)
	userID, _ := c.Get("user_id")

	userId, err := strconv.ParseUint(fmt.Sprint(userID), 10, 64)
	if err != nil {
		return nil, ErrUserNotFound
	}
	user, err := u.userRepo.GetByID(userId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.AvatarUrl != "" {
		user.AvatarUrl = req.AvatarUrl
	}
	if req.PassWord != "" {
		// 生成密码哈希
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PassWord), bcrypt.DefaultCost)
		if err != nil { // 报错退出
			log.Println("generate hashedPassword failed")
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
		hasChangePassword = true
	}

	err = u.userRepo.Update(user)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return &UpdateUserProfileResp{
		Nickname:          user.Nickname,
		Email:             user.Email,
		AvatarUrl:         user.AvatarUrl,
		HasChangePassword: hasChangePassword,
	}, nil
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &User{userRepo: userRepo}
}
