package handlerV1

import (
	"context"
	"net/http"

	payloadV1 "example.com/backend/api/v1/payload"
	usecaseV1 "example.com/backend/api/v1/usecase"
	"example.com/backend/internal/exception"
	"example.com/backend/pkg/response"
	"example.com/backend/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type AuthHandlerV1 struct {
	authUseCase usecaseV1.IAuthUseCaseV1
	log         zerolog.Logger
}

func NewAuthHandlerV1(
	authUseCase usecaseV1.IAuthUseCaseV1,
	log zerolog.Logger,
) *AuthHandlerV1 {
	return &AuthHandlerV1{
		authUseCase: authUseCase,
		log:         log,
	}
}

// SignUp godoc
// @Summary      Register new user
// @Description  Create new user account with name, email, and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body     payloadV1.SignUpRequest true "Sign Up Request"
// @Success      200  {object} payloadV1.SignUpResponse
// @Failure      400  {object} response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /auth/sign-up [post]
func (h *AuthHandlerV1) SignUp(ctx *gin.Context) {
	var req payloadV1.SignUpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("invalid sign-up request")
		msg := validation.ExtractValidationError(err)
		response.Error(ctx, http.StatusBadRequest, msg, err.Error())
		return
	}

	resp, err := h.authUseCase.SignUp(context.Background(), req)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			h.log.Error().Err(ex.Err).Msg(ex.Message)
			response.Error(ctx, ex.Code, ex.Message, ex.Err)
			return
		}

		response.Error(ctx, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "User registered successfully", resp)
}

// SignIn godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body     payloadV1.SignInRequest true "Sign In Request"
// @Success      200  {object} payloadV1.SignInResponse
// @Failure      400  {object} response.ErrorResponse
// @Failure      401  {object} response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /auth/sign-in [post]
func (h *AuthHandlerV1) SignIn(ctx *gin.Context) {
	var req payloadV1.SignInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Error().Err(err).Msg("invalid signin request")
		msg := validation.ExtractValidationError(err)
		response.Error(ctx, http.StatusBadRequest, msg, err.Error())
		return
	}

	resp, err := h.authUseCase.SignIn(context.Background(), req)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			h.log.Error().Err(ex.Err).Msg(ex.Message)
			response.Error(ctx, ex.Code, ex.Message, ex.Err)
			return
		}

		response.Error(ctx, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Login successful", resp)
}
