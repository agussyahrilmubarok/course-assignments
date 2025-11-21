package v1

import (
	"example.com/coupon-service/internal/api/middleware"
	"example.com/coupon-service/internal/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterAPIV1(group *echo.Group, pg *config.Postgres, logger *zap.Logger) {
	repository := NewRepository(pg, logger)
	service := NewService(repository, logger)
	handler := NewHandler(service, logger)

	coupons := group.Group("/coupons")
	coupons.POST("/issue", handler.IssueCoupon, middleware.UserIDMiddleware())
}
