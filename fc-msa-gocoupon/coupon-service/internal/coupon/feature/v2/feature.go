package v2

import (
	"context"
	"errors"
	"time"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/exception"
	"example.com/coupon/pkg/instrument"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockery --name=ICouponFeature
type ICouponFeature interface {
	IssueCoupon(ctx context.Context, couponPolicyCode string, userID string) (*coupon.Coupon, error)
	IssueCouponNoContextCanceled(ctx context.Context, couponPolicyCode string, userID string) (*coupon.Coupon, error)
	UseCoupon(ctx context.Context, couponID string, userID string, orderID string) (*coupon.Coupon, error)
	CancelCoupon(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error)
	FindCouponByCode(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error)
	FindCouponsByUserID(ctx context.Context, userID string) ([]coupon.Coupon, error)
	FindCouponsByCouponPolicyCode(ctx context.Context, couponPolicyCode string) ([]coupon.Coupon, error)
}

type couponFeature struct {
	db     *gorm.DB
	log    zerolog.Logger
	tracer trace.Tracer
}

func NewCouponFeature(db *gorm.DB, log zerolog.Logger, tracer trace.Tracer) ICouponFeature {
	return &couponFeature{
		db:     db,
		log:    log,
		tracer: tracer,
	}
}

// IssueCoupon generates a new coupon for a given user under a specified coupon policy.
// It validates the policy period, checks quota limits, and persists the coupon to the database.
// Returns the issued coupon on success or an appropriate error if the operation fails.
// Issued List:
//
// **ContextCanceledUnderHighLoad** (observed issue under load testing):
//
//	When many concurrent requests occur, transactions may wait for locks,
//	causing client timeouts and “context canceled” errors.
//	➜ Solution: This is expected under contention; can be mitigated with
//	  increased DB connection pool, per-request timeout, or retry logic.
func (f *couponFeature) IssueCoupon(ctx context.Context, couponPolicyCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := f.tracer.Start(ctx, "feature.IssueCoupon",
		trace.WithAttributes(
			attribute.String("coupon.policy_code", couponPolicyCode),
			attribute.String("user.id", userID),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, f.log)

	var issuedCoupon *coupon.Coupon
	err := f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the coupon policy to prevent concurrent quota violations
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

		// Check policy validity period (UTC-safe)
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

		// Check if user already has a coupon for this policy
		var existing coupon.Coupon
		if err := tx.Where("coupon_policy_id = ? AND user_id = ?", policy.ID, userID).
			First(&existing).Error; err == nil {
			span.RecordError(err)
			log.Warn().
				Str("user_id", userID).
				Str("coupon_policy_code", couponPolicyCode).
				Msg("User already has a coupon for this policy")
			return exception.NewBadRequest("User already has a coupon for this policy", nil)
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			span.RecordError(err)
			log.Error().
				Str("user_id", userID).
				Err(err).
				Msg("Failed to check existing coupon")
			return exception.NewInternal("Failed to check existing coupon", err)
		}

		// Check current issued count (real-time, without Preload)
		var issuedCount int64
		if err := tx.Model(&coupon.Coupon{}).
			Where("coupon_policy_id = ?", policy.ID).
			Count(&issuedCount).Error; err != nil {
			span.RecordError(err)
			log.Error().Err(err).Str("coupon_policy_code", policy.Code).Msg("Failed to count issued coupons")
			return exception.NewInternal("Failed to count issued coupons", err)
		}

		if issuedCount >= int64(policy.TotalQuantity) {
			err := coupon.ErrCouponPolicyQoutaExceeded
			span.RecordError(err)
			log.Warn().
				Str("coupon_policy_code", policy.Code).
				Int("total_quantity", policy.TotalQuantity).
				Int64("issued_quantity", issuedCount).
				Msg("Coupon policy quota exceeded")
			return exception.NewBadRequest("Coupon policy quota exceeded", err)
		}

		newCoupon := coupon.Coupon{
			ID:             uuid.NewString(),
			Code:           uuid.NewString(),
			Status:         coupon.CouponStatusAvailable,
			UserID:         userID,
			CouponPolicyID: policy.ID,
		}
		if err := f.db.WithContext(ctx).Create(&newCoupon).Error; err != nil {
			span.RecordError(err)
			log.Error().
				Str("coupon_policy_code", policy.Code).
				Str("user_id", userID).
				Err(err).
				Msg("Failed to issue new coupon")
			return exception.NewInternal("Failed to issue coupon", err)
		}

		instrument.CouponIssued.WithLabelValues(policy.Code).Inc()

		issuedCoupon = &newCoupon
		return nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "IssueCoupon")
	span.SetAttributes(attribute.String("coupon.code", issuedCoupon.Code))
	log.Info().
		Str("coupon_policy_code", couponPolicyCode).
		Str("coupon_code", issuedCoupon.Code).
		Str("user_id", userID).
		Msg("Coupon issued successfully")
	return issuedCoupon, nil
}

// IssueCouponNoContextCanceled issues a coupon without being affected by client-side context cancellation.
// This prevents "context canceled" errors under high concurrent load scenarios.
// Warning: Use carefully, as long-running transactions won't be canceled even if the client disconnects.
func (f *couponFeature) IssueCouponNoContextCanceled(ctx context.Context, couponPolicyCode string, userID string) (*coupon.Coupon, error) {
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

	var issuedCoupon *coupon.Coupon
	err := f.db.WithContext(bgCtx).Transaction(func(tx *gorm.DB) error {
		// Lock the coupon policy to prevent concurrent quota violations
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

		// Check policy validity period
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

		// Check if user already has a coupon for this policy
		var existing coupon.Coupon
		if err := tx.Where("coupon_policy_id = ? AND user_id = ?", policy.ID, userID).First(&existing).Error; err == nil {
			span.RecordError(err)
			log.Warn().
				Str("user_id", userID).
				Str("coupon_policy_code", couponPolicyCode).
				Msg("User already has a coupon for this policy")
			return exception.NewBadRequest("User already has a coupon for this policy", nil)
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			span.RecordError(err)
			log.Error().
				Str("user_id", userID).
				Err(err).
				Msg("Failed to check existing coupon")
			return exception.NewInternal("Failed to check existing coupon", err)
		}

		// Check current issued count
		var issuedCount int64
		if err := tx.Model(&coupon.Coupon{}).Where("coupon_policy_id = ?", policy.ID).Count(&issuedCount).Error; err != nil {
			span.RecordError(err)
			log.Error().Err(err).Str("coupon_policy_code", policy.Code).Msg("Failed to count issued coupons")
			return exception.NewInternal("Failed to count issued coupons", err)
		}

		if issuedCount >= int64(policy.TotalQuantity) {
			err := coupon.ErrCouponPolicyQoutaExceeded
			span.RecordError(err)
			log.Warn().
				Str("coupon_policy_code", policy.Code).
				Int("total_quantity", policy.TotalQuantity).
				Int64("issued_quantity", issuedCount).
				Msg("Coupon policy quota exceeded")
			return exception.NewBadRequest("Coupon policy quota exceeded", err)
		}

		newCoupon := coupon.Coupon{
			ID:             uuid.NewString(),
			Code:           uuid.NewString(),
			Status:         coupon.CouponStatusAvailable,
			UserID:         userID,
			CouponPolicyID: policy.ID,
		}
		if err := tx.Create(&newCoupon).Error; err != nil {
			span.RecordError(err)
			log.Error().
				Str("coupon_policy_code", policy.Code).
				Str("user_id", userID).
				Err(err).
				Msg("Failed to issue new coupon")
			return exception.NewInternal("Failed to issue coupon", err)
		}

		instrument.CouponIssued.WithLabelValues(policy.Code).Inc()

		issuedCoupon = &newCoupon
		return nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "IssueCouponNoContextCanceled")
	span.SetAttributes(attribute.String("coupon.code", issuedCoupon.Code))
	log.Info().
		Str("coupon_policy_code", couponPolicyCode).
		Str("coupon_code", issuedCoupon.Code).
		Str("user_id", userID).
		Msg("Coupon issued successfully (no context canceled)")
	return issuedCoupon, nil
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

	span.SetStatus(codes.Ok, "UseCoupon")
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
		return nil, exception.NewInternal("Failed to save use coupon", err)
	}

	span.SetStatus(codes.Ok, "CancelCoupon")
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

	span.SetStatus(codes.Ok, "FindCouponByCode")
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

	span.SetStatus(codes.Ok, "FindCouponsByUserID")
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

	span.SetStatus(codes.Ok, "FindCouponsByCouponPolicyCode")
	log.Info().
		Str("coupon_policy_code", couponPolicyCode).
		Int("coupon_count", len(coupons)).
		Msg("Fetched coupons for policy successfully")
	return coupons, nil
}
