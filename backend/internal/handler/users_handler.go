package handler

import (
	"fmt"
	"github.com/fynn404/gin-demo/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req = &service.GetUserProfileReq{
		UserID: fmt.Sprint(userID),
	}
	resp, err := h.userService.GetUserProfile(c, req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound, service.ErrUserIdInvalid:
			NotFound(c, "User not found")
		case service.ErrInvalidCredentials:
			Unauthorized(c, "Invalid credentials")
		default:
			InternalServerError(c, "Failed to login", err)
		}
		return
	}

	Success(c, resp)
}

func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	var req = &service.UpdateUserProfileReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body")
		return
	}

	resp, err := h.userService.UpdateUserProfile(c, req)
	if err != nil {
		switch err {
		case service.ErrUserExists, service.ErrUserIdInvalid:
			Error(c, 409, "User already exists", err)
		default:
			InternalServerError(c, "Failed to register user", err)
		}
		return
	}

	Success(c, resp)
}
