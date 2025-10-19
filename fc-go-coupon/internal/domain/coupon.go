package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Coupon struct {
	ID             string       `gorm:"primaryKey;not null;unique;column:id" json:"id"`
	Code           string       `gorm:"not null;unique;size:50" json:"code"`
	Status         CouponStatus `gorm:"not null;column:status" json:"status"`
	UsedAt         *time.Time   `gorm:"column:used_at" json:"usedAt,omitempty"`
	UserID         string       `gorm:"not null;column:user_id" json:"userId"`
	OrderID        *string      `gorm:"column:order_id" json:"orderId,omitempty"`
	CouponPolicyID string       `gorm:"not null;column:coupon_policy_id" json:"couponPolicyId"`
	CreatedAt      time.Time    `gorm:"autoCreateTime;column:created_at" json:"createdAt"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`
	CouponPolicy   CouponPolicy `gorm:"foreignKey:CouponPolicyID" json:"couponPolicy"`
}

type CouponStatus string

const (
	CouponStatusAvailable CouponStatus = "AVAILABLE"
	CouponStatusUsed      CouponStatus = "USED"
	CouponStatusExpired   CouponStatus = "EXPIRED"
	CouponStatusCanceled  CouponStatus = "CANCELED"
)

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
