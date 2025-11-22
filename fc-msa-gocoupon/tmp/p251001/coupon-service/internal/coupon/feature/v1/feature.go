package v1

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
)

//go:generate mockery --name=ICouponFeature
type ICouponFeature interface {
	IssueCoupon(ctx context.Context, couponPolicyCode string, userID string) (*coupon.Coupon, error)
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
// Issue List:
// **RaceConditionOnQuotaCheck**: The function checks coupon quota (GetIssuedQuantity() >= TotalQuantity) before inserting a new record, but without locking or a transaction. Multiple concurrent requests could issue more coupons than the quota allows.
// **NoUserDuplicateCheck**: The function doesn’t verify if the user already has a coupon under the same policy, allowing multiple coupons per user when the business rule may only allow one.
// **HeavyPreloadCoupons**: The code uses Preload("Coupons"), which loads all coupons under the policy just to check issued count — inefficient for large datasets.
// **TruncatedUUIDCollisionRisk**: Using only the first 8 characters of a UUID for the coupon code increases the probability of collisions.
// **NoDatabaseTransaction**: The code performs multiple dependent operations (read policy → check quota → insert coupon) without wrapping them in a transaction, risking data inconsistency if an error occurs mid-process.
// **NonRealtimeIssuedQuantity**: The method GetIssuedQuantity() likely depends on preloaded data, which may not reflect the actual count in the database at query time.
func (f *couponFeature) IssueCoupon(ctx context.Context, couponPolicyCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := f.tracer.Start(ctx, "feature.IssueCoupon",
		trace.WithAttributes(
			attribute.String("coupon.policy_code", couponPolicyCode),
			attribute.String("user.id", userID),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, f.log)
	log = log.With().Str("func", "IssueCoupon").Logger()

	var couponPolicy coupon.CouponPolicy
	if err := f.db.WithContext(ctx).
		Preload("Coupons").
		Where("code = ?", couponPolicyCode).
		First(&couponPolicy).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("coupon_policy_code", couponPolicyCode).
				Msg("coupon policy not found")
			return nil, exception.NewNotFound("coupon policy not found", err)
		}
		log.Error().
			Str("coupon_policy_code", couponPolicyCode).
			Err(err).
			Msg("Failed to get coupon policy")
		return nil, exception.NewInternal("failed to get coupon policy", err)
	}

	if !couponPolicy.IsValidPeriodUnix() {
		err := coupon.ErrCouponPolicyInvalidPeriod
		span.RecordError(err)
		log.Warn().
			Str("coupon_policy_code", couponPolicy.Code).
			Str("coupon_policy_start_time", couponPolicy.StartTime.Format(time.RFC3339)).
			Str("coupon_policy_end_time", couponPolicy.EndTime.Format(time.RFC3339)).
			Msg("coupon policy is not valid in the current period")
		return nil, exception.NewBadRequest("coupon policy is not valid in current period", err)
	}

	if couponPolicy.GetIssuedQuantity() >= couponPolicy.TotalQuantity {
		err := coupon.ErrCouponPolicyQoutaExceeded
		span.RecordError(err)
		log.Warn().
			Str("coupon_policy_code", couponPolicy.Code).
			Int("coupon_policy_total_quantity", couponPolicy.TotalQuantity).
			Int("coupon_policy_issued_quantity", couponPolicy.GetIssuedQuantity()).
			Msg("coupon policy quota exceeded")
		return nil, exception.NewBadRequest("coupon policy quota exceeded", err)
	}

	newCoupon := coupon.Coupon{
		ID:             uuid.NewString(),
		Code:           uuid.NewString(),
		Status:         coupon.CouponStatusAvailable,
		UserID:         userID,
		CouponPolicyID: couponPolicy.ID,
	}
	if err := f.db.WithContext(ctx).Create(&newCoupon).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_policy_code", couponPolicy.Code).
			Str("user_id", userID).
			Err(err).
			Msg("Failed to issue new coupon")
		return nil, exception.NewInternal("failed to issue coupon", err)
	}

	instrument.CouponQuota.WithLabelValues(couponPolicy.Code).Set(float64(couponPolicy.TotalQuantity))
	instrument.CouponIssued.WithLabelValues(couponPolicy.Code).Inc()

	span.SetStatus(codes.Ok, "IssueCoupon")
	span.SetAttributes(attribute.String("coupon.code", newCoupon.Code))
	log.Info().
		Str("coupon_policy_code", couponPolicyCode).
		Str("coupon_code", newCoupon.Code).
		Str("user_id", userID).
		Msg("coupon issued successfully")
	return &newCoupon, nil
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
	log = log.With().Str("func", "UseCoupon").Logger()

	var coupon coupon.Coupon
	if err := f.db.WithContext(ctx).
		Preload("CouponPolicy").
		First(&coupon, "code = ? AND user_id = ?", couponCode, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.RecordError(err)
			log.Warn().
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Msg("coupon not found")
			return nil, exception.NewNotFound("coupon not found", err)
		}
		log.Error().
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Err(err).
			Msg("failed to fetch coupon")
		return nil, exception.NewInternal("failed to fetch coupon", err)
	}

	if err := coupon.Use(orderID); err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Err(err).
			Msg("failed to use coupon")
		return nil, exception.NewBadRequest("failed to use coupon", err)
	}

	if err := f.db.WithContext(ctx).Save(&coupon).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Str("order_id", orderID).
			Err(err).
			Msg("failed to save use coupon")
		return nil, exception.NewInternal("failed to save use coupon", err)
	}

	span.SetStatus(codes.Ok, "UseCoupon")
	log.Info().
		Str("coupon_code", coupon.Code).
		Str("coupon_status", string(coupon.Status)).
		Str("user_id", userID).
		Str("order_id", orderID).
		Msg("coupon used successfully")
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
	log = log.With().Str("func", "CancelCoupon").Logger()

	var coupon coupon.Coupon
	if err := f.db.WithContext(ctx).
		Preload("CouponPolicy").
		First(&coupon, "code = ? AND user_id = ?", couponCode, userID).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Msg("coupon not found")
			return nil, exception.NewNotFound("coupon not found", err)
		}
		log.Error().
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Err(err).
			Msg("failed to fetch coupon")
		return nil, exception.NewInternal("failed to fetch coupon", err)
	}

	if err := coupon.Cancel(); err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Err(err).
			Msg("failed to cancel coupon")
		return nil, exception.NewBadRequest("failed to cancel coupon", err)
	}

	if err := f.db.WithContext(ctx).Save(&coupon).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_code", coupon.Code).
			Str("coupon_status", string(coupon.Status)).
			Str("user_id", userID).
			Err(err).
			Msg("failed to save use coupon")
		return nil, exception.NewInternal("failed to save cancel coupon", err)
	}

	span.SetStatus(codes.Ok, "CancelCoupon")
	log.Info().
		Str("coupon_code", coupon.Code).
		Str("coupon_status", string(coupon.Status)).
		Str("user_id", userID).
		Msg("coupon cancel successfully")
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
	log = log.With().Str("func", "FindCouponByCode").Logger()

	var c coupon.Coupon
	if err := f.db.WithContext(ctx).
		Preload("CouponPolicy").
		First(&c, "code = ? AND user_id = ?", couponCode, userID).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Msg("coupon not found")
			return nil, exception.NewNotFound("coupon not found", err)
		}
		log.Error().
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Err(err).
			Msg("failed to fetch coupon")
		return nil, exception.NewInternal("failed to fetch coupon", err)
	}

	span.SetStatus(codes.Ok, "FindCouponByCode")
	log.Info().
		Str("coupon_code", c.Code).
		Str("user_id", userID).
		Msg("coupon found successfully")
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
	log = log.With().Str("func", "FindCouponsByUserID").Logger()

	var coupons []coupon.Coupon
	if err := f.db.WithContext(ctx).
		Preload("CouponPolicy").
		Where("user_id = ?", userID).
		Find(&coupons).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("user_id", userID).
			Err(err).
			Msg("failed to fetch coupons by user")
		return nil, exception.NewInternal("failed to fetch coupons by user", err)
	}

	span.SetStatus(codes.Ok, "FindCouponsByUserID")
	log.Info().
		Str("user_id", userID).
		Int("coupon_count", len(coupons)).
		Msg("fetched coupons for user successfully")
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
	log = log.With().Str("func", "FindCouponsByCouponPolicyCode").Logger()

	var policy coupon.CouponPolicy
	if err := f.db.WithContext(ctx).
		Where("code = ?", couponPolicyCode).
		First(&policy).Error; err != nil {
		span.RecordError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("coupon_policy_code", couponPolicyCode).
				Msg("coupon policy not found")
			return nil, exception.NewNotFound("coupon policy not found", err)
		}
		log.Error().
			Str("coupon_policy_code", couponPolicyCode).
			Err(err).
			Msg("failed to fetch coupon policy")
		return nil, exception.NewInternal("failed to fetch coupon policy", err)
	}

	var coupons []coupon.Coupon
	if err := f.db.WithContext(ctx).
		Where("coupon_policy_id = ?", policy.ID).
		Find(&coupons).Error; err != nil {
		span.RecordError(err)
		log.Error().
			Str("coupon_policy_code", couponPolicyCode).
			Err(err).
			Msg("failed to fetch coupons by policy code")
		return nil, exception.NewInternal("failed to fetch coupons by policy code", err)
	}

	instrument.CouponQuota.WithLabelValues(policy.Code).Set(float64(policy.TotalQuantity))
	instrument.CouponIssued.WithLabelValues(policy.Code).Set(float64(len(coupons)))

	span.SetStatus(codes.Ok, "FindCouponsByCouponPolicyCode")
	log.Info().
		Str("coupon_policy_code", couponPolicyCode).
		Int("coupon_count", len(coupons)).
		Msg("fetched coupons for policy successfully")
	return coupons, nil
}
