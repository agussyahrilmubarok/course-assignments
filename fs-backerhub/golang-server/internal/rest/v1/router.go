package restV1

import (
	"example.com.backend/internal/config"
	"example.com.backend/internal/repos"

	"example.com.backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	handlerV1 "example.com.backend/internal/rest/v1/handler"
	usecaseV1 "example.com.backend/internal/rest/v1/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "example.com.backend/internal/rest/v1/docs"
)

// @title Backerhub Backend API V1
// @version 1.0
// @description Backerhub Backend API V1
// @BasePath /api/v1
func Register(ginEngine *gin.Engine, cfg *config.Config, db *gorm.DB) {
	r := ginEngine

	userRepo := repos.NewUserRepository(db)
	// campaignRepo := repos.NewCampaignRepository(db)
	// campaignImageRepo := repos.NewCampaignImageRepository(db)
	// transactionRepo := repos.NewTransactionRepository(db)

	// uploadService := service.NewUploadService()
	jwtService := service.NewJwtService(cfg)
	// midtransService := service.NewMidtransService(&cfg.Midtrans)

	authUseCase := usecaseV1.NewAuthUseCaseV1(userRepo, jwtService)
	// userUseCase := usecaseV1.NewUserUseCaseV1(userRepo)
	// campaignUseCase := usecaseV1.NewCampaignUseCaseV1(campaignRepo, campaignImageRepo, uploadService)
	// transactionUseCase := usecaseV1.NewTransactionUseCaseV1(transactionRepo, midtransService, userRepo, campaignRepo)

	authHandler := handlerV1.NewAuthHandlerV1(authUseCase)

	apiV1 := r.Group("/api/v1")

	authApi := apiV1.Group("/auth")
	{
		authApi.POST("/sign-up", authHandler.SignUp)
		authApi.POST("/sign-in", authHandler.SignIn)
	}

	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.InstanceName("backerhubAPIV1"),
		ginSwagger.URL("/api/v1/swagger/doc.json"),
	))
}
