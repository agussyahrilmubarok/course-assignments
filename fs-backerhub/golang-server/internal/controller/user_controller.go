package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	"example.com.backend/internal/service"
	"example.com.backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	USER_UPLOAD_PATH = "public/uploads/users"
)

type userController struct {
	baseController
	userService   service.IUserService
	uploadService service.IUploadService
}

func NewUserController(
	userService service.IUserService,
	uploadService service.IUploadService,
) *userController {
	return &userController{
		userService:   userService,
		uploadService: uploadService,
	}
}

func (h *userController) Index(c *gin.Context) {
	data := gin.H{"title": "Users"}

	ctx := c.Request.Context()

	h.showOnlyUsersIndex(c, ctx, data)
}

func (h *userController) Add(c *gin.Context) {
	data := gin.H{"title": "New User"}

	h.renderHTML(c, http.StatusOK, "user_add.html", data)
}

func (h *userController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "New User"}

	var input struct {
		Name       string `form:"name" binding:"required"`
		Email      string `form:"email" binding:"required,email"`
		Occupation string `form:"occupation" binding:"required"`
		Password   string `form:"password" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		log.Warn("failed to bind add user input", zap.Error(err))
		data["form"] = input
		data["error"] = "Invalid input."
		h.renderHTML(c, http.StatusBadRequest, "user_add.html", data)
		return
	}

	exists, _ := h.userService.ExistsByEmailIgnoreCase(ctx, input.Email)
	if exists {
		log.Warn("attempt to register with existing email", zap.String("user_email", input.Email))
		data["form"] = input
		data["error"] = "Email already used."
		h.renderHTML(c, http.StatusBadRequest, "user_add.html", data)
		return
	}

	userDto := model.UserDTO{
		Name:       input.Name,
		Email:      input.Email,
		Occupation: input.Occupation,
		Password:   input.Password,
	}
	if err := h.userService.Create(ctx, userDto); err != nil {
		log.Error("failed to create user", zap.Error(err), zap.String("user_email", input.Email))
		data["form"] = input
		data["error"] = "An error occurred while processing the request."
		h.renderHTML(c, http.StatusInternalServerError, "user_add.html", data)
		return
	}

	data["title"] = "Users"
	data["form"] = nil
	data["success"] = fmt.Sprintf("New user created successfully: %s", input.Email)

	h.showOnlyUsersIndex(c, ctx, data)
}

func (h *userController) Edit(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Edit User"}

	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		log.Error("user not found", zap.Error(err), zap.String("user_id", idStr))
		data["error"] = "An error occurred while retrieving user details."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	data["user"] = userDto

	h.renderHTML(c, http.StatusOK, "user_edit.html", data)
}

func (h *userController) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Edit User"}

	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		log.Error("user not found", zap.Error(err), zap.String("id", idStr))
		data["error"] = "User not found."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	var input struct {
		ID         string
		Name       string `form:"name" binding:"required"`
		Email      string `form:"email" binding:"required,email"`
		Occupation string `form:"occupation" binding:"required"`
		Password   string `form:"password" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		log.Warn("invalid bind user update", zap.Error(err))
		data["user"] = input
		data["error"] = "Invalid input."
		h.renderHTML(c, http.StatusBadRequest, "user_edit.html", data)
		return
	}

	input.ID = idStr
	if input.Email != userDto.Email {
		exists, _ := h.userService.ExistsByEmailIgnoreCase(ctx, input.Email)
		if exists {
			log.Warn("email already used by another user", zap.String("user_email", input.Email))
			data["user"] = input
			data["error"] = "Email is already in use."
			h.renderHTML(c, http.StatusBadRequest, "user_edit.html", data)
			return
		}
	}

	userDto.Name = input.Name
	userDto.Email = input.Email
	userDto.Password = input.Password
	userDto.Occupation = input.Occupation

	if err := h.userService.Update(ctx, *userDto); err != nil {
		log.Error("user update failed", zap.Error(err), zap.String("id", userDto.ID))
		data["user"] = input
		data["error"] = "Failed to update user."
		h.renderHTML(c, http.StatusInternalServerError, "user_edit.html", data)
		return
	}

	data["title"] = "Users"
	data["success"] = fmt.Sprintf("User updated successfully: %s", input.Email)

	h.showOnlyUsersIndex(c, ctx, data)
}

func (h *userController) Avatar(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Edit Avatar User"}

	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		log.Error("user not found for avatar", zap.Error(err), zap.String("user_id", idStr))
		data["error"] = "User not found."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	data["user"] = userDto

	h.renderHTML(c, http.StatusOK, "user_avatar.html", data)
}

func (h *userController) Upload(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Edit Avatar User"}

	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		log.Error("user not found for avatar upload", zap.Error(err), zap.String("user_id", idStr))
		data["error"] = "User not found."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	avatarFile, err := c.FormFile("avatar")
	if err != nil {
		log.Error("failed to retrieve avatar", zap.Error(err))
		data["error"] = "Failed to get avatar file."
		h.renderHTML(c, http.StatusBadRequest, "user_avatar.html", data)
		return
	}

	avatarPath, err := h.uploadService.SaveLocal(ctx, USER_UPLOAD_PATH, avatarFile, idStr)
	if err != nil {
		log.Error("failed to upload avatar file", zap.Error(err), zap.String("id", idStr))
		data["error"] = "Failed to upload avatar file."
		h.renderHTML(c, http.StatusInternalServerError, "user_avatar.html", data)
		return
	}

	trimmedPath := strings.TrimPrefix(avatarPath, fmt.Sprintf("%v/", USER_UPLOAD_PATH))
	userDto.ImageName = trimmedPath

	if err := h.userService.Update(ctx, *userDto); err != nil {
		h.uploadService.RemoveLocal(ctx, fmt.Sprintf("%v/", USER_UPLOAD_PATH), trimmedPath)
		log.Error("failed to save avatar user", zap.Error(err), zap.String("id", idStr))
		data["error"] = "Failed to save avatar user."
		h.renderHTML(c, http.StatusInternalServerError, "user_avatar.html", data)
		return
	}

	data["title"] = "Users"
	data["success"] = "User avatar updated successfully"

	h.showOnlyUsersIndex(c, ctx, data)
}

func (h *userController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Delete User"}

	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		log.Error("user not found for deletion", zap.Error(err), zap.String("user_id", idStr))
		data["error"] = "User not found."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	if userDto.ImageName != "default.png" {
		if err := h.uploadService.RemoveLocal(ctx, fmt.Sprintf("%v/", USER_UPLOAD_PATH), userDto.ImageName); err != nil {
			log.Warn("failed to remove user image", zap.Error(err), zap.String("image", userDto.ImageName))
		}
	}

	if err := h.userService.DeleteByID(ctx, userDto.ID); err != nil {
		log.Error("user delete failed", zap.Error(err), zap.String("id", userDto.ID))
		data["error"] = "Failed to delete a user."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	data["title"] = "Users"
	data["success"] = "User has been deleted successfully"

	h.showOnlyUsersIndex(c, ctx, data)
}

func (h *userController) showOnlyUsersIndex(c *gin.Context, ctx context.Context, data gin.H) {
	userDtos, err := h.userService.FindAll(ctx)
	if err != nil {
		data["users"] = []model.UserDTO{}
		h.renderHTML(c, http.StatusBadRequest, "user_index.html", data)
		return
	}

	var onlyUsers []model.UserDTO
	for _, userDto := range userDtos {
		if userDto.Role == domain.RoleUser {
			onlyUsers = append(onlyUsers, userDto)
		}
	}

	data["users"] = onlyUsers

	h.renderHTML(c, http.StatusOK, "user_index.html", data)
}
