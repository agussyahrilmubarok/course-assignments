package coupon

type IssueCouponRequest struct {
	PolicyCode string `json:"policy_code"`
}

type UseCouponRequest struct {
	CouponCode string `json:"coupon_code"`
	OrderID    string `json:"order_id"`
}

type CancelCouponRequest struct {
	CouponCode string `json:"coupon_code"`
}
