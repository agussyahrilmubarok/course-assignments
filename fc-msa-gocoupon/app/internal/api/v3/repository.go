package v3

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"example.com/coupon-service/internal/config"
	"example.com/coupon-service/internal/coupon"
	"example.com/coupon-service/internal/instrument"
	"example.com/coupon-service/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type IRepository interface {
	FindCouponPolicyByCode(ctx context.Context, code string) (*coupon.CouponPolicy, error)
	CountIssuedCoupons(ctx context.Context, policyID string) (int, error)
	CreateCoupon(ctx context.Context, coupon *coupon.Coupon) (*coupon.Coupon, error)
	FindCouponByCode(ctx context.Context, code string) (*coupon.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon *coupon.Coupon) (*coupon.Coupon, error)
	FindCouponPolicyByID(ctx context.Context, id string) (*coupon.CouponPolicy, error)
	SetCouponPolicyCache(ctx context.Context, policy *coupon.CouponPolicy) error
	GetCouponPolicyCache(ctx context.Context, code string) (*coupon.CouponPolicy, error)
	SetCouponPolicyIssuedCache(ctx context.Context, code string, quantity int, endTime time.Time) error
	GetCouponPolicyIssuedCache(ctx context.Context, code string) (int, error)
	IncrCouponPolicyIssuedCache(ctx context.Context, code string) error
	DecrCouponPolicyIssuedCache(ctx context.Context, code string) error
	AcquireRedisLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	ReleaseRedisLock(ctx context.Context, key string) error
}

var (
	CouponPolicyKeyPrefix            = "coupon:policy:"
	CouponPolicyIssuedCountKeyPrefix = "coupon:policy:issued:"
)

type repository struct {
	pg  *config.Postgres
	rdb *config.Redis
}

func NewRepository(pg *config.Postgres, rdb *config.Redis) IRepository {
	return &repository{
		pg:  pg,
		rdb: rdb,
	}
}

func (r *repository) FindCouponPolicyByCode(ctx context.Context, code string) (*coupon.CouponPolicy, error) {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.FindCouponPolicyByCode")
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
		span.RecordError(err)
		log.Error("failed to fetch coupon policy by code", zap.String("policy_code", code), zap.Error(err))
		return nil, coupon.ErrCouponPolicyNotFound
	}

	log.Info("fetched coupon policy successfully", zap.String("policy_id", policy.ID), zap.String("policy_code", code))
	return &policy, nil
}

func (r *repository) CountIssuedCoupons(ctx context.Context, policyID string) (int, error) {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.CountIssuedCoupons")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	row := r.pg.Pool.QueryRow(ctx, `
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

func (r *repository) CreateCoupon(ctx context.Context, c *coupon.Coupon) (*coupon.Coupon, error) {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.CreateCoupon")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

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
		span.RecordError(err)
		log.Error("failed to create coupon", zap.Error(err))
		return nil, errors.New("")
	}

	log.Info("coupon created successfully", zap.String("coupon_id", result.ID), zap.String("coupon_code", result.Code))
	return &result, nil
}

func (r *repository) FindCouponByCode(ctx context.Context, code string) (*coupon.Coupon, error) {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.FindCouponByCode")
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
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.UpdateCoupon")
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
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.FindCouponPolicyByID")
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

func (r *repository) SetCouponPolicyCache(ctx context.Context, policy *coupon.CouponPolicy) error {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.SetCouponPolicyCache")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	key := CouponPolicyKeyPrefix + policy.Code
	ttl := time.Until(policy.EndTime)
	if ttl <= 0 {
		ttl = time.Millisecond
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to marshal policy", zap.String("policy_code", policy.Code), zap.Error(err))
	}

	if err := r.rdb.Client.Set(ctx, key, policyJSON, ttl).Err(); err != nil {
		span.RecordError(err)
		log.Error("failed to set coupon policy", zap.String("policy_code", policy.Code), zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) GetCouponPolicyCache(ctx context.Context, policyCode string) (*coupon.CouponPolicy, error) {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.GetCouponPolicyCache")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	key := CouponPolicyKeyPrefix + policyCode
	val, err := r.rdb.Client.Get(ctx, key).Result()
	if err != nil {
		span.RecordError(err)
		log.Error("failed to get coupon policy", zap.String("policy_code", policyCode), zap.Error(err))
		return nil, err
	}

	var policy coupon.CouponPolicy
	err = json.Unmarshal([]byte(val), &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (r *repository) SetCouponPolicyIssuedCache(ctx context.Context, code string, quantity int, endTime time.Time) error {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.SetCouponPolicyIssuedCache")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	key := CouponPolicyIssuedCountKeyPrefix + code
	ttl := time.Until(endTime)
	if ttl <= 0 {
		ttl = time.Millisecond
	}

	if err := r.rdb.Client.Set(ctx, key, quantity, ttl).Err(); err != nil {
		span.RecordError(err)
		log.Error("failed to set issued coupon count in cache", zap.String("policy_code", code), zap.Int("quantity", quantity), zap.Error(err))
		return err
	}

	log.Info("set issued coupon count successfully", zap.String("policy_code", code), zap.Int("quantity", quantity))
	return nil
}

func (r *repository) GetCouponPolicyIssuedCache(ctx context.Context, code string) (int, error) {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.GetCouponPolicyIssuedCache")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	key := CouponPolicyIssuedCountKeyPrefix + code
	val, err := r.rdb.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // belum ada di cache
		}
		span.RecordError(err)
		log.Error("failed to get issued coupon count from cache", zap.String("policy_code", code), zap.Error(err))
		return 0, err
	}

	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to parse issued coupon count", zap.String("policy_code", code), zap.String("value", val), zap.Error(err))
		return 0, err
	}

	return int(count), nil
}

func (r *repository) IncrCouponPolicyIssuedCache(ctx context.Context, code string) error {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.IncrCouponPolicyIssuedCache")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	key := CouponPolicyIssuedCountKeyPrefix + code
	newVal, err := r.rdb.Client.Incr(ctx, key).Result()
	if err != nil {
		span.RecordError(err)
		log.Error("failed to increment issued coupon count", zap.String("policy_code", code), zap.Error(err))
		return err
	}

	log.Info("incremented issued coupon count", zap.String("policy_code", code), zap.Int64("new_value", newVal))
	return nil
}

func (r *repository) DecrCouponPolicyIssuedCache(ctx context.Context, code string) error {
	ctx, span := instrument.StartSpan(ctx, "V3.Repository.DecrCouponPolicyIssuedCache")
	defer span.End()

	log := logger.GetLoggerFromContext(ctx)

	key := CouponPolicyIssuedCountKeyPrefix + code
	newVal, err := r.rdb.Client.Decr(ctx, key).Result()
	if err != nil {
		span.RecordError(err)
		log.Error("failed to decrement issued coupon count", zap.String("policy_code", code), zap.Error(err))
		return err
	}

	log.Info("decremented issued coupon count", zap.String("policy_code", code), zap.Int64("new_value", newVal))
	return nil
}

func (r *repository) AcquireRedisLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	ok, err := r.rdb.Client.SetNX(ctx, key, "locked", ttl).Result()
	return ok, err
}

func (r *repository) ReleaseRedisLock(ctx context.Context, key string) error {
	_, err := r.rdb.Client.Del(ctx, key).Result()
	return err
}
