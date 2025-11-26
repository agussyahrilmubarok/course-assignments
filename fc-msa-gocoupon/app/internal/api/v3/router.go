package v3

import (
	"example.com/coupon-service/internal/api/middleware"
	"example.com/coupon-service/internal/config"
	"github.com/labstack/echo/v4"

	_ "example.com/coupon-service/internal/api/v3/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Coupon API V3
// @version 3.0
// @description Coupon API V3
// @BasePath /api/V3

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-USER-ID
func RegisterAPIV3(group *echo.Group, pg *config.Postgres, rdb *config.Redis) {
	repository := NewRepository(pg, rdb)
	service := NewService(repository)
	handler := NewHandler(service)

	coupons := group.Group("/v3/coupons")
	coupons.POST("/issue", handler.IssueCoupon, middleware.UserIDMiddleware())
	coupons.POST("/use", handler.UseCoupon, middleware.UserIDMiddleware())
	coupons.POST("/cancel", handler.CancelCoupon, middleware.UserIDMiddleware())
	coupons.GET("/:coupon_code", handler.FindCouponByCode, middleware.UserIDMiddleware())

	coupons.GET("/swagger/*", echoSwagger.EchoWrapHandler(
		echoSwagger.InstanceName("couponsApiV3"),
		echoSwagger.URL("/api/swagger/doc.json"),
	))
}
