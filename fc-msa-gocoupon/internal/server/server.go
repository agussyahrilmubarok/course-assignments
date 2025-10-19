package server

import (
	v1 "example.com/api/v1/handler"
	v2 "example.com/api/v2/handler"
	v3 "example.com/api/v3/handler"
	v4 "example.com/api/v4/handler"

	"github.com/gin-gonic/gin"

	_ "example.com/api/v1/docs"
	_ "example.com/api/v2/docs"
	_ "example.com/api/v3/docs"
	_ "example.com/api/v4/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")

	couponHandlerV1 := v1.NewCouponHandler()

	// V1 routes
	v1Group := api.Group("/v1")
	{
		v1Group.GET("/coupons", couponHandlerV1.GetCoupons)
		v1Group.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("swagger_v1"),
			ginSwagger.URL("/api/v1/swagger/doc.json"),
		))
	}

	couponHandlerV2 := v2.NewCouponHandler()

	// V2 routes
	v2Group := api.Group("/v2")
	{
		v2Group.GET("/coupons", couponHandlerV2.GetCoupons)
		v2Group.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("swagger_v2"),
			ginSwagger.URL("/api/v2/swagger/doc.json"),
		))
	}

	couponHandlerV3 := v3.NewCouponHandler()

	// V3 routes
	v3Group := api.Group("/v3")
	{
		v3Group.GET("/coupons", couponHandlerV3.GetCoupons)
		v3Group.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("swagger_v3"),
			ginSwagger.URL("/api/v3/swagger/doc.json"),
		))
	}

	couponHandlerV4 := v4.NewCouponHandler()

	// V4 routes
	v4Group := api.Group("/v4")
	{
		v4Group.GET("/coupons", couponHandlerV4.GetCoupons)
		v4Group.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName("swagger_v4"),
			ginSwagger.URL("/api/v4/swagger/doc.json"),
		))
	}

	return r
}
