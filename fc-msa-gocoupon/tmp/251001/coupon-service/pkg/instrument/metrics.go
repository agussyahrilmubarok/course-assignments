package instrument

import (
	"fmt"
	"net/http"

	"example.com/coupon/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricServer(cfg *config.Config) *http.Server {
	r := gin.New()
	r.Use(gin.Recovery())

	gin.SetMode(gin.ReleaseMode)

	r.GET("/metrics/prometheus", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	serverAddr := fmt.Sprintf(":%v", cfg.App.Metric.Port)

	return &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}
}
