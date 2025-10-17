package order

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Handler struct {
	store IStore
	log   zerolog.Logger
}

func NewHandler(
	store IStore,
	log zerolog.Logger,
) *Handler {
	return &Handler{
		store: store,
		log:   log,
	}
}

func (h *Handler) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateOrderRequest
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

func (h *Handler) CancelOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
