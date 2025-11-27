package dummy

import (
	"errors"
	"fmt"
	"time"

	"example.com/coupon-service/internal/config"
	"example.com/coupon-service/internal/coupon"
	"example.com/coupon-service/internal/instrument/logging"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler struct {
	pg  *config.Postgres
	rdb *config.Redis
}

func NewHandler(
	pg *config.Postgres,
	rdb *config.Redis,
) *Handler {
	return &Handler{
		pg:  pg,
		rdb: rdb,
	}
}

// Dummy save in DB
func (h *Handler) InitDummyDB(c echo.Context) error {
	log := logging.GetLogger()
	ctx := c.Request().Context()
	now := time.Now().UTC()

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
		// CouponPolicy `ongoing` with 10 qoutas
		{"Black Friday Mega Sale 10", "BF-C10", 50, coupon.DiscountTypePercentage, 10, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 100 qoutas
		{"Black Friday Mega Sale 100", "BF-C100", 50, coupon.DiscountTypePercentage, 100, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 1_000 qoutas
		{"Black Friday Mega Sale 1000", "BF-C1k", 50, coupon.DiscountTypePercentage, 1000, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 10_000 qoutas
		{"Black Friday Mega Sale 1000", "BF-C10k", 50, coupon.DiscountTypePercentage, 10000, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 1_000_000 qoutas
		{"Black Friday Mega Sale 1m", "BF-C1m", 50, coupon.DiscountTypePercentage, 1000000, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 1_000_000 + 1 qoutas
		{"Black Friday Mega Sale 1m+1", "BF-C1m+1", 50, coupon.DiscountTypePercentage, 1000001, -1 * time.Hour, 24 * time.Hour},

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
		}

		// Insert Coupon Policy Records
		_, err := h.pg.Pool.Exec(ctx, `
			INSERT INTO coupon_policies (
				id, code, name, description, total_quantity,
				start_time, end_time, discount_type, discount_value,
				minimum_order_amount, maximum_discount_amount
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		`, policy.ID, policy.Code, policy.Name, policy.Description, policy.TotalQuantity,
			policy.StartTime.UTC(), policy.EndTime.UTC(), policy.DiscountType, policy.DiscountValue,
			policy.MinimumOrderAmount, policy.MaximumDiscountAmount)
		if err != nil {
			log.Error("failed to insert policy", zap.Error(err))
			return c.JSON(500, map[string]string{"error": err.Error()})
		}
	}

	log.Info("coupon policy dummy data db successfully inserted")
	return c.JSON(200, map[string]string{"status": "coupon policy dummy data db initialized"})
}

func (h *Handler) CleanDummyDB(c echo.Context) error {
	log := logging.GetLogger()
	ctx := c.Request().Context()

	// Delete CouponPolicy records in Postgres
	_, err := h.pg.Pool.Exec(ctx, `DELETE FROM coupon_policies`)
	if err != nil {
		log.Error("failed to clean dummy data", zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	log.Info("coupon policy dummy data db cleaned successfully")
	return c.JSON(200, map[string]string{"status": "coupon policy dummy data db cleaned"})
}

func (h *Handler) InitDummyRedisAndDB(c echo.Context) error {
	log := logging.GetLogger()
	ctx := c.Request().Context()
	now := time.Now().UTC()

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
		// CouponPolicy `ongoing` with 10 qoutas
		{"Black Friday Mega Sale 10", "BF-C10", 50, coupon.DiscountTypePercentage, 10, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 100 qoutas
		{"Black Friday Mega Sale 100", "BF-C100", 50, coupon.DiscountTypePercentage, 100, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 1_000 qoutas
		{"Black Friday Mega Sale 1000", "BF-C1k", 50, coupon.DiscountTypePercentage, 1000, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 10_000 qoutas
		{"Black Friday Mega Sale 1000", "BF-C10k", 50, coupon.DiscountTypePercentage, 10000, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 1_000_000 qoutas
		{"Black Friday Mega Sale 1m", "BF-C1m", 50, coupon.DiscountTypePercentage, 1000000, -1 * time.Hour, 24 * time.Hour},
		// CouponPolicy `ongoing` with 1_000_000 + 1 qoutas
		{"Black Friday Mega Sale 1m+1", "BF-C1m+1", 50, coupon.DiscountTypePercentage, 1000001, -1 * time.Hour, 24 * time.Hour},

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
		}

		// Insert CouponPolicy records in Postgres
		_, err := h.pg.Pool.Exec(ctx, `
			INSERT INTO coupon_policies (
				id, code, name, description, total_quantity,
				start_time, end_time, discount_type, discount_value,
				minimum_order_amount, maximum_discount_amount
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		`, policy.ID, policy.Code, policy.Name, policy.Description, policy.TotalQuantity,
			policy.StartTime.UTC(), policy.EndTime.UTC(), policy.DiscountType, policy.DiscountValue,
			policy.MinimumOrderAmount, policy.MaximumDiscountAmount)
		if err != nil {
			log.Error("failed to insert policy", zap.Error(err))
			return c.JSON(500, map[string]string{"error": err.Error()})
		}

		ttl := time.Until(now.Add(e.EndOffset))
		if ttl <= 0 {
			ttl = time.Millisecond
		}

		// Insert CouponPolicy quantity
		policyQuantityKey := "coupon:policy:quantity:" + e.Code
		if err := h.rdb.Client.Set(ctx, policyQuantityKey, e.TotalQty, ttl).Err(); err != nil {
			log.Error("failed to set policy quantity in Redis", zap.Error(err))
			continue
		}
	}

	log.Info("coupon policy dummy data redis and db successfully inserted")
	return c.JSON(200, map[string]string{"status": "coupon policy dummy data redis and db initialized"})
}

func (h *Handler) CleanDummyRedisAndDB(c echo.Context) error {
	log := logging.GetLogger()
	ctx := c.Request().Context()

	rows, err := h.pg.Pool.Query(ctx, `SELECT code FROM coupon_policies`)
	if err != nil {
		log.Error("failed to fetch policy codes", zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil {
			codes = append(codes, code)
		}
	}

	// Delete CouponPolicy in Redis
	for _, code := range codes {
		policyQuantityKey := "coupon:policy:quantity:" + code
		if err := h.rdb.Client.Del(ctx, policyQuantityKey).Err(); err != nil {
			log.Warn("failed to delete redis key", zap.String("key", policyQuantityKey), zap.Error(err))
		} else {
			log.Info("successfully deleted redis key", zap.String("key", policyQuantityKey))
		}
	}

	// Delete CouponPolicy records in Postgres
	_, err = h.pg.Pool.Exec(ctx, `DELETE FROM coupon_policies`)
	if err != nil {
		log.Error("failed to clean dummy data", zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	log.Info("coupon policy dummy data redis and db cleaned successfully")
	return c.JSON(200, map[string]string{"status": "coupon policy dummy data redis and db cleaned"})
}

// Dummy Validation

func (h *Handler) CheckQuantity(c echo.Context) error {
	log := logging.GetLogger()

	policyCode := c.Param("policy_code")
	if policyCode == "" {
		log.Error("invalid policy_code")
		return c.JSON(400, map[string]string{"error": "policy_code is required"})
	}

	ctx := c.Request().Context()

	var cp coupon.CouponPolicy
	err := h.pg.Pool.QueryRow(
		ctx,
		`SELECT id, code, name, description, total_quantity, start_time, end_time,
		        discount_type, discount_value, minimum_order_amount, maximum_discount_amount,
		        created_at, updated_at
		 FROM coupon_policies
		 WHERE code = $1
		 LIMIT 1`,
		policyCode,
	).Scan(
		&cp.ID,
		&cp.Code,
		&cp.Name,
		&cp.Description,
		&cp.TotalQuantity,
		&cp.StartTime,
		&cp.EndTime,
		&cp.DiscountType,
		&cp.DiscountValue,
		&cp.MinimumOrderAmount,
		&cp.MaximumDiscountAmount,
		&cp.CreatedAt,
		&cp.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(404, map[string]string{"error": "policy not found"})
		}
		log.Error("failed to get coupon policy", zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	var totalIssued int
	err = h.pg.Pool.QueryRow(
		ctx,
		`SELECT COUNT(*) 
		 FROM coupons 
		 WHERE coupon_policy_id = $1`,
		cp.ID,
	).Scan(&totalIssued)
	if err != nil {
		log.Error("failed to count coupon", zap.Error(err))
		return c.JSON(500, map[string]string{"error": err.Error()})
	}

	result := map[string]interface{}{
		"policy_code":        cp.Code,
		"policy_name":        cp.Name,
		"start_time":         cp.StartTime,
		"end_time":           cp.EndTime,
		"total_quantity":     cp.TotalQuantity,
		"total_issued":       totalIssued,
		"remaining_quantity": cp.TotalQuantity - totalIssued,
	}
	return c.JSON(200, result)
}
