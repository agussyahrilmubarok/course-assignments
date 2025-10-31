package server

import (
	"fmt"
	"net/http"
	"os"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/config"
	"example.com/coupon/pkg/instrument"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"

	v1 "example.com/coupon/api/v1/handler"
	v2 "example.com/coupon/api/v2/handler"
	v3 "example.com/coupon/api/v3/handler"
	v4 "example.com/coupon/api/v4/handler"
	featureV1 "example.com/coupon/internal/coupon/feature/v1"
	featureV2 "example.com/coupon/internal/coupon/feature/v2"
	featureV3 "example.com/coupon/internal/coupon/feature/v3"
	featureV4 "example.com/coupon/internal/coupon/feature/v4"
	couponMiddleware "example.com/coupon/internal/middleware"

	_ "example.com/coupon/api/v1/docs"
	_ "example.com/coupon/api/v2/docs"
	_ "example.com/coupon/api/v3/docs"
	_ "example.com/coupon/api/v4/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewGinRouter(cfg *config.Config, log zerolog.Logger) *http.Server {
	db, err := config.NewPostgres(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	if err := db.AutoMigrate(&coupon.Coupon{}, &coupon.CouponPolicy{}); err != nil {
		log.Fatal().Err(err).Msg("AutoMigrate failed")
		os.Exit(1)
	}

	tracer := otel.Tracer(cfg.App.Name)

	r := gin.Default()
	r.Use(instrument.Middleware(tracer, log))
	r.Use(instrument.MetricAppMiddleware)

	apiRoute := r.Group("/api")

	couponFeatureV1 := featureV1.NewCouponFeature(db, log, tracer)
	couponPolicyHandlerV1 := v1.NewCouponPolicyHandler(db, log, tracer)
	couponHandlerV1 := v1.NewCouponHandler(couponFeatureV1, log, tracer)

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

	couponFeatureV2 := featureV2.NewCouponFeature(db, log, tracer)
	couponPolicyHandlerV2 := v2.NewCouponPolicyHandler(db, log, tracer)
	couponHandlerV2 := v2.NewCouponHandler(couponFeatureV2, log, tracer)

	// V2 routes
	v2Group := apiRoute.Group("/v2")
	{
		v2Group.POST("/couponPolicies/dummy", couponPolicyHandlerV2.CreateCouponPolicyDummy)
		v2Group.GET("/couponPolicies", couponPolicyHandlerV2.SearchCouponPolicy)

		v2Group.POST("/coupons/issue", couponMiddleware.UserIDMiddleware(), couponHandlerV2.IssueCoupon)
		v2Group.GET("/coupons/:code", couponMiddleware.UserIDMiddleware(), couponHandlerV2.FindCouponByCode)
		v2Group.POST("/coupons/:code/use", couponMiddleware.UserIDMiddleware(), couponHandlerV2.UseCoupon)
		v2Group.POST("/coupons/:code/cancel", couponMiddleware.UserIDMiddleware(), couponHandlerV2.CancelCoupon)
		v2Group.GET("/coupons/user", couponMiddleware.UserIDMiddleware(), couponHandlerV2.FindCouponsByUserID)
		v2Group.GET("/coupons/policy/:policyCode", couponHandlerV2.FindCouponsByCouponPolicyCode)

		v2Group.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("couponSwaggerV2"),
			ginSwagger.URL("/api/v2/swagger/doc.json"),
		))
	}

	couponFeatureV3 := featureV3.NewCouponFeature(db, log, tracer)
	couponPolicyHandlerV3 := v3.NewCouponPolicyHandler(db, log, tracer)
	couponHandlerV3 := v3.NewCouponHandler(couponFeatureV3, log, tracer)

	// V3 routes
	v3Group := apiRoute.Group("/v3")
	{
		v3Group.POST("/couponPolicies/dummy", couponPolicyHandlerV3.CreateCouponPolicyDummy)
		v3Group.GET("/couponPolicies", couponPolicyHandlerV3.SearchCouponPolicy)

		v3Group.POST("/coupons/issue", couponMiddleware.UserIDMiddleware(), couponHandlerV3.IssueCoupon)
		v3Group.GET("/coupons/:code", couponMiddleware.UserIDMiddleware(), couponHandlerV3.FindCouponByCode)
		v3Group.POST("/coupons/:code/use", couponMiddleware.UserIDMiddleware(), couponHandlerV3.UseCoupon)
		v3Group.POST("/coupons/:code/cancel", couponMiddleware.UserIDMiddleware(), couponHandlerV3.CancelCoupon)
		v3Group.GET("/coupons/user", couponMiddleware.UserIDMiddleware(), couponHandlerV3.FindCouponsByUserID)
		v3Group.GET("/coupons/policy/:policyCode", couponHandlerV3.FindCouponsByCouponPolicyCode)

		v3Group.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("couponSwaggerV3"),
			ginSwagger.URL("/api/v3/swagger/doc.json"),
		))
	}

	couponFeatureV4 := featureV4.NewCouponFeature(db, log, tracer)
	couponPolicyHandlerV4 := v4.NewCouponPolicyHandler(db, log, tracer)
	couponHandlerV4 := v4.NewCouponHandler(couponFeatureV4, log, tracer)

	// V4 routes
	v4Group := apiRoute.Group("/v4")
	{
		v4Group.POST("/couponPolicies/dummy", couponPolicyHandlerV4.CreateCouponPolicyDummy)
		v4Group.GET("/couponPolicies", couponPolicyHandlerV4.SearchCouponPolicy)

		v4Group.POST("/coupons/issue", couponMiddleware.UserIDMiddleware(), couponHandlerV4.IssueCoupon)
		v4Group.GET("/coupons/:code", couponMiddleware.UserIDMiddleware(), couponHandlerV4.FindCouponByCode)
		v4Group.POST("/coupons/:code/use", couponMiddleware.UserIDMiddleware(), couponHandlerV4.UseCoupon)
		v4Group.POST("/coupons/:code/cancel", couponMiddleware.UserIDMiddleware(), couponHandlerV4.CancelCoupon)
		v4Group.GET("/coupons/user", couponMiddleware.UserIDMiddleware(), couponHandlerV4.FindCouponsByUserID)
		v4Group.GET("/coupons/policy/:policyCode", couponHandlerV4.FindCouponsByCouponPolicyCode)

		v4Group.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("couponSwaggerV4"),
			ginSwagger.URL("/api/v4/swagger/doc.json"),
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
