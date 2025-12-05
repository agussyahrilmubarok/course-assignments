package restV2

import (
	"github.com/gin-gonic/gin"

	_ "example.com.backend/internal/rest/v2/docs"
	handlerV2 "example.com.backend/internal/rest/v2/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Backerhub Backend API V2
// @version 2.0
// @description Backerhub Backend API V2
// @BasePath /api/v2
func Register(ginEngine *gin.Engine) {
	r := ginEngine

	handler := handlerV2.NewHandler()

	apiV2 := r.Group("/api/v2")
	apiV2.GET("", handler.V2)

	apiV2.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.InstanceName("backerhubAPIV2"),
		ginSwagger.URL("/api/v2/swagger/doc.json"),
	))
}
