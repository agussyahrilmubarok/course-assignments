package coupon

import "time"

type CouponStatus string

const (
	CouponStatusPending   CouponStatus = "PENDING"
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

// Use marks the coupon as used with given orderId, or returns an error
func (c *Coupon) Use(orderId string) error {
	if c.Status == CouponStatusUsed {
		return ErrCouponAlreadyUsed
	}
	if c.Status == CouponStatusExpired {
		return ErrCouponExpired
	}
	if c.Status == CouponStatusCanceled {
		return ErrCouponCanceled
	}
	if c.Status == CouponStatusPending {
		return ErrCouponPending
	}

	now := time.Now()
	c.Status = CouponStatusUsed
	c.OrderID = &orderId
	c.UsedAt = &now
	return nil
}

// Cancel reverts the coupon to CANCELED if previously used, or returns an error
func (c *Coupon) Cancel() error {
	if c.Status != CouponStatusUsed {
		return ErrCouponNotUsed
	}

	c.Status = CouponStatusCanceled
	c.OrderID = nil
	c.UsedAt = nil
	return nil
}
