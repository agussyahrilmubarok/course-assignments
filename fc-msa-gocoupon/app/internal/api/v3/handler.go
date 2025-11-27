package v3

import (
	"errors"

	"example.com/coupon-service/internal/api/middleware"
	"example.com/coupon-service/internal/coupon"
	"example.com/coupon-service/internal/instrument/logging"
	"example.com/coupon-service/internal/instrument/tracing"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	service IService
}

func NewHandler(service IService) *Handler {
	return &Handler{
		service: service,
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
// @Router       /coupons/issue [post]
func (h *Handler) IssueCoupon(c echo.Context) error {
	ctx, span := tracing.StartSpan(c.Request().Context(), "V3.Handler.IssueCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	var payload coupon.IssueCouponRequest
	if err := c.Bind(&payload); err != nil {
		span.RecordError(err)
		log.Warn("invalid body request", zap.Error(err))
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		err := errors.New("invalid x-user-id header")
		span.RecordError(err)
		log.Warn("invalid x-user-id header")
		return c.JSON(401, map[string]string{"error": err.Error()})
	}

	result, err := h.service.IssueCoupon(ctx, payload.PolicyCode, userID)
	if err != nil || result == nil {
		span.RecordError(err)
		log.Error("failed to issue coupon",
			zap.String("policy_code", payload.PolicyCode),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	log.Info("issue coupon successfully",
		zap.String("policy_code", payload.PolicyCode),
		zap.String("user_id", userID),
		zap.String("coupon_code", result.Code),
	)
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
// @Router       /coupons/use [post]
func (h *Handler) UseCoupon(c echo.Context) error {
	ctx, span := tracing.StartSpan(c.Request().Context(), "V3.Handler.IssueCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	var payload coupon.UseCouponRequest
	if err := c.Bind(&payload); err != nil {
		span.RecordError(err)
		log.Warn("invalid body request", zap.Error(err))
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		err := errors.New("invalid x-user-id header")
		span.RecordError(err)
		log.Warn("invalid x-user-id header")
		return c.JSON(401, map[string]string{"error": "invalid user id"})
	}

	result, err := h.service.UseCoupon(ctx, payload.CouponCode, userID, payload.OrderID)
	if err != nil || result == nil {
		span.RecordError(err)
		log.Error("failed to use coupon", zap.String("coupon_code", payload.CouponCode), zap.String("user_id", userID), zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	log.Info("use coupon successfully", zap.String("coupon_code", payload.CouponCode), zap.String("user_id", userID), zap.String("coupon_code", result.Code))
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
// @Router       /coupons/cancel [post]
func (h *Handler) CancelCoupon(c echo.Context) error {
	ctx, span := tracing.StartSpan(c.Request().Context(), "V3.Handler.CancelCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	var payload coupon.CancelCouponRequest
	if err := c.Bind(&payload); err != nil {
		span.RecordError(err)
		log.Warn("invalid body request", zap.Error(err))
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		err := errors.New("invalid x-user-id header")
		span.RecordError(err)
		log.Warn("invalid x-user-id header")
		return c.JSON(401, map[string]string{"error": "invalid user id"})
	}

	result, err := h.service.CancelCoupon(ctx, payload.CouponCode, userID)
	if err != nil || result == nil {
		span.RecordError(err)
		log.Error("failed to cancel coupon", zap.String("coupon_code", payload.CouponCode), zap.String("user_id", userID), zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	log.Info("cancel coupon successfully", zap.String("coupon_code", payload.CouponCode), zap.String("user_id", userID), zap.String("coupon_code", result.Code))
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
// @Router       /coupons/{coupon_code} [get]
func (h *Handler) FindCouponByCode(c echo.Context) error {
	ctx, span := tracing.StartSpan(c.Request().Context(), "V3.Handler.FindCouponByCode")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	couponCode := c.Param("coupon_code")
	if couponCode == "" {
		err := errors.New("invalid coupon_code")
		span.RecordError(err)
		log.Error("invalid coupon_code")
		return c.JSON(400, map[string]string{"error": "coupon_code is required"})
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		err := errors.New("invalid x-user-id header")
		span.RecordError(err)
		log.Warn("invalid x-user-id header")
		return c.JSON(401, map[string]string{"error": "invalid user id"})
	}

	result, err := h.service.FindCouponByCode(ctx, couponCode, userID)
	if err != nil || result == nil {
		span.RecordError(err)
		log.Error("failed to find coupon by code", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	log.Info("find coupon by code successfully", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.String("coupon_code", result.Code))
	return c.JSON(200, result)
}
