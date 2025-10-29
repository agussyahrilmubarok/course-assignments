package server

import (
	"fmt"
	"net/http"
	"os"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	v1 "example.com/coupon/api/v1/handler"
	featureV1 "example.com/coupon/internal/coupon/feature/v1"
	couponMiddleware "example.com/coupon/internal/middleware"

	_ "example.com/coupon/api/v1/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewGinRouter(cfg *config.Config, log zerolog.Logger) *http.Server {
	r := gin.Default()

	db, err := config.NewPostgres(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	if err := db.AutoMigrate(&coupon.Coupon{}, &coupon.CouponPolicy{}); err != nil {
		log.Fatal().Err(err).Msg("AutoMigrate failed")
		os.Exit(1)
	}

	apiRoute := r.Group("/api")

	couponFeatureV1 := featureV1.NewCouponFeature(db)
	couponPolicyHandlerV1 := v1.NewCouponPolicyHandler(db)
	couponHandlerV1 := v1.NewCouponHandler(couponFeatureV1)

	// V1 routes
	v1Group := apiRoute.Group("/v1")
	{
		v1Group.POST("/couponPolicies/dummy", couponPolicyHandlerV1.CreateCouponPolicyDummy)
		v1Group.GET("/couponPolicies", couponPolicyHandlerV1.SearchCouponPolicy)

		v1Group.POST("/coupons/issue", couponMiddleware.UserIDMiddleware(), couponHandlerV1.IssueCoupon)
		v1Group.GET("/coupons/:code", couponMiddleware.UserIDMiddleware(), couponHandlerV1.FindCouponByCode)
		v1Group.POST("/coupons/:code/use", couponMiddleware.UserIDMiddleware(), couponHandlerV1.UseCoupon)
		v1Group.POST("/coupons/:code/cancel", couponMiddleware.UserIDMiddleware(), couponHandlerV1.CancelCoupon)
		v1Group.GET("/coupons/user", couponMiddleware.UserIDMiddleware(), couponHandlerV1.FindCouponsByUserID)
		v1Group.GET("/coupons/policy/:policyCode", couponHandlerV1.FindCouponsByCouponPolicyCode)

		v1Group.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("couponSwaggerV1"),
			ginSwagger.URL("/api/v1/swagger/doc.json"),
		))
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": cfg.App.Name,
		})
	})

	serverAddr := fmt.Sprintf(":%v", cfg.App.Port)

	return &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}
}
