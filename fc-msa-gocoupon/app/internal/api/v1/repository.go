package v1

import (
	"context"
	"errors"

	"example.com/coupon-service/internal/config"
	"example.com/coupon-service/internal/coupon"
	"go.uber.org/zap"
)

type IRepository interface {
	FindCouponPolicyByCode(ctx context.Context, code string) (*coupon.CouponPolicy, error)
	CreateCoupon(ctx context.Context, coupon *coupon.Coupon) (*coupon.Coupon, error)
	CountIssuedCoupons(ctx context.Context, policyID string) (int, error)
}

type repository struct {
	pg     *config.Postgres
	logger *zap.Logger
}

func NewRepository(
	pg *config.Postgres,
	logger *zap.Logger,
) IRepository {
	return &repository{
		pg:     pg,
		logger: logger,
	}
}

func (r *repository) FindCouponPolicyByCode(ctx context.Context, code string) (*coupon.CouponPolicy, error) {
	row := r.pg.Pool.QueryRow(ctx, `
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
		r.logger.Error("failed to fetch coupon policy by code", zap.String("policy_code", code), zap.Error(err))
		return nil, coupon.ErrCouponPolicyNotFound
	}

	return &policy, nil
}

func (r *repository) CreateCoupon(ctx context.Context, c *coupon.Coupon) (*coupon.Coupon, error) {
	row := r.pg.Pool.QueryRow(ctx, `
		INSERT INTO coupons (
			id,
			code,
			status,
			used_at,
			user_id,
			order_id,
			coupon_policy_id,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
		)
		RETURNING 
			id,
			code,
			status,
			used_at,
			user_id,
			order_id,
			coupon_policy_id,
			created_at,
			updated_at
	`,
		c.ID,
		c.Code,
		c.Status,
		c.UsedAt,
		c.UserID,
		c.OrderID,
		c.CouponPolicyID,
	)

	var result coupon.Coupon
	err := row.Scan(
		&result.ID,
		&result.Code,
		&result.Status,
		&result.UsedAt,
		&result.UserID,
		&result.OrderID,
		&result.CouponPolicyID,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("failed to create coupon", zap.Error(err))
		return nil, errors.New("")
	}

	return &result, nil
}

func (r *repository) CountIssuedCoupons(ctx context.Context, policyID string) (int, error) {
	row := r.pg.Pool.QueryRow(ctx, `
        SELECT COUNT(*) 
        FROM coupons
        WHERE coupon_policy_id = $1
    `, policyID)

	var count int
	if err := row.Scan(&count); err != nil {
		r.logger.Error("failed to count issued coupons", zap.String("policy_id", policyID), zap.Error(err))
		return 0, coupon.ErrCouponCounted
	}

	return count, nil
}
