package v1

import (
	"context"

	"example.com/coupon-service/internal/config"
	"example.com/coupon-service/internal/coupon"
	"go.uber.org/zap"
)

type IRepository interface {
	FindCouponPolicyByCode(ctx context.Context, code string) (*coupon.CouponPolicy, error)
	FindCouponPolicyByCodeLock(ctx context.Context, code string) (*coupon.CouponPolicy, error)
}

type repository struct {
	db     *config.Postgres
	logger *zap.Logger
}

func NewRepository(
	db *config.Postgres,
	logger *zap.Logger,
) IRepository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) FindCouponPolicyByCode(ctx context.Context, code string) (*coupon.CouponPolicy, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT 
			id,
			code,
			name,
			description,
			total_quantity,
			start_time,
			end_time,
			discount_type,
			discount_value,
			minimum_order_amount,
			maximum_discount_amount,
			created_at,
			updated_at
		FROM coupon_policies
		WHERE code = $1
		LIMIT 1
	`, code)

	var policy coupon.CouponPolicy

	err := row.Scan(
		&policy.ID,
		&policy.Code,
		&policy.Name,
		&policy.Description,
		&policy.TotalQuantity,
		&policy.StartTime,
		&policy.EndTime,
		&policy.DiscountType,
		&policy.DiscountValue,
		&policy.MinimumOrderAmount,
		&policy.MaximumDiscountAmount,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("failed to fetch coupon policy by code", zap.String("code", code), zap.Error(err))
		return nil, err
	}

	return &policy, nil
}

func (r *repository) FindCouponPolicyByCodeLock(ctx context.Context, code string) (*coupon.CouponPolicy, error) {
	panic("unimplemented")
}
