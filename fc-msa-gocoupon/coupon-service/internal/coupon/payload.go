package coupon

type IssueCouponRequest struct {
	CouponPolicyCode string `json:"couponPolicyCode"`
}

type IssueCouponMessage struct {
	CouponPolicyCode string `json:"couponPolicyCode"`
	UserID           string `json:"userId"`
}

type UseCouponRequest struct {
	OrderID string `json:"order_id"`
}
