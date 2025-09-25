package handler

import (
	"ecommerce/user-service/model"
	"ecommerce/user-service/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	authUseCase usecase.IAuthUseCase
	log         zerolog.Logger
}

func NewAuthHandler(authUseCase usecase.IAuthUseCase, log zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		log:         log,
	}
}

func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/auth")
	api.Post("/signup", h.SignUp)
	api.Post("/signin", h.SignIn)
}

// SignUp godoc
// @Summary      User registration
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.SignUpRequest true "Sign Up Request"
// @Success      200 {object} model.SignUpResponse
// @Failure      400 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/signup [post]
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var req model.SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := req.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.authUseCase.SignUp(c.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("SignUp failed")

		switch err {
		case usecase.ErrEmailAlreadyExists:
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
	}

	return c.Status(http.StatusOK).JSON(res)
}

// SignIn godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.SignInRequest true "Sign In Request"
// @Success      200 {object} model.SignInResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/signin [post]
func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
	var req model.SignInRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := req.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.authUseCase.SignIn(c.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("SignIn failed")

		switch err {
		case usecase.ErrInvalidCredentials:
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
	}

	return c.Status(http.StatusOK).JSON(res)
}
