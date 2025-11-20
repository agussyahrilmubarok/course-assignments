package coupon

import "time"

type CouponStatus string

const (
	CouponStatusAvailable CouponStatus = "AVAILABLE"
	CouponStatusUsed      CouponStatus = "USED"
	CouponStatusExpired   CouponStatus = "EXPIRED"
	CouponStatusCanceled  CouponStatus = "CANCELED"
)

type Coupon struct {
	ID             string       `json:"id"`
	Code           string       `json:"code"`
	Status         CouponStatus `json:"status"`
	UsedAt         *time.Time   `json:"used_at,omitempty"`
	UserID         string       `json:"user_id"`
	OrderID        *string      `json:"order_id,omitempty"`
	CouponPolicyID string       `json:"coupon_policy_id"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`

	CouponPolicy *CouponPolicy `json:"coupon_policy,omitempty"`
}
