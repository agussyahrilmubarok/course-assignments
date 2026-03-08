package handlerV1

import (
	"errors"
	"net/http"

	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/logger"
	"example.com.backend/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	usecaseV1 "example.com.backend/internal/rest/v1/usecase"
)

type ProfileHandlerV1 struct {
	userUseCase usecaseV1.IUserUseCaseV1
}

func NewProfileHandlerV1(
	userUseCase usecaseV1.IUserUseCaseV1,
) *ProfileHandlerV1 {
	return &ProfileHandlerV1{
		userUseCase: userUseCase,
	}
}

// GetMe godoc
// @Summary      Get Profile
// @Description  Get logged in user profile
// @Tags         Profile
// @Produce      json
// @Success      200  {object} payloadV1.UserResponse
// @Failure      401  {object} response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /profiles/me [get]
// @Security     BearerAuth
func (h *ProfileHandlerV1) GetMe(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing user id in context"))
		log.Error("unauthorized access", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid user id type"))
		log.Error("unauthorized access", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	resp, err := h.userUseCase.GetMe(ctx, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error("failed to get profile", zap.String("message", ex.Message), zap.Error(ex.Err))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}

		log.Error("unexpected error during get profile", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	log.Info("profile retrieved successfully", zap.String("user_id", userID))
	response.Success(c, http.StatusOK, "Profile found", resp)
}
