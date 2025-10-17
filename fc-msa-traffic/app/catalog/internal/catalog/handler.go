package catalog

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

// GetProducts godoc
// @Summary      Get all products
// @Description  Retrieve list of all products
// @Tags         products
// @Produce      json
// @Success      200  {array}   Product
// @Failure      400  {object}  map[string]interface{}
// @Router       /products [get]
func (h *Handler) GetProducts(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.store.FindProducts(ctx)
	if err != nil {
		h.log.Warn().Msg("Failed to find all products")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Empty products"})
	}

	h.log.Info().Int("products_count", len(products)).Msg("Returning products list")
	return c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary      Get product by ID
// @Description  Retrieve product details by product ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  Product
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func (h *Handler) GetProduct(c echo.Context) error {
	ctx := c.Request().Context()

	productID := c.Param("id")
	product, err := h.store.FindProductByID(ctx, productID)
	if err != nil || product == nil {
		h.log.Warn().Str("product_id", productID).Msg("Failed to find product by id")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Product not found"})
	}

	h.log.Info().Str("product_id", productID).Msg("Returning product detail")
	return c.JSON(http.StatusOK, product)
}

// ReverseProductStock godoc
// @Summary      Reverse product stock
// @Description  Decrease product stock by quantity
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        request  body      ReserveStockRequest  true  "Reverse Stock Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /products/reverse [post]
func (h *Handler) ReverseProductStock(c echo.Context) error {
	ctx := c.Request().Context()

	var req ReserveStockRequest
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

	err := h.store.ReverseProductByID(ctx, req.ProductID, req.Quantity)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to reverse product stock")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to reverse product stock"})
	}

	h.log.Info().Str("product_id", req.ProductID).Int("quantity", req.Quantity).Msg("Reversed product stock successfully")
	return c.JSON(http.StatusOK, echo.Map{"message": "Stock reversed successfully"})
}

// ReleaseProductStock godoc
// @Summary      Release product stock
// @Description  Increase product stock by quantity
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        request  body      ReleaseStockRequest  true  "Release Stock Request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /products/release [post]
func (h *Handler) ReleaseProductStock(c echo.Context) error {
	ctx := c.Request().Context()

	var req ReleaseStockRequest
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

	err := h.store.ReleaseProductByID(ctx, req.ProductID, req.Quantity)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to release product stock")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to release product stock"})
	}

	h.log.Info().Str("product_id", req.ProductID).Int("quantity", req.Quantity).Msg("Released product stock successfully")
	return c.JSON(http.StatusOK, echo.Map{"message": "Stock released successfully"})
}

// GetProductStock godoc
// @Summary      Get product stock by ID
// @Description  Retrieve product stock by product ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func (h *Handler) GetProductStock(c echo.Context) error {
	ctx := c.Request().Context()

	productID := c.Param("id")
	product, err := h.store.FindProductByID(ctx, productID)
	if err != nil || product == nil {
		h.log.Warn().Str("product_id", productID).Msg("Failed to find product by id")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Product not found"})
	}

	h.log.Info().Str("product_id", productID).Int("stock", product.Stock).Msg("Returning product stock")
	return c.JSON(http.StatusOK, echo.Map{"stock": product.Stock})
}