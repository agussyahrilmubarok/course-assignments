package dummy

import (
	"fmt"
	"time"

	"example.com/coupon-service/internal/config"
	"example.com/coupon-service/internal/coupon"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	db     *config.Postgres
	logger *zap.Logger
}

func NewHandler(db *config.Postgres, logger *zap.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
	}
}

func (h *Handler) InitDummyV1(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()

	// 5 real-world cases, mix of ongoing, future, past events
	events := []struct {
		Name        string
		Code        string
		Discount    int
		Type        coupon.DiscountType
		TotalQty    int
		StartOffset time.Duration
		EndOffset   time.Duration
	}{
		{"Black Friday Mega Sale", "BF-2025", 50, coupon.DiscountTypePercentage, 100, -1 * time.Hour, 24 * time.Hour},         // ongoing
		{"Christmas Special", "XMAS-2025", 20000, coupon.DiscountTypeFixedAmount, 50, 24 * time.Hour, 10 * 24 * time.Hour},    // future
		{"New Year Promo", "NY-2025", 15, coupon.DiscountTypePercentage, 75, -48 * time.Hour, -24 * time.Hour},                // past
		{"Regular Discount", "REG-2025", 10000, coupon.DiscountTypeFixedAmount, 200, -7 * 24 * time.Hour, 7 * 24 * time.Hour}, // ongoing
		{"Expired Promo", "EXP-2025", 5000, coupon.DiscountTypeFixedAmount, 20, -30 * 24 * time.Hour, -10 * 24 * time.Hour},   // expired
	}

	for _, e := range events {
		policyID := uuid.New().String()
		policy := &coupon.CouponPolicy{
			ID:                    policyID,
			Code:                  e.Code,
			Name:                  e.Name,
			Description:           fmt.Sprintf("%s promo", e.Name),
			TotalQuantity:         e.TotalQty,
			StartTime:             now.Add(e.StartOffset),
			EndTime:               now.Add(e.EndOffset),
			DiscountType:          e.Type,
			DiscountValue:         e.Discount,
			MinimumOrderAmount:    50000,
			MaximumDiscountAmount: 100000,
			CreatedAt:             now,
			UpdatedAt:             now,
		}

		_, err := h.db.Pool.Exec(ctx, `
			INSERT INTO coupon_policies (
				id, code, name, description, total_quantity,
				start_time, end_time, discount_type, discount_value,
				minimum_order_amount, maximum_discount_amount, created_at, updated_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		`, policy.ID, policy.Code, policy.Name, policy.Description, policy.TotalQuantity,
			policy.StartTime, policy.EndTime, policy.DiscountType, policy.DiscountValue,
			policy.MinimumOrderAmount, policy.MaximumDiscountAmount, policy.CreatedAt, policy.UpdatedAt)
		if err != nil {
			h.logger.Error("failed to insert policy", zap.Error(err))
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
	}

	h.logger.Info("coupon policy dummy data v1 successfully inserted")
	return c.JSON(200, map[string]string{"status": "coupon policy dummy data v1 initialized"})
}

func (h *Handler) CleanDummyV1(c echo.Context) error {
	ctx := c.Request().Context()
	_, err := h.db.Pool.Exec(ctx, `DELETE FROM coupon_policies`)
	if err != nil {
		h.logger.Error("failed to clean dummy data", zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	h.logger.Info("coupon policy dummy data v1 cleaned successfully")
	return c.JSON(200, map[string]string{"status": "coupon policy dummy data v1 cleaned"})
}

func (h *Handler) InitDummyV2(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()

	// 5 real-world cases, mix of ongoing, future, past events
	events := []struct {
		Name        string
		Code        string
		Discount    int
		Type        coupon.DiscountType
		TotalQty    int
		StartOffset time.Duration
		EndOffset   time.Duration
	}{
		{"Black Friday Mega Sale", "BF-2025", 50, coupon.DiscountTypePercentage, 100, -1 * time.Hour, 24 * time.Hour},         // ongoing
		{"Christmas Special", "XMAS-2025", 20000, coupon.DiscountTypeFixedAmount, 50, 24 * time.Hour, 10 * 24 * time.Hour},    // future
		{"New Year Promo", "NY-2025", 15, coupon.DiscountTypePercentage, 75, -48 * time.Hour, -24 * time.Hour},                // past
		{"Regular Discount", "REG-2025", 10000, coupon.DiscountTypeFixedAmount, 200, -7 * 24 * time.Hour, 7 * 24 * time.Hour}, // ongoing
		{"Expired Promo", "EXP-2025", 5000, coupon.DiscountTypeFixedAmount, 20, -30 * 24 * time.Hour, -10 * 24 * time.Hour},   // expired
	}

	for _, e := range events {
		policyID := uuid.New().String()
		policy := &coupon.CouponPolicy{
			ID:                    policyID,
			Code:                  e.Code,
			Name:                  e.Name,
			Description:           fmt.Sprintf("%s promo", e.Name),
			TotalQuantity:         e.TotalQty,
			StartTime:             now.Add(e.StartOffset),
			EndTime:               now.Add(e.EndOffset),
			DiscountType:          e.Type,
			DiscountValue:         e.Discount,
			MinimumOrderAmount:    50000,
			MaximumDiscountAmount: 100000,
			CreatedAt:             now,
			UpdatedAt:             now,
		}

		_, err := h.db.Pool.Exec(ctx, `
			INSERT INTO coupon_policies (
				id, code, name, description, total_quantity,
				start_time, end_time, discount_type, discount_value,
				minimum_order_amount, maximum_discount_amount, created_at, updated_at
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		`, policy.ID, policy.Code, policy.Name, policy.Description, policy.TotalQuantity,
			policy.StartTime, policy.EndTime, policy.DiscountType, policy.DiscountValue,
			policy.MinimumOrderAmount, policy.MaximumDiscountAmount, policy.CreatedAt, policy.UpdatedAt)
		if err != nil {
			h.logger.Error("failed to insert policy", zap.Error(err))
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
	}

	h.logger.Info("coupon policy dummy data v1 successfully inserted")
	return c.JSON(200, map[string]string{"status": "coupon policy dummy data v1 initialized"})
}

func (h *Handler) CleanDummyV2(c echo.Context) error {
	ctx := c.Request().Context()
	_, err := h.db.Pool.Exec(ctx, `DELETE FROM coupon_policies`)
	if err != nil {
		h.logger.Error("failed to clean dummy data", zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	h.logger.Info("coupon policy dummy data v1 cleaned successfully")
	return c.JSON(200, map[string]string{"status": "coupon policy dummy data v1 cleaned"})
}