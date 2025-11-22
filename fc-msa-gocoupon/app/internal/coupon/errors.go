package coupon

import "errors"

var (
	ErrCouponPolicyNotActive       = errors.New("coupon policy not active yet")
	ErrCouponPolicyExpired         = errors.New("coupon policy expired")
	ErrCouponPolicyQuantityExceed  = errors.New("coupon quantity exhausted")
	ErrCouponAlreadyUsed           = errors.New("coupon has already been used")
	ErrCouponCanceled              = errors.New("coupon canceled")
	ErrCouponExpired               = errors.New("coupon has expired")
	ErrCouponNotUsed               = errors.New("coupon has not been used")
	ErrCouponNotOwner              = errors.New("not the owner of this coupon")
	ErrCouponTooManyRequests       = errors.New("too many concurrent coupon requests")
	ErrCouponInvalidForOrder       = errors.New("coupon not valid for this order")
	ErrCouponUserLimitExceeded     = errors.New("user has already claimed this coupon")
	ErrCouponOrderAmountTooLow     = errors.New("order amount is below coupon minimum requirement")
	ErrCouponInvalidForProduct     = errors.New("coupon not applicable for selected product")
	ErrCouponQuantityRaceCondition = errors.New("coupon quantity limit reached (race condition)")
	ErrCouponUserAlreadyClaimed    = errors.New("user has already claimed this coupon")
)

var (
	ErrCouponInternal       = errors.New("internal coupon service error")
	ErrCouponNotFound       = errors.New("coupon not found")
	ErrCouponCounted        = errors.New("failed to count issued coupons")
	ErrCouponCreated        = errors.New("failed to create coupon")
	ErrCouponPolicyNotFound = errors.New("coupon policy not found")
	ErrDatabaseUnavailable  = errors.New("database unavailable")
	ErrTransactionFailed    = errors.New("transaction failed")
	ErrTimeout              = errors.New("timeout during database operation")
	ErrUnknown              = errors.New("unknown technical error")
)
