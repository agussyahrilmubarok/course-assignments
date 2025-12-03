package v3

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
// Previous Failure Scenario:
// Root Causes:
// Fix Implemented:
// Potential Issues / What could go wrong:
func (s *service) IssueCoupon(ctx context.Context, policyCode string, userID string) (*coupon.Coupon, error) {
	ctx, span := tracing.StartSpan(ctx, "V3.Service.IssueCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	couponIssueDuration := prometheus.NewTimer(
		metrics.CouponIssueDuration.WithLabelValues(policyCode, "v3"),
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
		available, err := s.repo.GetCouponPolicyQuantity(ctx, policy.Code)
		if err != nil {
			issued, err := s.repo.CountIssuedCouponsTx(ctx, tx, policy.ID)
			if err != nil {
				span.RecordError(err)
				log.Error("failed to count issued coupons", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
				return coupon.ErrCouponInternal
			}

			available = policy.TotalQuantity - issued
			_ = s.repo.SetCouponPolicyQuantity(ctx, policy.Code, available, policy.EndTime)
		}

		if available <= 0 {
			err := fmt.Errorf("%w, %v quotas", coupon.ErrCouponPolicyQuantityExceed, policy.TotalQuantity)
			span.RecordError(err)
			log.Warn("coupon quantity exhausted (redis)", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
			return err
		}

		err = s.repo.DecrCouponPolicyQuantity(ctx, policy.Code)
		if err != nil {
			span.RecordError(err)
			log.Error("failed to decrement redis quota", zap.String("policy_code", policyCode), zap.Error(err))
			return coupon.ErrCouponInternal
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
			_ = s.repo.IncrCouponPolicyQuantity(ctx, policy.Code)
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
	ctx, span := tracing.StartSpan(ctx, "V3.Service.UseCoupon")
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
	ctx, span := tracing.StartSpan(ctx, "V3.Service.CancelCoupon")
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
	ctx, span := tracing.StartSpan(ctx, "V3.Service.FindCouponByCode")
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
