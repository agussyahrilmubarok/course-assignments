package instrument

import "github.com/prometheus/client_golang/prometheus"

var (
	CouponIssued = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coupon_issued_total",
			Help: "Number of coupons issued per coupon policy",
		},
		[]string{"coupon_policy_code"},
	)

	CouponQuota = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "coupon_policy_quota",
			Help: "Total quota of coupons per coupon policy",
		},
		[]string{"coupon_policy_code"},
	)
)

func init() {
	prometheus.MustRegister(CouponIssued)
	prometheus.MustRegister(CouponQuota)
}
