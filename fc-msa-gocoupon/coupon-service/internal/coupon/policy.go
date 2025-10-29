package coupon

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountType string

const (
	DiscountTypeFixedAmount DiscountType = "FIXED_AMOUNT"
	DiscountTypePercentage  DiscountType = "PERCENTAGE"
)

type CouponPolicy struct {
	ID                    string       `json:"id" gorm:"primaryKey;column:id"`
	Code                  string       `json:"code" gorm:"unique;size:50;not null"`
	Name                  string       `json:"name" gorm:"not null"`
	Description           string       `json:"description" gorm:"type:text;not null"`
	TotalQuantity         int          `json:"totalQuantity" gorm:"not null;column:total_quantity"`
	StartTime             time.Time    `json:"startTime" gorm:"not null;column:start_time"`
	EndTime               time.Time    `json:"endTime" gorm:"not null;column:end_time"`
	DiscountType          DiscountType `json:"discountType" gorm:"not null;column:discount_type"`
	DiscountValue         int          `json:"discountValue" gorm:"not null;column:discount_value"`
	MinimumOrderAmount    int          `json:"minimumOrderAmount" gorm:"not null;column:minimum_order_amount"`
	MaximumDiscountAmount int          `json:"maximumDiscountAmount" gorm:"not null;column:maximum_discount_amount"`
	CreatedAt             time.Time    `json:"createdAt" gorm:"autoCreateTime;column:created_at"`
	UpdatedAt             time.Time    `json:"updatedAt" gorm:"autoUpdateTime;column:updated_at"`
	Coupons               []Coupon     `json:"coupons,omitempty" gorm:"foreignKey:CouponPolicyID"`
}

func (c *CouponPolicy) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}

var (
	ErrCouponPolicyQoutaExceeded = errors.New("coupon policy quota exceeded")
	ErrCouponPolicyInvalidPeriod = errors.New("coupon policy is not valid in current period")
)

// IsValidPeriod returns true if the current time is within the start and end time of the coupon policy.
func (c *CouponPolicy) IsValidPeriod() bool {
	now := time.Now()
	return !now.Before(c.StartTime) && !now.After(c.EndTime)
}

// GetIssuedQuantity returns the number of issued coupons linked to this coupon policy.
func (c *CouponPolicy) GetIssuedQuantity() int {
	if c.Coupons != nil {
		return len(c.Coupons)
	}
	return 0
}
