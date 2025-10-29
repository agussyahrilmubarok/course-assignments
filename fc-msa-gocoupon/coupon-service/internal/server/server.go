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

	couponPolicyHandlerV1 := v1.NewCouponPolicyHandler(db)
	couponHandlerV1 := v1.NewCouponHandler()

	// V1 routes
	v1Group := apiRoute.Group("/v1")
	{
		v1Group.POST("/couponPolicies/dummy", couponPolicyHandlerV1.CreateCouponPolicyDummy)
		v1Group.GET("/couponPolicies", couponPolicyHandlerV1.SearchCouponPolicy)

		v1Group.GET("/coupons", couponHandlerV1.FindCoupon)

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
