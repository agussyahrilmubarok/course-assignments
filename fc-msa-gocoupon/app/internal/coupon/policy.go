package coupon

import (
	"fmt"
	"time"
)

type DiscountType string

const (
	DiscountTypeFixedAmount DiscountType = "FIXED_AMOUNT"
	DiscountTypePercentage  DiscountType = "PERCENTAGE"
)

type CouponPolicy struct {
	ID                    string       `json:"id"`
	Code                  string       `json:"code"`
	Name                  string       `json:"name"`
	Description           string       `json:"description"`
	TotalQuantity         int          `json:"total_quantity"`
	StartTime             time.Time    `json:"start_time"`
	EndTime               time.Time    `json:"end_time"`
	DiscountType          DiscountType `json:"discount_type"`
	DiscountValue         int          `json:"discount_value"`
	MinimumOrderAmount    int          `json:"minimum_order_amount"`
	MaximumDiscountAmount int          `json:"maximum_discount_amount"`
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`

	Coupons []Coupon `json:"coupons,omitempty"`
}

func (c *CouponPolicy) IsValidPeriod() error {
	now := time.Now().UTC()
	start := c.StartTime.UTC()
	end := c.EndTime.UTC()

	if now.Before(start) {
		return fmt.Errorf("coupon policy not active yet, starts at %s", start)
	}
	
	if now.After(end) {
		return fmt.Errorf("coupon policy expired at %s", end)
	}

	return nil
}
