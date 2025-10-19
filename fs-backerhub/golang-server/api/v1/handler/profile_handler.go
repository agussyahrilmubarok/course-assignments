package handlerV1

import (
	"errors"
	"net/http"

	usecaseV1 "example.com/backend/api/v1/usecase"
	"example.com/backend/internal/exception"
	"example.com/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type ProfileHandlerV1 struct {
	userUseCase usecaseV1.IUserUseCaseV1
	log         zerolog.Logger
}

func NewProfileHandlerV1(
	userUseCase usecaseV1.IUserUseCaseV1,
	log zerolog.Logger,
) *ProfileHandlerV1 {
	return &ProfileHandlerV1{
		userUseCase: userUseCase,
		log:         log,
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
func (h *ProfileHandlerV1) GetMe(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	resp, err := h.userUseCase.GetMe(ctx.Request.Context(), userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			h.log.Error().Err(ex.Err).Msg(ex.Message)
			response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
			return
		}

		h.log.Error().Err(err).Msg("failed to get profile")
		response.Error(ctx, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Profile found", resp)
}
