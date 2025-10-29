package v1

import (
	"context"

	"example.com/coupon/internal/coupon"
)

//go:generate mockery --name=ICouponService
type ICouponService interface {
	IssueCoupon(ctx context.Context) (*coupon.Coupon, error)
	UseCoupon(ctx context.Context) (*coupon.Coupon, error)
	CancelCoupon(ctx context.Context) (*coupon.Coupon, error)
	FindCoupon(ctx context.Context) (*coupon.Coupon, error)
	FindCouponsByUserID(ctx context.Context) ([]coupon.Coupon, error)
	FindCouponsByCouponPolicyID(ctx context.Context) ([]coupon.Coupon, error)
}
