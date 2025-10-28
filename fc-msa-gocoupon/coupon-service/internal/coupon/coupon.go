package coupon

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CouponStatus string

const (
	CouponStatusAvailable CouponStatus = "AVAILABLE"
	CouponStatusUsed      CouponStatus = "USED"
	CouponStatusExpired   CouponStatus = "EXPIRED"
	CouponStatusCanceled  CouponStatus = "CANCELED"
)

type Coupon struct {
	ID             string        `json:"id" gorm:"primaryKey;column:id"`
	Code           string        `json:"code" gorm:"unique;size:50;not null"`
	Status         CouponStatus  `json:"status" gorm:"not null"`
	UsedAt         *time.Time    `json:"usedAt,omitempty" gorm:"column:used_at"`
	UserID         string        `json:"userId" gorm:"not null;column:user_id"`
	OrderID        *string       `json:"orderId,omitempty" gorm:"column:order_id"`
	CouponPolicyID string        `json:"couponPolicyId" gorm:"not null;column:coupon_policy_id"`
	CreatedAt      time.Time     `json:"createdAt" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt      time.Time     `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at"`
	CouponPolicy   *CouponPolicy `json:"couponPolicy,omitempty" gorm:"foreignKey:CouponPolicyID"`
}

var (
	ErrCouponAlreadyUsed = errors.New("coupon has already been used")
	ErrCouponExpired     = errors.New("coupon has expired")
	ErrCouponNotUsed     = errors.New("coupon has not been used")
)

func (c *Coupon) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}

// IsExpired returns true if current time is before start or after end
func (c *Coupon) IsExpired() bool {
	now := time.Now()
	return now.Before(c.CouponPolicy.StartTime) || now.After(c.CouponPolicy.EndTime)
}

// IsUsed returns true if coupon status is USED
func (c *Coupon) IsUsed() bool {
	return c.Status == CouponStatusUsed
}

// Use marks the coupon as used with given orderId, or returns an error
func (c *Coupon) Use(orderId string) error {
	if c.IsUsed() {
		return ErrCouponAlreadyUsed
	}
	if c.IsExpired() {
		return ErrCouponExpired
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
