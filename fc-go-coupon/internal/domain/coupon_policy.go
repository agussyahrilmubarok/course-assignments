package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CouponPolicy struct {
	ID                    string       `gorm:"primaryKey;not null;unique;column:id" json:"id"`
	Name                  string       `gorm:"not null" json:"name"`
	Description           string       `gorm:"not null;type:text" json:"description"`
	TotalQuantity         int          `gorm:"not null;column:total_quantity" json:"totalQuantity"`
	StartTime             time.Time    `gorm:"not null;column:start_time" json:"startTime"`
	EndTime               time.Time    `gorm:"not null;column:end_time" json:"endTime"`
	DiscountType          DiscountType `gorm:"not null;column:discount_type" json:"discountType"`
	DiscountValue         int          `gorm:"not null;column:discount_value" json:"discountValue"`
	MinimumOrderAmount    int          `gorm:"not null;column:minimum_order_amount" json:"minimumOrderAmount"`
	MaximumDiscountAmount int          `gorm:"not null;column:maximum_discount_amount" json:"maximumDiscountAmount"`
	CreatedAt             time.Time    `gorm:"autoCreateTime;column:created_at" json:"createdAt"`
	UpdatedAt             time.Time    `gorm:"autoUpdateTime;column:updated_at" json:"updatedAt"`
	Coupons               []Coupon     `gorm:"foreignKey:CouponPolicyID" json:"coupons"`
}

type DiscountType string

const (
	DiscountTypeFixedAmount DiscountType = "FIXED_AMOUNT"
	DiscountTypePercentage  DiscountType = "PERCENTAGE"
)

func (c *CouponPolicy) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}

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
