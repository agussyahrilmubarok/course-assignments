package coupon

import "time"

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
