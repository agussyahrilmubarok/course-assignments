package v2

import (
	"context"
	"fmt"

	"example.com/coupon-service/internal/coupon"
	"example.com/coupon-service/internal/instrument/logging"
	"example.com/coupon-service/internal/instrument/metrics"
	"example.com/coupon-service/internal/instrument/tracing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func NewService(
	repo IRepository,
) IService {
	return &service{
		repo: repo,
	}
}

// Problem Summary:
// The IssueCoupon flow historically suffered from race conditions when multiple
// concurrent requests attempted to issue coupons for the same coupon policy.
// The main issue was caused by a read-before-insert logic that wasn't atomic.
//
// Previous Failure Scenario:
// - Multiple concurrent requests executed CountIssuedCoupons() simultaneously.
// - All requests observed issued < TotalQuantity (quota still available).
// - All proceeded to CreateCoupon(), causing overshoot of quota.
// - Some transactions failed later when COUNT(*) exceeded TotalQuantity.
// - Result: inconsistent quota, random success_count, and user-level unfairness.
//
// Root Causes:
// - The quota check and coupon creation were not executed atomically.
// - Missing transactional locking (no SELECT ... FOR UPDATE).
// - Database was not protecting quota at the row level.
//
// Fix Implemented:
//   - The code now wraps the entire quota check and coupon creation in a single
//     SERIALIZABLE / REPEATABLE READ transaction.
//   - Coupon policy row is locked via SELECT ... FOR UPDATE.
//   - Ensures that only one request can pass quota validation at a time.
//   - Prevents race conditions and quota overshoot.
//
// Remaining Risks / What Could Still Go Wrong:
// NOTE: Even with row-level locking fixing race conditions, several important
//
//	business and infrastructure issues remain possible:
//
// 1. Performance Bottlenecks:
//   - High lock contention on coupon_policies row during traffic spikes.
//   - Requests may wait 5–20 seconds on locks under heavy concurrency.
//   - Potential transaction timeout or cancellation.
//
// 2. Database Hotspot:
//   - One coupon policy = one locked row.
//   - All clients must serialize through the same row, creating a hotspot.
//   - Throughput becomes DB-bound, not service-bound.
//
// 3. CountIssuedCoupons() Overhead:
//   - COUNT(*) per request becomes expensive for large coupon tables.
//   - May cause slow sequential scans or heavy index scans over time.
//   - May increase lock wait time inside transaction.
//
// 4. Potential Deadlocks:
//   - If other transactions lock rows in a different order.
//   - Even with SELECT ... FOR UPDATE, ordering mismatch can deadlock.
//   - Must ensure consistent lock ordering everywhere in the codebase.
//
// 5. User Eligibility / Abuse:
//   - No per-user coupon limit enforced (user may issue multiple coupons).
//   - Missing unique constraint on (user_id, policy_id) for one-per-user rule.
//   - No idempotency key (users can double-click and issue two coupons).
//   - No rate-limiting → user can spam requests.
//
// 6. Business Rule Edge Cases:
//   - Validity period check may pass even if request enters near expiration
//     but finishes after expiration.
//   - No checks for product/order requirements in policy.
//
// 7. Coupon Integrity Risks:
//   - UUID duplication extremely unlikely but still requires UNIQUE index.
//   - Insert failure after quota decrement (if using issued_count later).
//
// 8. Logging & Observability Impact:
//   - Logging inside the transaction increases transaction duration.
//   - High volume logs under load amplify lock contention.
//
// 9. Infrastructure Risks:
//   - Connection pool exhaustion if many requests block on locked row.
//   - Slow queries cause goroutine buildup if context timeouts not aligned.
//   - Retry/backoff not implemented for transient DB errors.
//
// 10. Scalability Limitations:
//   - This approach works but does not scale to high-traffic flash-sale
//     scenarios (tens of thousands of QPS).
//   - Requires additional techniques like Redis atomic counters,
//     queue-based issuance, or pre-generated coupon pools.
//
// Summary:
// - Race condition is fixed.
// - But the system still has significant bottlenecks and missing business rules.
//
// Recommendation:
// - Implement per-user constraints, unique indexes, and idempotency keys.
// - Replace COUNT(*) with atomic increment counters.
// - Consider Redis or queue-based issuance for high scale.
// - Add retry/backoff and refine timeouts to avoid lock starvation.
// - Improve transaction efficiency to prevent DB hotspot issues.
func (s *service) IssueCoupon(ctx context.Context, policyCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := tracing.StartSpan(ctx, "V2.Service.IssueCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	couponIssueDuration := prometheus.NewTimer(
		metrics.CouponIssueDuration.WithLabelValues(policyCode, "v2"),
	)
	defer couponIssueDuration.ObserveDuration()

	var createdCoupon *coupon.Coupon

	err := s.repo.WithTx(ctx, func(tx pgx.Tx) error {
		// Retrieve Coupon Policy
		policy, err := s.repo.FindCouponPolicyByCodeForUpdateTx(ctx, tx, policyCode)
		if err != nil || policy == nil {
			span.RecordError(err)
			log.Warn("failed to get coupon policy not found", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
			return coupon.ErrCouponPolicyNotFound
		}

		// Check Valid Period
		if err := policy.IsValidPeriod(); err != nil {
			span.RecordError(err)
			log.Warn("coupon policy not valid period", zap.String("policy_code", policyCode), zap.Error(err))
			return err
		}

		// Check Available Quantity
		issued, err := s.repo.CountIssuedCouponsTx(ctx, tx, policy.ID)
		if err != nil {
			span.RecordError(err)
			log.Error("failed to count issued coupons", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
			return coupon.ErrCouponInternal
		}

		if issued >= policy.TotalQuantity {
			err := fmt.Errorf("%w, %v qoutas", coupon.ErrCouponPolicyQuantityExceed, policy.TotalQuantity)
			span.RecordError(err)
			log.Warn("coupon quantity exhausted", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
			return err
		}

		// TODO: Check User Eligibility
		// TODO: Check Order / Product Requirements (optional)

		// Create New Coupon
		tempCoupon := &coupon.Coupon{
			ID:             uuid.New().String(),
			Code:           uuid.New().String(),
			Status:         coupon.CouponStatusAvailable,
			UsedAt:         nil,
			UserID:         userID,
			OrderID:        nil,
			CouponPolicyID: policy.ID,
		}

		tempCoupon, err = s.repo.CreateCouponTx(ctx, tx, tempCoupon)
		if err != nil {
			span.RecordError(err)
			log.Error("failed to issue coupon not created", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
			return coupon.ErrCouponInternal
		}

		createdCoupon = tempCoupon
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return New Coupon
	log.Info("issue coupon successfully", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.String("coupon_code", createdCoupon.Code))
	return createdCoupon, nil
}

func (s *service) UseCoupon(ctx context.Context, couponCode string, userID string, orderID string) (*coupon.Coupon, error) {
	ctx, span := tracing.StartSpan(ctx, "V2.Service.UseCoupon")
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
	ctx, span := tracing.StartSpan(ctx, "V2.Service.CancelCoupon")
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
	ctx, span := tracing.StartSpan(ctx, "V2.Service.FindCouponByCode")
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
