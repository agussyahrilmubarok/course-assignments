package user

import (
	"context"

	"github.com/agussyahrilmubarok/gox/pkg/xexception"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	service IService
	logger  zerolog.Logger
}

func NewHandler(
	service IService,
	logger zerolog.Logger,
) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// SignUp godoc
// @Summary User Sign Up
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SignUpParam true "User sign-up payload"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Router /auth/signup [post]
func (h *Handler) SignUp(c *fiber.Ctx) error {
	ctx := context.Background()

	var req SignUpParam
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Msg("invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	resp, err := h.service.SignUp(ctx, req)
	if err != nil {
		h.logger.Err(err).Msg("sign up user failed")
		if httpErr, ok := err.(*xexception.Http); ok {
			return c.Status(httpErr.Code).JSON(fiber.Map{"error": httpErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// SignIn godoc
// @Summary User Sign In
// @Description Log in a user and obtain an authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SignInParam true "User sign-in payload"
// @Success 200 {object} UserWithTokenResponse
// @Failure 400 {object} map[string]string
// @Router /auth/signin [post]
func (h *Handler) SignIn(c *fiber.Ctx) error {
	ctx := context.Background()

	var req SignInParam
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Msg("invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	resp, err := h.service.SignIn(ctx, req)
	if err != nil {
		h.logger.Err(err).Msg("sign in user failed")
		if httpErr, ok := err.(*xexception.Http); ok {
			return c.Status(httpErr.Code).JSON(fiber.Map{"error": httpErr.Message})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
