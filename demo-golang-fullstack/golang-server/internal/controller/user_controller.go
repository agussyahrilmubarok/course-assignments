package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"example.com/backend/internal/domain"
	"example.com/backend/internal/model"
	"example.com/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	USER_UPLOAD_PATH = "public/uploads/users"
)

type userController struct {
	baseController
	userService   service.IUserService
	uploadService service.IUploadService
	log           zerolog.Logger
}

func NewUserController(
	userService service.IUserService,
	uploadService service.IUploadService,
	log zerolog.Logger,
) *userController {
	return &userController{
		userService:   userService,
		uploadService: uploadService,
		log:           log,
	}
}

func (h *userController) Index(c *gin.Context) {
	data := gin.H{
		"title": "Users",
	}
	ctx := c.Request.Context()
	h.showOnlyUsersIndex(c, ctx, data)
}

func (h *userController) Add(c *gin.Context) {
	data := gin.H{
		"title": "New User",
	}
	h.renderHTML(c, http.StatusOK, "user_add.html", data)
}

func (h *userController) Create(c *gin.Context) {
	data := gin.H{
		"title": "New User",
	}

	var input struct {
		Name       string `form:"name" binding:"required"`
		Email      string `form:"email" binding:"required,email"`
		Occupation string `form:"occupation" binding:"required"`
		Password   string `form:"password" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		h.log.Warn().Err(err).Msgf("failed to bind add user input")
		data["form"] = input
		data["error"] = "Invalid input."
		h.renderHTML(c, http.StatusBadRequest, "user_add.html", data)
		return
	}

	ctx := c.Request.Context()
	user, _ := h.userService.FindByEmail(ctx, input.Email)
	if user != nil {
		h.log.Error().Msgf("attempt to register with existing email %v", input.Email)
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
	err := h.userService.Create(ctx, userDto)
	if err != nil {
		h.log.Error().Err(err).Msgf("failed to create user %s", input.Email)
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
	data := gin.H{
		"title": "Edit User",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")

	userDto, err := h.userService.FindByID(ctx, idStr)
	if userDto == nil || err != nil {
		h.log.Error().Err(err).Msgf("user is not found id %s", idStr)
		data["title"] = "Users"
		data["error"] = "An error occurred while retrieving user details."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	data["user"] = userDto

	h.renderHTML(c, http.StatusOK, "user_edit.html", data)
}

func (h *userController) Update(c *gin.Context) {
	data := gin.H{
		"title": "Edit User",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		h.log.Error().Err(err).Msgf("user is not found id %s", idStr)
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
		h.log.Warn().Msgf("invalid bind user update")
		data["user"] = input
		data["error"] = "Invalid input."
		h.renderHTML(c, http.StatusBadRequest, "user_edit.html", data)
		return
	}

	input.ID = idStr
	if input.Email != userDto.Email {
		existinguser, _ := h.userService.FindByEmail(ctx, input.Email)
		if existinguser != nil {
			h.log.Warn().Msgf("email is already used by another user %s", input.Email)
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

	err = h.userService.Update(ctx, *userDto)
	if err != nil {
		h.log.Error().Err(err).Msgf("user update failed to id %s", userDto.ID)
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
	data := gin.H{
		"title": "Edit Avatar User",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		h.log.Error().Err(err).Msgf("user is not found id %s", idStr)
		data["error"] = "User not found."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	data["user"] = userDto

	h.renderHTML(c, http.StatusOK, "user_avatar.html", data)
}

func (h *userController) Upload(c *gin.Context) {
	data := gin.H{
		"title": "Edit Avatar User",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		h.log.Error().Err(err).Msgf("user is not found id %s", idStr)
		data["error"] = "User not found."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	avatarFile, err := c.FormFile("avatar")
	if err != nil {
		h.log.Error().Err(err).Msgf("failed to retrieve avatar")
		data["error"] = "Failed to get avatar file."
		h.renderHTML(c, http.StatusBadRequest, "user_avatar.html", data)
		return
	}

	avatarPath, err := h.uploadService.SaveLocal(USER_UPLOAD_PATH, avatarFile, idStr)
	if err != nil {
		h.log.Error().Err(err).Msgf("failed to upload avatar file %s", idStr)
		data["error"] = "Failed to upload avatar file."
		h.renderHTML(c, http.StatusInternalServerError, "user_avatar.html", data)
		return
	}

	trimmedPath := strings.TrimPrefix(avatarPath, fmt.Sprintf("%v/", USER_UPLOAD_PATH))
	userDto.ImageName = trimmedPath

	err = h.userService.Update(ctx, *userDto)
	if err != nil {
		h.uploadService.RemoveLocal(fmt.Sprintf("%v/", USER_UPLOAD_PATH), trimmedPath)
		h.log.Error().Err(err).Msgf("failed to save avatar user %s", idStr)
		data["error"] = "Failed to save avatar user."
		h.renderHTML(c, http.StatusInternalServerError, "user_avatar.html", data)
		return
	}

	data["title"] = "Users"
	data["success"] = "User avatar updated successfully"

	h.showOnlyUsersIndex(c, ctx, data)
}

func (h *userController) Delete(c *gin.Context) {
	data := gin.H{
		"title": "Delete User",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")
	userDto, err := h.userService.FindByID(ctx, idStr)
	if err != nil || userDto == nil {
		h.log.Error().Err(err).Msgf("user is not found id %s", idStr)
		data["error"] = "User not found."
		h.showOnlyUsersIndex(c, ctx, data)
		return
	}

	if userDto.ImageName != "default.png" {
		if err := h.uploadService.RemoveLocal(fmt.Sprintf("%v/", USER_UPLOAD_PATH), userDto.ImageName); err != nil {
			h.log.Error().Err(err).Msgf("failed to remove user image for id: %v", userDto.ImageName)
		}
	}

	err = h.userService.DeleteByID(ctx, userDto.ID)
	if err != nil {
		h.log.Error().Err(err).Msgf("user delete failed %s", userDto.ID)
		data["error"] = "Failed to delete a user."
		h.showOnlyUsersIndex(c, ctx, data)
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
