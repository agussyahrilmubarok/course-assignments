package v1

import (
	"context"
	"errors"

	"example.com/coupon-service/internal/coupon"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IService interface {
	IssueCoupon(ctx context.Context, policyCode string, userID string) (*coupon.Coupon, error)
}

type service struct {
	repo   IRepository
	logger *zap.Logger
}

func NewService(
	repo IRepository,
	logger *zap.Logger,
) IService {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) IssueCoupon(ctx context.Context, policyCode string, userID string) (*coupon.Coupon, error) {
	policy, err := s.repo.FindCouponPolicyByCode(ctx, policyCode)
	if err != nil || policy == nil {
		s.logger.Warn("failed to get coupon policy not found", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	if err := policy.IsValidPeriod(); err != nil {
		s.logger.Warn("coupon policy not valid period", zap.String("policy_code", policyCode), zap.Error(err))
		return nil, err
	}

	issued, err := s.repo.CountIssuedCoupons(ctx, policy.ID)
	if err != nil {
		s.logger.Error("failed to count issued coupons", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	if issued >= policy.TotalQuantity {
		s.logger.Warn("coupon quantity exhausted", zap.String("policy_code", policyCode), zap.String("user_id", userID))
		return nil, errors.New("coupon quantity exhausted")
	}

	coupon := &coupon.Coupon{
		ID:             uuid.New().String(),
		Code:           uuid.New().String(),
		Status:         coupon.CouponStatusAvailable,
		UsedAt:         nil,
		UserID:         userID,
		OrderID:        nil,
		CouponPolicyID: policy.ID,
	}

	coupon, err = s.repo.CreateCoupon(ctx, coupon)
	if err != nil {
		s.logger.Error("failed to issue coupon not created", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	s.logger.Info("issue coupon successfully", zap.String("policy_code", policyCode), zap.String("user_id", userID), zap.String("coupon_code", coupon.Code))
	return coupon, nil
}
