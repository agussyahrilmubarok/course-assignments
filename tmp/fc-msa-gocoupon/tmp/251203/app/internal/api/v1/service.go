package v1

import (
	"context"
	"fmt"

	"example.com/coupon-service/internal/coupon"
	"example.com/coupon-service/internal/instrument/logging"
	"example.com/coupon-service/internal/instrument/metrics"
	"example.com/coupon-service/internal/instrument/tracing"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type IService interface {
	IssueCoupon(ctx context.Context, policyCode string, userID string) (*coupon.Coupon, error)
	UseCoupon(ctx context.Context, couponCode string, userID string, orderID string) (*coupon.Coupon, error)
	CancelCoupon(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error)
	FindCouponByCode(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error)
}

type service struct {
	repo IRepository
}

func NewService(repo IRepository) IService {
	return &service{
		repo: repo,
	}
}

// Problem Current Code:
// When handler receives concurrent requests (e.g., 100 requests at the same time)
// Got race condition:
//   - Multiple requests call CountIssuedCoupons() at the same time
//   - All see issued < TotalQuantity (quota) and proceed to create coupon
//   - Result: more coupons are created than allowed quota
//   - Some requests eventually fail because COUNT in DB exceeds TotalQuantity
//
// Cause:
//   - The check (issued >= TotalQuantity) and insert (CreateCoupon) are NOT atomic
//   - No row locking or transaction isolation is applied
//
// Consequence:
//   - With high concurrency, success_count may be random (less than quota)
//   - Violates coupon policy quota
//
// Solutions:
//  1. Use database transaction with SELECT ... FOR UPDATE to lock row
//  2. Use DB constraint or stored procedure to enforce quota atomically
//  3. Implement optimistic locking with a version field
//
// Potential Issues / What could go wrong:
//   - Duplicate coupon codes (even with UUID, need unique constraint in DB)
//   - Quota enforcement inaccurate due to read-before-insert (more coupons than allowed)
//   - Transaction failures / partial writes leading to inconsistent state
//   - Database timeout or high latency under load
//   - User eligibility check not implemented (users might get more coupons than allowed)
//   - Policy validity edge cases (coupon issued outside valid period if request timing is tight)
//   - Logging overhead slowing down request handling under high concurrency
//   - Multiple requests for same user could bypass intended per-user limits
//   - No retry or backoff on DB conflicts or transient errors
//   - Potential deadlocks if DB row-level locking implemented incorrectly
func (s *service) IssueCoupon(ctx context.Context, policyCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := tracing.StartSpan(ctx, "V1.Service.IssueCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	couponIssueDuration := prometheus.NewTimer(
		metrics.CouponIssueDuration.WithLabelValues(policyCode, "v1"),
	)
	defer couponIssueDuration.ObserveDuration()

	// Retrieve Coupon Policy
	policy, err := s.repo.FindCouponPolicyByCode(ctx, policyCode)
	if err != nil || policy == nil {
		span.RecordError(err)
		log.Warn("failed to get coupon policy not found", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
		return nil, coupon.ErrCouponPolicyNotFound
	}

	// Check Valid Period
	if err := policy.IsValidPeriod(); err != nil {
		span.RecordError(err)
		log.Warn("coupon policy not valid period", zap.String("policy_code", policyCode), zap.Error(err))
		return nil, err
	}

	// Check Available Quantity
	issued, err := s.repo.CountIssuedCoupons(ctx, policy.ID)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to count issued coupons", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
		return nil, coupon.ErrCouponInternal
	}

	if issued >= policy.TotalQuantity {
		err := fmt.Errorf("%w, %v qoutas", coupon.ErrCouponPolicyQuantityExceed, policy.TotalQuantity)
		span.RecordError(err)
		log.Warn("coupon quantity exhausted", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	// TODO: Check User Eligibility
	// TODO: Check Order / Product Requirements (optional)

	// Create New Coupon
	newCoupon := &coupon.Coupon{
		ID:             uuid.New().String(),
		Code:           uuid.New().String(),
		Status:         coupon.CouponStatusAvailable,
		UsedAt:         nil,
		UserID:         userID,
		OrderID:        nil,
		CouponPolicyID: policy.ID,
	}
	newCoupon, err = s.repo.CreateCoupon(ctx, newCoupon)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to issue coupon not created", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
		return nil, coupon.ErrCouponInternal
	}

	// Return New Coupon
	log.Info("issue coupon successfully", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.String("coupon_code", newCoupon.Code))
	return newCoupon, nil
}

func (s *service) UseCoupon(ctx context.Context, couponCode string, userID string, orderID string) (*coupon.Coupon, error) {
	ctx, span := tracing.StartSpan(ctx, "V1.Service.UseCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	// Retrieve Coupon Policy
	c, err := s.repo.FindCouponByCode(ctx, couponCode)
	if err != nil || c == nil {
		span.RecordError(err)
		log.Warn("failed to get coupon by code", zap.String("coupon_code", couponCode), zap.Error(err))
		return nil, coupon.ErrCouponNotFound
	}

	// Check Coupon Owner
	if c.UserID != userID {
		err := coupon.ErrCouponNotOwner
		span.RecordError(err)
		log.Warn("failed to use coupon not owner", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	// TODO: Retrieve Coupon Policy
	// TODO: Check Policy Validity

	// Check Coupon Status
	if err := c.Use(orderID); err != nil {
		span.RecordError(err)
		log.Warn("failed to use coupon not match status", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.String("order_id", orderID), zap.Error(err))
		return nil, err
	}

	// Update Coupon
	updatedCoupon, err := s.repo.UpdateCoupon(ctx, c)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to update coupon", zap.String("coupon_code", couponCode), zap.Error(err))
		return nil, coupon.ErrCouponInternal
	}

	log.Info("coupon used successfully", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.String("order_id", orderID))
	return updatedCoupon, nil
}

func (s *service) CancelCoupon(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := tracing.StartSpan(ctx, "V1.Service.CancelCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	// Retrieve Coupon Policy
	c, err := s.repo.FindCouponByCode(ctx, couponCode)
	if err != nil || c == nil {
		span.RecordError(err)
		log.Warn("failed to get coupon by code", zap.String("coupon_code", couponCode), zap.Error(err))
		return nil, coupon.ErrCouponNotFound
	}

	// Check Coupon Owner
	if c.UserID != userID {
		err := coupon.ErrCouponNotOwner
		span.RecordError(err)
		log.Warn("failed to cancel coupon not owner", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Check Coupon Status
	if err := c.Cancel(); err != nil {
		span.RecordError(err)
		log.Warn("failed to cancel coupon not match status", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Update Coupon
	updatedCoupon, err := s.repo.UpdateCoupon(ctx, c)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to update coupon", zap.String("coupon_code", couponCode), zap.Error(err))
		return nil, coupon.ErrCouponInternal
	}

	log.Info("coupon cancel successfully", zap.String("coupon_code", couponCode), zap.String("user_id", userID))
	return updatedCoupon, nil
}

func (s *service) FindCouponByCode(ctx context.Context, couponCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := tracing.StartSpan(ctx, "V1.Service.FindCouponByCode")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	// Retrieve Coupon Policy
	c, err := s.repo.FindCouponByCode(ctx, couponCode)
	if err != nil || c == nil {
		span.RecordError(err)
		log.Warn("failed to get coupon by code", zap.String("coupon_code", couponCode), zap.Error(err))
		return nil, coupon.ErrCouponNotFound
	}

	// Check Coupon Owner
	if c.UserID != userID {
		err := coupon.ErrCouponNotOwner
		span.RecordError(err)
		log.Warn("failed to get coupon not owner", zap.String("coupon_code", couponCode), zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	// Retrieve Coupon Policy
	policy, err := s.repo.FindCouponPolicyByID(ctx, c.CouponPolicyID)
	if err != nil || policy == nil {
		span.RecordError(err)
		log.Warn("failed to get coupon policy", zap.String("coupon_id", c.ID), zap.String("coupon_code", couponCode), zap.Error(err))
		return nil, coupon.ErrCouponPolicyNotFound
	}
	c.CouponPolicy = policy

	// Return coupon with policy
	log.Info("returning coupon with attached policy", zap.String("coupon_id", c.ID), zap.String("coupon_code", couponCode), zap.String("user_id", userID))
	return c, nil
}
