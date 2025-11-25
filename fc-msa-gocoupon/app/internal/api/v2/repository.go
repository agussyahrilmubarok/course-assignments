package v2

import (
	"context"
	"errors"

	"example.com/coupon-service/internal/config"
	"example.com/coupon-service/internal/coupon"
	"example.com/coupon-service/internal/instrument"
	"example.com/coupon-service/internal/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type IRepository interface {
	FindCouponPolicyByCodeForUpdateTx(ctx context.Context, tx pgx.Tx, code string) (*coupon.CouponPolicy, error)
	CountIssuedCouponsTx(ctx context.Context, tx pgx.Tx, policyID string) (int, error)
	CreateCouponTx(ctx context.Context, tx pgx.Tx, c *coupon.Coupon) (*coupon.Coupon, error)
	FindCouponByCode(ctx context.Context, code string) (*coupon.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon *coupon.Coupon) (*coupon.Coupon, error)
	FindCouponPolicyByID(ctx context.Context, id string) (*coupon.CouponPolicy, error)
	WithTx(ctx context.Context, fn func(pgx.Tx) error) error
}

type repository struct {
	pg *config.Postgres
}

func NewRepository(pg *config.Postgres) IRepository {
	return &repository{
		pg: pg,
	}
}

func (r *repository) FindCouponPolicyByCodeForUpdateTx(ctx context.Context, tx pgx.Tx, code string) (*coupon.CouponPolicy, error) {
	ctx, span := instrument.StartSpan(ctx, "V2.Repository.FindCouponPolicyByCodeForUpdateTx")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	row := tx.QueryRow(ctx, `
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
		FOR UPDATE
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
		span.RecordError(err)
		log.Error("failed to fetch coupon policy by code", zap.String("policy_code", code), zap.Error(err))
		return nil, coupon.ErrCouponPolicyNotFound
	}

	log.Info("fetched coupon policy successfully", zap.String("policy_id", policy.ID), zap.String("policy_code", code))
	return &policy, nil
}

func (r *repository) CountIssuedCouponsTx(ctx context.Context, tx pgx.Tx, policyID string) (int, error) {
	ctx, span := instrument.StartSpan(ctx, "V2.Repository.CountIssuedCouponsTx")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	row := tx.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM coupons
        WHERE coupon_policy_id = $1
    `, policyID)

	var count int
	if err := row.Scan(&count); err != nil {
		span.RecordError(err)
		log.Error("failed to count issued coupons", zap.String("policy_id", policyID), zap.Error(err))
		return 0, coupon.ErrCouponCounted
	}

	log.Info("counted issued coupons successfully", zap.String("policy_id", policyID), zap.Int("issued_count", count))
	return count, nil
}

func (r *repository) CreateCouponTx(ctx context.Context, tx pgx.Tx, c *coupon.Coupon) (*coupon.Coupon, error) {
	ctx, span := instrument.StartSpan(ctx, "V2.Repository.CreateCouponTx")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	row := tx.QueryRow(ctx, `
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
		span.RecordError(err)
		log.Error("failed to create coupon", zap.Error(err))
		return nil, errors.New("")
	}

	log.Info("coupon created successfully", zap.String("coupon_id", result.ID), zap.String("coupon_code", result.Code))
	return &result, nil
}

func (r *repository) FindCouponByCode(ctx context.Context, code string) (*coupon.Coupon, error) {
	ctx, span := instrument.StartSpan(ctx, "V2.Repository.FindCouponByCode")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	row := r.pg.Pool.QueryRow(ctx, `
		SELECT 
			id,
			code,
			status,
			used_at,
			user_id,
			order_id,
			coupon_policy_id,
			created_at,
			updated_at
		FROM coupons
		WHERE code = $1
		LIMIT 1
	`, code)

	var c coupon.Coupon
	err := row.Scan(
		&c.ID,
		&c.Code,
		&c.Status,
		&c.UsedAt,
		&c.UserID,
		&c.OrderID,
		&c.CouponPolicyID,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to fetch coupon by code", zap.String("coupon_code", code), zap.Error(err))
		return nil, coupon.ErrCouponNotFound
	}

	log.Info("fetched coupon successfully", zap.String("coupon_id", c.ID), zap.String("coupon_code", c.Code))
	return &c, nil
}

func (r *repository) UpdateCoupon(ctx context.Context, c *coupon.Coupon) (*coupon.Coupon, error) {
	ctx, span := instrument.StartSpan(ctx, "V2.Repository.UpdateCoupon")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	row := r.pg.Pool.QueryRow(ctx, `
		UPDATE coupons
		SET
			status = $1,
			used_at = $2,
			user_id = $3,
			order_id = $4,
			updated_at = NOW()
		WHERE id = $5
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
		c.Status,
		c.UsedAt,
		c.UserID,
		c.OrderID,
		c.ID,
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
		span.RecordError(err)
		log.Error("failed to update coupon", zap.String("coupon_id", c.ID), zap.Error(err))
		return nil, errors.New("failed to update coupon")
	}

	log.Info("coupon updated successfully", zap.String("coupon_id", result.ID), zap.String("coupon_code", result.Code), zap.String("status", string(result.Status)))
	return &result, nil
}

func (r *repository) FindCouponPolicyByID(ctx context.Context, id string) (*coupon.CouponPolicy, error) {
	ctx, span := instrument.StartSpan(ctx, "V2.Repository.FindCouponPolicyByID")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

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
		WHERE id = $1
		LIMIT 1
	`, id)

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
		span.RecordError(err)
		log.Error("failed to fetch coupon policy by code", zap.String("policy_id", id), zap.Error(err))
		return nil, coupon.ErrCouponPolicyNotFound
	}

	log.Info("fetched coupon policy successfully", zap.String("policy_id", policy.ID), zap.String("policy_code", policy.Code))
	return &policy, nil
}

func (r *repository) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := r.pg.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
