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
