package v1

import (
	"example.com/coupon-service/internal/api/middleware"
	"example.com/coupon-service/internal/coupon"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	service IService
	logger  *zap.Logger
}

func NewHandler(
	service IService,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// IssueCoupon godoc
// @Summary      Issue a coupon for a user
// @Description  Issues a coupon under a specific policy code for the authenticated user
// @Tags         coupons
// @Accept       json
// @Produce      json
// @Param        X-USER-ID  header  string  true  "User ID"
// @Param        payload    body    coupon.IssueCouponRequest  true  "Issue coupon payload"
// @Success      200  {object}  coupon.Coupon
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/coupons/issue [post]
func (h *Handler) IssueCoupon(c echo.Context) error {
	var payload coupon.IssueCouponRequest
	if err := c.Bind(&payload); err != nil {
		h.logger.Warn("invalid body request", zap.Error(err))
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	ctx := c.Request().Context()
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.logger.Warn("invalid x-user-id header")
		return c.JSON(401, map[string]string{"error": "invalid user id"})
	}

	result, err := h.service.IssueCoupon(ctx, payload.PolicyCode, userID)
	if err != nil || result == nil {
		h.logger.Error("failed to issue coupon", zap.String("policy_code", payload.PolicyCode), zap.String("user_id", userID), zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	h.logger.Info("issue coupon successfully", zap.String("policy_code", payload.PolicyCode), zap.String("user_id", userID), zap.String("coupon_code", result.Code))
	return c.JSON(200, result)
}

// UseCoupon godoc
// @Summary      Use a coupon for an order
// @Description  Marks a coupon as used for the given order by the authenticated user
// @Tags         coupons
// @Accept       json
// @Produce      json
// @Param        X-USER-ID  header  string  true  "User ID"
// @Param        payload    body    coupon.UseCouponRequest  true  "Use coupon payload"
// @Success      200  {object}  coupon.Coupon
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/coupons/use [post]
func (h *Handler) UseCoupon(c echo.Context) error {
	var payload coupon.UseCouponRequest
	if err := c.Bind(&payload); err != nil {
		h.logger.Warn("invalid body request", zap.Error(err))
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	ctx := c.Request().Context()
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.logger.Warn("invalid x-user-id header")
		return c.JSON(401, map[string]string{"error": "invalid user id"})
	}

	result, err := h.service.UseCoupon(ctx, payload.CouponCode, userID, payload.OrderID)
	if err != nil || result == nil {
		h.logger.Error("failed to use coupon", zap.String("coupon_code", payload.CouponCode), zap.String("user_id", userID), zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	h.logger.Info("use coupon successfully", zap.String("coupon_code", payload.CouponCode), zap.String("user_id", userID), zap.String("coupon_code", result.Code))
	return c.JSON(200, result)
}

// CancelCoupon godoc
// @Summary      Cancel a coupon
// @Description  Cancels a coupon for the authenticated user
// @Tags         coupons
// @Accept       json
// @Produce      json
// @Param        X-USER-ID  header  string  true  "User ID"
// @Param        payload    body    coupon.CancelCouponRequest  true  "Cancel coupon payload"
// @Success      200  {object}  coupon.Coupon
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/coupons/cancel [post]
func (h *Handler) CancelCoupon(c echo.Context) error {
	var payload coupon.CancelCouponRequest
	if err := c.Bind(&payload); err != nil {
		h.logger.Warn("invalid body request", zap.Error(err))
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	ctx := c.Request().Context()
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.logger.Warn("invalid x-user-id header")
		return c.JSON(401, map[string]string{"error": "invalid user id"})
	}

	result, err := h.service.CancelCoupon(ctx, payload.CouponCode, userID)
	if err != nil || result == nil {
		h.logger.Error("failed to cancel coupon", zap.String("coupon_code", payload.CouponCode), zap.String("user_id", userID), zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	h.logger.Info("cancel coupon successfully", zap.String("coupon_code", payload.CouponCode), zap.String("user_id", userID), zap.String("coupon_code", result.Code))
	return c.JSON(200, result)
}

// FindCouponByCode godoc
// @Summary      Find coupon by code
// @Description  Retrieves coupon information for the authenticated user
// @Tags         coupons
// @Accept       json
// @Produce      json
// @Param        X-USER-ID     header  string  true  "User ID"
// @Param        coupon_code   path    string  true  "Coupon Code"
// @Success      200  {object}  coupon.Coupon
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/coupons/{coupon_code} [get]
func (h *Handler) FindCouponByCode(c echo.Context) error {
	couponCode := c.Param("coupon_code")
	if couponCode == "" {
		h.logger.Error("invalid coupon_code")
		return c.JSON(400, map[string]string{"error": "coupon_code is required"})
	}

	ctx := c.Request().Context()
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		h.logger.Warn("invalid x-user-id header")
		return c.JSON(401, map[string]string{"error": "invalid user id"})
	}

	result, err := h.service.FindCouponByCode(ctx, couponCode, userID)
	if err != nil || result == nil {
		h.logger.Error("failed to find coupon by code", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	h.logger.Info("find coupon by code successfully", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.String("coupon_code", result.Code))
	return c.JSON(200, result)
}
