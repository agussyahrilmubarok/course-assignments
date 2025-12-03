package metrics

import (
	"example.com/coupon-service/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	CouponIssueDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "coupon_issue_duration_seconds",
			Help: "Duration to issue coupon",
			Buckets: []float64{
				0.005, 0.01, 0.025, 0.05, 0.1,
				0.2, 0.5, 1, 2, 5,
			},
		},
		[]string{"policy_code", "version"},
	)
)

func init() {
	prometheus.MustRegister(CouponIssueDuration)
}

func NewMetricServer(cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.GET("/metrics/prometheus", func(c echo.Context) error {
		promhttp.Handler().ServeHTTP(c.Response(), c.Request())
		return nil
	})
	return e
}
