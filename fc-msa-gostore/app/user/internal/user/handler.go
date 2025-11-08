package user

import (
	"context"

	"example.com/user/pkg/exception"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	service IService
	log     *zerolog.Logger
}

func NewHandler(service IService, log *zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

// SignUp godoc
// @Summary User Sign Up
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "User sign-up payload"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Router /auth/signup [post]
func (h *Handler) SignUp(c *fiber.Ctx) error {
	ctx := context.Background()

	var req SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn().Msg("Invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	resp, err := h.service.SignUp(ctx, req)
	if err != nil {
		if httpErr, ok := err.(*exception.Http); ok {
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
// @Param request body SignInRequest true "User sign-in payload"
// @Success 200 {object} UserWithTokenResponse
// @Failure 400 {object} map[string]string
// @Router /auth/signin [post]
func (h *Handler) SignIn(c *fiber.Ctx) error {
	var req SignInRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	resp, err := h.service.SignIn(context.Background(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// FindByID godoc
// @Summary Get User by ID
// @Description Retrieve a user's data by their unique ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 404 {object} map[string]string
// @Router /{id} [get]
func (h *Handler) FindByID(c *fiber.Ctx) error {
	userID := c.Params("id")

	resp, err := h.service.FindByID(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
