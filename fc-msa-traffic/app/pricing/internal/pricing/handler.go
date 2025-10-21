package pricing

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Handler struct {
	store  IStore
	client IClient
	log    zerolog.Logger
}

func NewHandler(
	store IStore,
	client IClient,
	log zerolog.Logger,
) *Handler {
	return &Handler{
		store:  store,
		client: client,
		log:    log,
	}
}

// CreatePricingRule godoc
// @Summary      Create or update pricing rule for a product
// @Description  Set pricing rule including markup, discount, and thresholds for a specific product
// @Tags         pricing
// @Accept       json
// @Produce      json
// @Param        request  body      PricingRuleRequest  true  "Pricing Rule Request"
// @Success      201      {object}  map[string]interface{}  "Successfully created pricing rule"
// @Failure      400      {object}  map[string]interface{}  "Validation or binding error"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /rules [post]
func (h *Handler) CreatePricingRule(c echo.Context) error {
	ctx := c.Request().Context()

	var req PricingRuleRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind request")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
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

	price, err := h.client.GetPriceProduct(ctx, req.ProductID)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get product price when creating pricing rule")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get product price"})
	}

	if err := h.store.SetPricingRule(ctx, req.ToPricingRule(price)); err != nil {
		h.log.Error().Err(err).Msg("Failed to create pricing rule")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create pricing rule"})
	}

	h.log.Info().Str("product_id", req.ProductID).Msg("Pricing rule created successfully")
	return c.JSON(http.StatusCreated, echo.Map{"message": "Pricing rule created successfully"})
}

// GetPricing godoc
// @Summary      Get calculated pricing for a product
// @Description  Retrieve pricing rule and calculate final price based on current stock level
// @Tags         pricing
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  Pricing  "Calculated pricing result"
// @Failure      400  {object}  map[string]interface{}  "Invalid request"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /{id} [get]
func (h *Handler) GetPricing(c echo.Context) error {
	ctx := c.Request().Context()

	productID := c.Param("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Product ID is required"})
	}

	pricingRule, err := h.store.GetPricingRule(ctx, productID)
	if err != nil || pricingRule == nil {
		h.log.Error().Err(err).Msg("Failed to get pricing rule")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Pricing rule not found"})
	}

	stock, err := h.client.GetStockProduct(ctx, productID)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get product stock")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get product stock"})
	}

	// Calculate price based on stock
	markup := pricingRule.DefaultMarkup
	discount := pricingRule.DefaultDiscount

	// Apply adjustment if stock is below threshold
	if stock < pricingRule.StockThreshold {
		markup += pricingRule.MarkupIncrease
		discount -= pricingRule.DiscountReduction
		if discount < 0 {
			discount = 0 // Prevent negative discount
		}
	}

	productPrice := pricingRule.ProductPrice
	finalPrice := productPrice * (1 + markup) * (1 - discount)

	h.log.Info().Str("product_id", productID).Float64("final_price", finalPrice).Float64("markup", markup).Float64("discount", discount).Int("stock", stock).Msg("Pricing calculated successfully")
	return c.JSON(http.StatusOK, &Pricing{
		ProductID:  productID,
		Discount:   discount,
		Markup:     markup,
		FinalPrice: finalPrice,
	})
}
