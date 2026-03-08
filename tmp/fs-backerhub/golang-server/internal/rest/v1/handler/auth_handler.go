package handlerV1

import (
	"net/http"

	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/logger"
	"example.com.backend/pkg/response"
	"example.com.backend/pkg/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	payloadV1 "example.com.backend/internal/rest/v1/payload"
	usecaseV1 "example.com.backend/internal/rest/v1/usecase"
)

type AuthHandlerV1 struct {
	authUseCase usecaseV1.IAuthUseCaseV1
}

func NewAuthHandlerV1(
	authUseCase usecaseV1.IAuthUseCaseV1,
) *AuthHandlerV1 {
	return &AuthHandlerV1{
		authUseCase: authUseCase,
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
func (h *AuthHandlerV1) SignUp(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	var req payloadV1.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid sign-up request", zap.Error(err))
		msg := validation.ExtractValidationError(err)
		response.Error(c, http.StatusBadRequest, msg, err.Error())
		return
	}

	resp, err := h.authUseCase.SignUp(ctx, req)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error("sign-up failed", zap.String("message", ex.Message), zap.Error(ex.Err))
			response.Error(c, ex.Code, ex.Message, ex.Err)
			return
		}

		log.Error("unexpected error during sign-up", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	log.Info("user registered successfully", zap.String("user_email", req.Email))
	response.Success(c, http.StatusCreated, "User registered successfully", resp)
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
func (h *AuthHandlerV1) SignIn(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	var req payloadV1.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("invalid sign-in request", zap.Error(err))
		msg := validation.ExtractValidationError(err)
		response.Error(c, http.StatusBadRequest, msg, err.Error())
		return
	}

	resp, err := h.authUseCase.SignIn(ctx, req)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error("sign-in failed", zap.String("message", ex.Message), zap.Error(ex.Err))
			response.Error(c, ex.Code, ex.Message, ex.Err)
			return
		}

		log.Error("unexpected error during sign-in", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	log.Info("user signed in successfully", zap.String("user_email", req.Email))
	response.Success(c, http.StatusOK, "Login successful", resp)
}
