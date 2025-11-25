package v1

import (
	"example.com/coupon-service/internal/api/middleware"
	"example.com/coupon-service/internal/config"
	"github.com/labstack/echo/v4"

	_ "example.com/coupon-service/internal/api/v1/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Coupon API V1
// @version 1.0
// @description Coupon API V1
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-USER-ID
func RegisterAPIV1(group *echo.Group, pg *config.Postgres) {
	repository := NewRepository(pg)
	service := NewService(repository)
	handler := NewHandler(service)

	coupons := group.Group("/v1/coupons")
	coupons.POST("/issue", handler.IssueCoupon, middleware.UserIDMiddleware())
	coupons.POST("/use", handler.UseCoupon, middleware.UserIDMiddleware())
	coupons.POST("/cancel", handler.CancelCoupon, middleware.UserIDMiddleware())
	coupons.GET("/:coupon_code", handler.FindCouponByCode, middleware.UserIDMiddleware())

	coupons.GET("/swagger/*", echoSwagger.EchoWrapHandler(
		echoSwagger.InstanceName("couponsApiV1"),
		echoSwagger.URL("/api/swagger/doc.json"),
	))
}
