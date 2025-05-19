package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fynn404/gin-demo/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	// Get the file from form data
	file, err := c.FormFile("avatar")
	if err != nil {
		BadRequest(c, "No file uploaded")
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		BadRequest(c, "Invalid file type. Only JPG, JPEG and PNG are allowed")
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + ext

	// Create uploads directory if it doesn't exist
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		InternalServerError(c, "Failed to create upload directory", err)
		return
	}

	// Save the file
	filepath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		InternalServerError(c, "Failed to save file", err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := c.GetString("user_id")

	// Update user's avatar URL in database
	avatarURL := "/uploads/avatars/" + filename
	req := &service.UpdateUserProfileReq{
		Id:        userID,
		AvatarUrl: avatarURL,
	}

	resp, err := h.userService.UpdateUserProfile(c, req)
	if err != nil {
		InternalServerError(c, "Failed to update user avatar", err)
		return
	}

	Success(c, resp)
}
