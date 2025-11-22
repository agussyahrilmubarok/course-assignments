package cache

import (
	"context"
	"encoding/json"
	"time"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/instrument"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

//go:generate mockery --name=ICache
type ICache interface {
	SetCouponPolicyQuantity(ctx context.Context, couponPolicyCode string, quantity int64, endTime time.Time) error
	GetCouponPolicyQuantity(ctx context.Context, couponPolicyCode string) (int64, error)
	DeleteCouponPolicyQuantity(ctx context.Context, couponPolicyCode string) error
	DecrementAndGetCouponPolicyQuantity(ctx context.Context, couponPolicyCode string) (int64, error)
	IncrementAndGetCouponPolicyQuantity(ctx context.Context, couponPolicyCode string) (int64, error)

	SetCouponState(ctx context.Context, coupon coupon.Coupon, endTime time.Time) error
	GetCouponState(ctx context.Context, couponCode string) (*coupon.Coupon, error)
	DeleteCouponState(ctx context.Context, couponCode string) error
}

type cache struct {
	rdb    *redis.Client
	log    zerolog.Logger
	tracer trace.Tracer
}

func NewCache(
	rdb *redis.Client,
	log zerolog.Logger,
	tracer trace.Tracer,
) ICache {
	return &cache{
		rdb:    rdb,
		log:    log,
		tracer: tracer,
	}
}

func (c *cache) SetCouponPolicyQuantity(ctx context.Context, couponPolicyCode string, quantity int64, endTime time.Time) error {
	ctx, span := c.tracer.Start(ctx, "cache.SetCouponPolicyQuantity",
		trace.WithAttributes(
			attribute.String("coupon.policy_code", couponPolicyCode),
			attribute.Int64("coupon.policy_quantity", quantity),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, c.log)
	log = log.With().Str("func", "SetCouponPolicyQuantity").Logger()

	key := "coupon:policy:quantity:" + couponPolicyCode
	ttl := time.Until(endTime)
	if ttl <= 0 {
		ttl = 1
	}

	if err := c.rdb.Set(ctx, key, quantity, ttl).Err(); err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("key", key).Msg("failed to set coupon policy quantity")
		return err
	}

	return nil
}

func (c *cache) GetCouponPolicyQuantity(ctx context.Context, couponPolicyCode string) (int64, error) {
	ctx, span := c.tracer.Start(ctx, "cache.GetCouponPolicyQuantity",
		trace.WithAttributes(
			attribute.String("coupon.policy_code", couponPolicyCode),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, c.log)
	log = log.With().Str("func", "GetCouponPolicyQuantity").Logger()

	key := "coupon:policy:quantity:" + couponPolicyCode
	val, err := c.rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("key", key).Msg("failed to get coupon policy quantity")
		return 0, err
	}

	return val, nil
}

func (c *cache) DeleteCouponPolicyQuantity(ctx context.Context, couponPolicyCode string) error {
	ctx, span := c.tracer.Start(ctx, "cache.DeleteCouponPolicyQuantity",
		trace.WithAttributes(attribute.String("coupon.policy_code", couponPolicyCode)),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, c.log)
	log = log.With().Str("func", "DeleteCouponPolicyQuantity").Logger()

	key := "coupon:policy:quantity:" + couponPolicyCode
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("key", key).Msg("failed to delete coupon policy quantity")
		return err
	}

	return nil
}

func (c *cache) DecrementAndGetCouponPolicyQuantity(ctx context.Context, couponPolicyCode string) (int64, error) {
	ctx, span := c.tracer.Start(ctx, "cache.DecrementAndGetCouponPolicyQuantity",
		trace.WithAttributes(attribute.String("coupon.policy_code", couponPolicyCode)),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, c.log)
	log = log.With().Str("func", "DecrementAndGetCouponPolicyQuantity").Logger()

	key := "coupon:policy:quantity:" + couponPolicyCode
	val, err := c.rdb.Decr(ctx, key).Result()
	if err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("key", key).Msg("failed to decrement and get coupon policy quantity")
		return 0, err
	}

	return val, nil
}

func (c *cache) IncrementAndGetCouponPolicyQuantity(ctx context.Context, couponPolicyCode string) (int64, error) {
	ctx, span := c.tracer.Start(ctx, "cache.IncrementAndGetCouponPolicyQuantity",
		trace.WithAttributes(attribute.String("coupon.policy_code", couponPolicyCode)),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, c.log)
	log = log.With().Str("func", "IncrementAndGetCouponPolicyQuantity").Logger()

	key := "coupon:policy:quantity:" + couponPolicyCode
	val, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("key", key).Msg("failed to increment and get coupon policy quantity")
		return 0, err
	}

	return val, nil
}

func (c *cache) SetCouponState(ctx context.Context, cp coupon.Coupon, endTime time.Time) error {
	ctx, span := c.tracer.Start(ctx, "cache.SetCouponState",
		trace.WithAttributes(attribute.String("coupon.code", cp.Code)),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, c.log)
	log = log.With().Str("func", "SetCouponState").Logger()

	key := "coupon:state:" + cp.Code
	ttl := time.Until(endTime)
	if ttl <= 0 {
		ttl = 1
	}

	// Encode coupon to JSON
	data, err := json.Marshal(cp)
	if err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("coupon_code", cp.Code).Msg("failed to marshal coupon state")
		return err
	}

	if err := c.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("key", key).Msg("failed to set coupon state")
		return err
	}

	log.Debug().Str("coupon_code", cp.Code).Msg("coupon state cached successfully")
	return nil
}

func (c *cache) GetCouponState(ctx context.Context, couponCode string) (*coupon.Coupon, error) {
	ctx, span := c.tracer.Start(ctx, "cache.GetCouponState",
		trace.WithAttributes(attribute.String("coupon.code", couponCode)),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, c.log)
	log = log.With().Str("func", "GetCouponState").Logger()

	key := "coupon:state:" + couponCode
	val, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // not found
	}
	if err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("key", key).Msg("failed to get coupon state")
		return nil, err
	}

	var cp coupon.Coupon
	if err := json.Unmarshal([]byte(val), &cp); err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("coupon_code", couponCode).Msg("failed to unmarshal coupon state")
		return nil, err
	}

	return &cp, nil
}

func (c *cache) DeleteCouponState(ctx context.Context, couponCode string) error {
	ctx, span := c.tracer.Start(ctx, "cache.DeleteCouponState",
		trace.WithAttributes(attribute.String("coupon.code", couponCode)),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, c.log)
	log = log.With().Str("func", "DeleteCouponState").Logger()

	key := "coupon:state:" + couponCode
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		span.RecordError(err)
		log.Error().Err(err).Str("key", key).Msg("failed to delete coupon state")
		return err
	}

	log.Debug().Str("coupon_code", couponCode).Msg("coupon state deleted successfully")
	return nil
}
