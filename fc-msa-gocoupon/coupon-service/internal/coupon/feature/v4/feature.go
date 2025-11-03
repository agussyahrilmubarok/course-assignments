package v4

import (
	"context"
	"errors"
	"time"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/internal/coupon/cache"
	"example.com/coupon/pkg/exception"
	"example.com/coupon/pkg/instrument"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=ICouponFeature
type ICouponFeature interface {
	IssueCoupon(ctx context.Context, couponPolicyCode string, userID string) error
	ProcessIssueCoupon(ctx context.Context, message coupon.IssueCouponMessage) error
	UseCoupon(ctx context.Context, couponID string, userID string, orderID string) (*coupon.Coupon, error)
	CancelCoupon(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error)
	FindCouponByCode(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error)
	FindCouponsByUserID(ctx context.Context, userID string) ([]coupon.Coupon, error)
	FindCouponsByCouponPolicyCode(ctx context.Context, couponPolicyCode string) ([]coupon.Coupon, error)
}

type couponFeature struct {
	db            *gorm.DB
	rdb           *redis.Client
	cache         cache.ICache
	kafkaProducer *KafkaProducer
	log           zerolog.Logger
	tracer        trace.Tracer
}

func NewCouponFeature(
	db *gorm.DB,
	rdb *redis.Client,
	cache cache.ICache,
	kafkaProducer *KafkaProducer,
	log zerolog.Logger,
	tracer trace.Tracer,
) ICouponFeature {
	return &couponFeature{
		db:            db,
		rdb:           rdb,
		cache:         cache,
		kafkaProducer: kafkaProducer,
		log:           log,
		tracer:        tracer,
	}
}

// IssueCoupon generates a new coupon for a given user under a specified coupon policy.
// It validates the policy period, checks quota limits, and persists the coupon to the database.
// Returns the issued coupon on success or an appropriate error if the operation fails.
// Issue List:
func (f *couponFeature) IssueCoupon(ctx context.Context, couponPolicyCode string, userID string) error {
	// Use a background context to avoid client-side cancellation
	bgCtx := context.Background()

	bgCtx, span := f.tracer.Start(bgCtx, "feature.IssueCouponNoContextCanceled",
		trace.WithAttributes(
			attribute.String("coupon.policy_code", couponPolicyCode),
			attribute.String("user.id", userID),
		),
	)
	defer span.End()

	log := instrument.GetLogger(bgCtx, f.log)
	lockKey := "coupon:lock:" + couponPolicyCode
	lockValue := uuid.NewString()

	// Step 1: Acquire distributed lock
	acquired, err := f.tryAcquireLock(bgCtx, lockKey, lockValue, 3*time.Second, 5*time.Second)
	if err != nil {
		span.RecordError(err)
		return exception.NewInternal("Failed to acquire Redis lock", err)
	}
	if !acquired {
		err := coupon.ErrCouponTooManyRequests
		span.RecordError(err)
		return exception.NewTooManyRequests("Please try again later", err)
	}
	defer func() {
		if err := f.releaseLock(bgCtx, lockKey, lockValue); err != nil {
			log.Warn().Err(err).Str("key", lockKey).Msg("Failed to release lock")
		}
	}()

	err = f.db.WithContext(bgCtx).Transaction(func(tx *gorm.DB) error {
		// Step 2: Fetch coupon policy
		var policy coupon.CouponPolicy
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("code = ?", couponPolicyCode).
			First(&policy).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Warn().Str("coupon_policy_code", couponPolicyCode).Msg("Coupon policy not found")
				return exception.NewNotFound("Coupon policy not found", err)
			}
			log.Error().Err(err).Str("coupon_policy_code", couponPolicyCode).Msg("Failed to load coupon policy")
			return exception.NewInternal("Failed to load coupon policy", err)
		}

		instrument.CouponQuota.WithLabelValues(policy.Code).Set(float64(policy.TotalQuantity))

		// Step 3: Validate period
		if !policy.IsValidPeriodUnix() {
			err := coupon.ErrCouponPolicyInvalidPeriod
			span.RecordError(err)
			log.Warn().
				Str("coupon_policy_code", policy.Code).
				Str("start_time", policy.StartTime.UTC().Format(time.RFC3339)).
				Str("end_time", policy.EndTime.UTC().Format(time.RFC3339)).
				Msg("Coupon policy is not valid in the current period")
			return exception.NewBadRequest("Coupon policy is not valid in current period", err)
		}

		// Step 4: Check if user already has this coupon
		var existing coupon.Coupon
		if err := tx.Where("coupon_policy_id = ? AND user_id = ?", policy.ID, userID).First(&existing).Error; err == nil {
			log.Warn().
				Str("user_id", userID).
				Str("coupon_policy_code", couponPolicyCode).
				Msg("User already has a coupon for this policy")
			return exception.NewBadRequest("User already has a coupon for this policy", nil)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			span.RecordError(err)
			log.Error().Err(err).Msg("Failed to check existing coupon")
			return exception.NewInternal("Failed to check existing coupon", err)
		}

		// Step 5: Decrement quota in Redis
		remaining, err := f.cache.DecrementAndGetCouponPolicyQuantity(bgCtx, policy.Code)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// Redis belum ada quantity → initialize
				if err := f.cache.SetCouponPolicyQuantity(bgCtx, policy.Code, int64(policy.TotalQuantity-1), policy.EndTime); err != nil {
					return exception.NewInternal("Failed to initialize Redis quantity", err)
				}
				remaining = int64(policy.TotalQuantity - 1)
			} else {
				span.RecordError(err)
				log.Error().Err(err).Msg("Failed to update coupon quota in Redis")
				return exception.NewInternal("Failed to update coupon quota", err)
			}
		}

		if remaining < 0 {
			// Rollback Redis decrement
			if _, err := f.cache.IncrementAndGetCouponPolicyQuantity(bgCtx, policy.Code); err != nil {
				log.Warn().Err(err).Msg("Failed to rollback Redis quota")
			}
			err := coupon.ErrCouponPolicyQoutaExceeded
			span.RecordError(err)
			log.Warn().Str("coupon_policy_code", policy.Code).Msg("Coupon policy quota exceeded")
			return exception.NewBadRequest("Coupon policy quota exceeded", err)
		}

		// Step 6: Publish issue coupon
		if err := f.kafkaProducer.SendCouponIssueRequest(bgCtx, coupon.IssueCouponMessage{
			CouponPolicyCode: policy.Code,
			UserID:           userID,
		}); err != nil {
			if _, err := f.cache.IncrementAndGetCouponPolicyQuantity(bgCtx, policy.Code); err != nil {
				log.Warn().Err(err).Msg("Failed to rollback Redis quota")
			}
			span.RecordError(err)
			log.Warn().Str("coupon_policy_code", policy.Code).Msg("Failed to process coupon issue request")
			return exception.NewInternal("Failed to process coupon issue request", err)
		}

		return nil
	})

	if err != nil {
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "Kafka message published successfully")
	log.Info().
		Str("coupon_policy_code", couponPolicyCode).
		Str("user_id", userID).
		Msg("Coupon issued successfully")
	return nil
}

func (f *couponFeature) ProcessIssueCoupon(ctx context.Context, message coupon.IssueCouponMessage) error {
	ctx, span := f.tracer.Start(ctx, "feature.ProcessIssueCoupon",
		trace.WithAttributes(
			attribute.String("coupon.policy_code", message.CouponPolicyCode),
			attribute.String("user.id", message.UserID),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, f.log)

	return f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Fetch policy
		var policy coupon.CouponPolicy
		if err := tx.Where("code = ?", message.CouponPolicyCode).First(&policy).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Warn().Str("coupon_policy_code", message.CouponPolicyCode).Msg("Policy not found during processing")
				return exception.NewNotFound("Coupon policy not found", err)
			}
			return exception.NewInternal("Failed to load policy during processing", err)
		}

		// Step 2: Create coupon
		newCoupon := coupon.Coupon{
			ID:             uuid.NewString(),
			Code:           uuid.NewString(),
			Status:         coupon.CouponStatusAvailable,
			UserID:         message.UserID,
			CouponPolicyID: policy.ID,
		}
		if err := tx.Create(&newCoupon).Error; err != nil {
			// Rollback Redis decrement
			if _, err := f.cache.IncrementAndGetCouponPolicyQuantity(ctx, policy.Code); err != nil {
				log.Warn().Err(err).Msg("Failed to rollback Redis quota")
			}
			span.RecordError(err)
			return exception.NewInternal("Failed to create coupon", err)
		}

		if err := f.cache.SetCouponState(ctx, newCoupon, policy.EndTime); err != nil {
			log.Warn().Err(err).Str("coupon_code", newCoupon.Code).Msg("Failed to cache coupon state")
			// Not fatal
		}

		instrument.CouponIssued.WithLabelValues(policy.Code).Inc()

		span.SetStatus(codes.Ok, "Coupon issued successfully from Kafka")
		log.Info().
			Str("coupon_policy_code", message.CouponPolicyCode).
			Str("user_id", message.UserID).
			Str("coupon_code", newCoupon.Code).
			Msg("Coupon issued successfully from Kafka message")
		return nil
	})
}

// UseCoupon marks a specific coupon as used for a given order by a user.
// It ensures the coupon exists, belongs to the user, and is in a valid state to be used.
// Returns the updated coupon or an error if usage fails or cannot be saved.
func (f *couponFeature) UseCoupon(ctx context.Context, couponCode string, userID string, orderID string) (*coupon.Coupon, error) {
	ctx, span := f.tracer.Start(ctx, "feature.UseCoupon",
		trace.WithAttributes(
			attribute.String("coupon.code", couponCode),
			attribute.String("user.id", userID),
			attribute.String("order.id", orderID),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, f.log)

	var coupon coupon.Coupon
	if err := f.db.WithContext(ctx).
		Preload("CouponPolicy").
		First(&coupon, "code = ? AND user_id = ?", couponCode, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.RecordError(err)
			log.Warn().
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Msg("Coupon not found")
			return nil, exception.NewNotFound("Coupon not found", err)
		}
		log.Error().
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to fetch coupon")
		return nil, exception.NewInternal("Failed to fetch coupon", err)
	}

	if err := coupon.Use(orderID); err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to use coupon")
		return nil, exception.NewBadRequest("Failed to use coupon", err)
	}

	if err := f.db.WithContext(ctx).Save(&coupon).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Str("order_id", orderID).
			Err(err).
			Msg("Failed to save use coupon")
		return nil, exception.NewInternal("Failed to save use coupon", err)
	}

	if err := f.cache.SetCouponState(ctx, coupon, coupon.CouponPolicy.EndTime); err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Str("order_id", orderID).
			Err(err).
			Msg("Failed to update coupon state")
		return nil, exception.NewInternal("Failed to save use coupon", err)
	}

	log.Info().
		Str("coupon_code", coupon.Code).
		Str("coupon_status", string(coupon.Status)).
		Str("user_id", userID).
		Str("order_id", orderID).
		Msg("Coupon used successfully")
	return &coupon, nil
}

// CancelCoupon reverses the usage of a coupon, marking it as available again.
// It validates the coupon's existence and current status before performing the cancellation.
// Returns the updated coupon or an error if cancellation fails or cannot be saved.
func (f *couponFeature) CancelCoupon(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := f.tracer.Start(ctx, "feature.CancelCoupon",
		trace.WithAttributes(
			attribute.String("coupon.code", couponCode),
			attribute.String("user.id", userID),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, f.log)

	var coupon coupon.Coupon
	if err := f.db.WithContext(ctx).
		Preload("CouponPolicy").
		First(&coupon, "code = ? AND user_id = ?", couponCode, userID).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Msg("Coupon not found")
			return nil, exception.NewNotFound("Coupon not found", err)
		}
		log.Error().
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to fetch coupon")
		return nil, exception.NewInternal("Failed to fetch coupon", err)
	}

	if err := coupon.Cancel(); err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to cancel coupon")
		return nil, exception.NewBadRequest("Failed to cancel coupon", err)
	}

	if err := f.db.WithContext(ctx).Save(&coupon).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to save use coupon")
		return nil, exception.NewInternal("Failed to save cancel coupon", err)
	}

	if _, err := f.cache.IncrementAndGetCouponPolicyQuantity(ctx, coupon.CouponPolicy.Code); err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to increment coupon policy quantity (rollback)")
		return nil, exception.NewInternal("Failed to save cancel coupon", err)
	}

	if err := f.cache.SetCouponState(ctx, coupon, coupon.CouponPolicy.EndTime); err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to set coupon state")
		return nil, exception.NewInternal("Failed to save cancel coupon", err)
	}

	log.Info().
		Str("coupon_code", coupon.Code).
		Str("coupon_status", string(coupon.Status)).
		Str("user_id", userID).
		Msg("Coupon used successfully")
	return &coupon, nil
}

// FindCouponByCode retrieves a single coupon for a user by its unique code.
// Returns the coupon if found, or a NotFound error if no matching coupon exists.
func (f *couponFeature) FindCouponByCode(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := f.tracer.Start(ctx, "feature.FindCouponByCode",
		trace.WithAttributes(
			attribute.String("coupon.code", couponCode),
			attribute.String("user.id", userID),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, f.log)

	var c coupon.Coupon
	if err := f.db.WithContext(ctx).
		Preload("CouponPolicy").
		First(&c, "code = ? AND user_id = ?", couponCode, userID).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Msg("Coupon not found")
			return nil, exception.NewNotFound("Coupon not found", err)
		}
		log.Error().
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to fetch coupon")
		return nil, exception.NewInternal("Failed to fetch coupon", err)
	}

	log.Info().
		Str("coupon_code", c.Code).
		Str("user_id", userID).
		Msg("Coupon found successfully")
	return &c, nil
}

// FindCouponsByUserID fetches all coupons associated with a specific user.
// Returns the list of coupons or an Internal error if the database query fails.
func (f *couponFeature) FindCouponsByUserID(ctx context.Context, userID string) ([]coupon.Coupon, error) {
	ctx, span := f.tracer.Start(ctx, "feature.FindCouponsByUserID",
		trace.WithAttributes(attribute.String("user.id", userID)),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, f.log)

	var coupons []coupon.Coupon
	if err := f.db.WithContext(ctx).
		Preload("CouponPolicy").
		Where("user_id = ?", userID).
		Find(&coupons).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("user_id", userID).
			Err(err).
			Msg("Failed to fetch coupons by user")
		return nil, exception.NewInternal("Failed to fetch coupons by user", err)
	}

	log.Info().
		Str("user_id", userID).
		Int("coupon_count", len(coupons)).
		Msg("Fetched coupons for user successfully")
	return coupons, nil
}

// FindCouponsByCouponPolicyCode fetches all coupons issued under a specific coupon policy.
// Returns the list of coupons or an appropriate error if the policy does not exist or the query fails.
func (f *couponFeature) FindCouponsByCouponPolicyCode(ctx context.Context, couponPolicyCode string) ([]coupon.Coupon, error) {
	ctx, span := f.tracer.Start(ctx, "feature.FindCouponsByCouponPolicyCode",
		trace.WithAttributes(attribute.String("coupon.policy_code", couponPolicyCode)),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, f.log)

	var policy coupon.CouponPolicy
	if err := f.db.WithContext(ctx).
		Where("code = ?", couponPolicyCode).
		First(&policy).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("coupon_policy_code", couponPolicyCode).
				Msg("Coupon policy not found")
			return nil, exception.NewNotFound("Coupon policy not found", err)
		}
		log.Error().
			Str("coupon_policy_code", couponPolicyCode).
			Err(err).
			Msg("Failed to fetch coupon policy")
		return nil, exception.NewInternal("Failed to fetch coupon policy", err)
	}

	var coupons []coupon.Coupon
	if err := f.db.WithContext(ctx).
		Where("coupon_policy_id = ?", policy.ID).
		Find(&coupons).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_policy_code", couponPolicyCode).
			Err(err).
			Msg("Failed to fetch coupons by policy code")
		return nil, exception.NewInternal("Failed to fetch coupons by policy code", err)
	}

	instrument.CouponQuota.WithLabelValues(policy.Code).Set(float64(policy.TotalQuantity))
	instrument.CouponIssued.WithLabelValues(policy.Code).Set(float64(len(coupons)))

	log.Info().
		Str("coupon_policy_code", couponPolicyCode).
		Int("coupon_count", len(coupons)).
		Msg("Fetched coupons for policy successfully")
	return coupons, nil
}

func (s *couponFeature) tryAcquireLock(ctx context.Context, key, value string, wait, lease time.Duration) (bool, error) {
	deadline := time.Now().Add(wait)
	for time.Now().Before(deadline) {
		ok, err := s.rdb.SetNX(ctx, key, value, lease).Result()
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false, nil
}

func (s *couponFeature) releaseLock(ctx context.Context, key, value string) error {
	script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
	return s.rdb.Eval(ctx, script, []string{key}, value).Err()
}
