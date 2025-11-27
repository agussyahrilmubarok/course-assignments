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

type IssueCouponMessage struct {
	PolicyID   string `json:"policy_id"`
	PolicyCode string `json:"policy_code"`
	CouponID   string `json:"coupon_id"`
	CouponCode string `json:"coupon_code"`
	UserID     string `json:"user_id"`
}
