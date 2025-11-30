package metrics

import (
	"example.com/coupon-service/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricServer(cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.GET("/metrics/prometheus", func(c echo.Context) error {
		promhttp.Handler().ServeHTTP(c.Response(), c.Request())
		return nil
	})
	return e
}
