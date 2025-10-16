package account

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Handler struct {
	store   IStore
	service IService
	log     zerolog.Logger
}

func NewHandler(
	store IStore,
	service IService,
	log zerolog.Logger,
) *Handler {
	return &Handler{
		store:   store,
		service: service,
		log:     log,
	}
}

func (h *Handler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()

	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind request")
		return c.JSON(400, echo.Map{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		h.log.Warn().Err(err).Msg("Validation failed")
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation error", "errors": errors})
	}

	tx := h.store.WithTx()
	if h.store.ExistsUserByEmailIgnoreCase(ctx, req.Email) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email already in use"})
	}

	var user *User
	user = req.ToUser()
	if err := h.store.SaveUser(ctx, user); err != nil {
		h.log.Error().Err(err).Msg("Failed to save user")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to sign up"})
	}

	var res *UserResponse
	res.FromUser(user)

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) SignIn(c echo.Context) error {
	var req SignInRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind request")
		return c.JSON(400, echo.Map{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		h.log.Warn().Err(err).Msg("Validation failed")
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation error", "errors": errors})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) Validate(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetMe(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
