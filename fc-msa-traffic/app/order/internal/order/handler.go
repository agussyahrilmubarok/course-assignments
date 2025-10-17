package order

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

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Create a new order with calculated pricing, markup, and discount
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request  body      CreateOrderRequest  true  "Order creation request"
// @Success      200      {object}  Order               "Successfully created order"
// @Failure      400      {object}  map[string]interface{}  "Invalid request or validation error"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /flash [post]
func (h *Handler) CreateFlashOrder(c echo.Context) error {
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

	order, err := h.service.CalculateAndCreateOrder(ctx, req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to create order")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create order"})
	}

	h.log.Info().Str("order_id", order.ID).Msg("Create order successfully")
	return c.JSON(http.StatusOK, order)
}

// CancelOrder godoc
// @Summary      Cancel an order
// @Description  Cancel an existing order and reverse its stock and pricing effects
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request  body      CancelOrderRequest  true  "Order cancellation request"
// @Success      200      {object}  Order               "Successfully canceled order"
// @Failure      400      {object}  map[string]interface{}  "Invalid request or validation error"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /cancel [post]
func (h *Handler) CancelOrder(c echo.Context) error {
	ctx := c.Request().Context()

	var req CancelOrderRequest
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

	order, err := h.service.CancelOrderAndRestockProduct(ctx, req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to cancel order")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to cancel order"})
	}

	h.log.Info().Str("order_id", order.ID).Msg("Cancel order successfully")
	return c.JSON(http.StatusOK, order)
}

// GetOrder godoc
// @Summary      Get order by ID
// @Description  Retrieve detailed order information including items and final amount
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  Order  "Order details"
// @Failure      400  {object}  map[string]interface{}  "Invalid request"
// @Failure      404  {object}  map[string]interface{}  "Order not found"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /{id} [get]
func (h *Handler) GetOrder(c echo.Context) error {
	ctx := c.Request().Context()

	orderID := c.Param("id")
	order, err := h.store.FindOrderByID(ctx, orderID)
	if err != nil {
		h.log.Error().Err(err).Str("order_id", orderID).Msg("Failed to cancel order")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to find order"})
	}

	h.log.Info().Str("order_id", order.ID).Msg("Cancel order successfully")
	return c.JSON(http.StatusOK, order)
}
