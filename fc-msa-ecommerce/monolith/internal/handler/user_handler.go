package handler

import (
	"ecommerce/internal/middleware"
	"ecommerce/internal/model"
	"ecommerce/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	UserService service.IUserService
	Logger      zerolog.Logger
}

func NewUserHandler(userService service.IUserService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Logger:      logger,
	}
}

// GetCurrentUser godoc
// @Summary Get current logged-in user
// @Description Get user info from JWT
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrUnauthorized())
		return
	}

	user, err := h.UserService.GetByID(c.Request.Context(), userID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to get current user")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to get user", err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(http.StatusOK, "User profile fetched", user))
}

// UpdateCurrentUser godoc
// @Summary Update current logged-in user
// @Description Update user profile info
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.UpdateUserRequest true "Update user payload"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/users/me [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrUnauthorized())
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error().Err(err).Msg("invalid request body")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	updatedUser, err := h.UserService.Update(c.Request.Context(), userID, req)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to update user")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to update user", err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(http.StatusOK, "User updated successfully", updatedUser))
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users (admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 403 {object} model.ErrorResponse
// @Router /api/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.UserService.GetAll(c.Request.Context())
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to get all users")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to get users", err))
		return
	}

	c.JSON(http.StatusOK, model.NewSuccessResponse(http.StatusOK, "Users retrieved", users))
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by ID (admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid user ID", err))
		return
	}

	err = h.UserService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, model.ErrNotFound())
			return
		}
		h.Logger.Error().Err(err).Msg("failed to delete user")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to delete user", err))
		return
	}

	c.Status(http.StatusNoContent)
}
