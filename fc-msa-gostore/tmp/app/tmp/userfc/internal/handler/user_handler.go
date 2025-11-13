package handler

import (
	"net/http"

	"example.com/pkg/model"
	"example.com/userfc/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type userHandler struct {
	userService service.IUserService
	log         *zap.Logger
}

func NewUserHandler(
	userService service.IUserService,
	log *zap.Logger,
) *userHandler {
	return &userHandler{
		userService: userService,
		log:         log,
	}
}

func (h *userHandler) SignUp(c *gin.Context) {
	var req model.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid signup request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	response, err := h.userService.SignUp(c.Request.Context(), &req)
	if err != nil || response == nil {
		h.log.Error("failed to sign up", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed sign up user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) SignIn(c *gin.Context) {
	var req model.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid signup request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	response, err := h.userService.SignIn(c.Request.Context(), &req)
	if err != nil || response == nil {
		h.log.Error("failed to sign in", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed sign in user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) FindUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.log.Warn("invalid param request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid param"})
		return
	}

	response, err := h.userService.FindByID(c.Request.Context(), id)
	if err != nil || response == nil {
		h.log.Error("failed to get user by id", zap.String("user_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
