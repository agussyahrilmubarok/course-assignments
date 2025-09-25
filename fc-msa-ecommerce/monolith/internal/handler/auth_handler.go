package handler

import (
	"net/http"

	"ecommerce/internal/model"
	"ecommerce/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	AuthService service.IAuthService
	Logger      zerolog.Logger
}

func NewAuthHandler(authService service.IAuthService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Logger:      logger,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.RegisterRequest true "Register payload"
// @Success 201 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error().Err(err).Msg("invalid register request")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	if err := h.AuthService.Register(c.Request.Context(), req); err != nil {
		h.Logger.Error().Err(err).Msg("failed to register user")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to register user", err))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Code:    http.StatusCreated,
		Message: "User registered successfully",
	})
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body model.LoginRequest true "Login payload"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error().Err(err).Msg("invalid login request")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	token, err := h.AuthService.Login(c.Request.Context(), req)
	if err != nil {
		h.Logger.Warn().Err(err).Msg("authentication failed")
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse(http.StatusUnauthorized, "Invalid email or password", err))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Login successful",
		Data:    gin.H{"token": token},
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout the current user (stateless JWT, no action performed)
// @Tags auth
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Stateless JWT logout (handled client-side)
	if err := h.AuthService.Logout(c.Request.Context()); err != nil {
		h.Logger.Error().Err(err).Msg("logout error")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Logout failed", err))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Logout successful",
	})
}
