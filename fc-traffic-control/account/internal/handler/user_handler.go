package handler

import (
	"net/http"

	"traffic-control/account/internal/domain"
	"traffic-control/account/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userService service.IUserService
	log         zerolog.Logger
}

// NewUserHandler initializes the handler with service + logger
func NewUserHandler(userService service.IUserService, log zerolog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log.With().Str("component", "handler.user").Logger(),
	}
}

// RegisterRoutes registers all user endpoints under /users
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	users := e.Group("/users")
	users.GET("", h.GetAllUsers)
	users.GET("/:id", h.GetUserByID)
	users.POST("", h.CreateUser)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}

// GetAllUsers handles GET /users
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := h.userService.FindAll(ctx)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to fetch users")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	user, err := h.userService.FindByID(ctx, id)
	if err != nil {
		h.log.Error().Err(err).Str("id", id).Msg("Failed to fetch user")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch user"})
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req domain.User
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Invalid request body")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	ctx := c.Request().Context()
	if err := h.userService.Create(ctx, &req); err != nil {
		h.log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfully",
		"user":    req,
	})
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var req domain.User

	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Invalid request body")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	req.ID = id
	ctx := c.Request().Context()

	if err := h.userService.Update(ctx, &req); err != nil {
		h.log.Error().Err(err).Str("id", id).Msg("Failed to update user")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User updated successfully",
		"user":    req,
	})
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	if err := h.userService.Delete(ctx, id); err != nil {
		h.log.Error().Err(err).Str("id", id).Msg("Failed to delete user")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete user"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}
