package restV1

import (
	"example.com.backend/internal/config"
	"example.com.backend/internal/repos"

	"example.com.backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	backendMiddleware "example.com.backend/internal/middleware"
	_ "example.com.backend/internal/rest/v1/docs"
	handlerV1 "example.com.backend/internal/rest/v1/handler"
	usecaseV1 "example.com.backend/internal/rest/v1/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Backerhub Backend API V1
// @version 1.0
// @description Backerhub Backend API V1
// @BasePath /api/v1
func Register(ginEngine *gin.Engine, cfg *config.Config, db *gorm.DB) {
	r := ginEngine

	userRepo := repos.NewUserRepository(db)
	campaignRepo := repos.NewCampaignRepository(db)
	campaignImageRepo := repos.NewCampaignImageRepository(db)
	transactionRepo := repos.NewTransactionRepository(db)

	uploadService := service.NewUploadService()
	jwtService := service.NewJwtService(cfg)
	midtransService := service.NewMidtransService(&cfg.Midtrans)

	authUseCase := usecaseV1.NewAuthUseCaseV1(userRepo, jwtService)
	userUseCase := usecaseV1.NewUserUseCaseV1(userRepo)
	campaignUseCase := usecaseV1.NewCampaignUseCaseV1(campaignRepo, campaignImageRepo, uploadService)
	transactionUseCase := usecaseV1.NewTransactionUseCaseV1(transactionRepo, midtransService, userRepo, campaignRepo)

	handler := handlerV1.NewHandler()
	authHandler := handlerV1.NewAuthHandlerV1(authUseCase)
	profileHandler := handlerV1.NewProfileHandlerV1(userUseCase)
	campaignHandler := handlerV1.NewCampaignHandlerV1(campaignUseCase)
	transactionHandler := handlerV1.NewTransactionHanderV1(transactionUseCase)
	paymentHandler := handlerV1.NewPaymentHanderV1(transactionUseCase)

	jwtAuthMiddleware := backendMiddleware.JwtAuthMiddleware(jwtService)

	apiV1 := r.Group("/api/v1")
	apiV1.GET("", handler.V1)

	authApi := apiV1.Group("/auth")
	{
		authApi.POST("/sign-up", authHandler.SignUp)
		authApi.POST("/sign-in", authHandler.SignIn)
	}

	profileApi := apiV1.Group("/profiles")
	{
		profileApi.GET("/me", jwtAuthMiddleware, profileHandler.GetMe)
	}

	campaignApi := apiV1.Group("/campaigns")
	{
		campaignApi.GET("", campaignHandler.FindAll)
		campaignApi.GET("/top", campaignHandler.FindAllTop)
		campaignApi.GET("/:id", campaignHandler.FindByID)
		campaignApi.GET("/me", jwtAuthMiddleware, campaignHandler.FindAllByUser)
		campaignApi.GET("/:id/me", jwtAuthMiddleware, campaignHandler.FindByIDByUser)
		campaignApi.POST("", jwtAuthMiddleware, campaignHandler.CreateByUser)
		campaignApi.PUT("/:id", jwtAuthMiddleware, campaignHandler.UpdateByIDByUser)
		campaignApi.POST("/:id/images", jwtAuthMiddleware, campaignHandler.UploadImageByUser)
		campaignApi.DELETE("/:id", jwtAuthMiddleware, campaignHandler.DeleteByUser)
	}

	transactionApi := apiV1.Group("/transactions")
	{
		transactionApi.GET("/me", jwtAuthMiddleware, transactionHandler.FindAllByUser)
		transactionApi.GET("/campaign/:id", jwtAuthMiddleware, transactionHandler.FindAllByCampaign)
		transactionApi.GET("/:id", jwtAuthMiddleware, transactionHandler.FindByID)
		transactionApi.POST("/donation", jwtAuthMiddleware, transactionHandler.Donation)
	}

	paymentApi := apiV1.Group("/payments")
	{
		paymentApi.POST("/midtrans/notification", paymentHandler.MidtransPaymentNotification)
	}

	apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.InstanceName("backerhubAPIV1"),
		ginSwagger.URL("/api/v1/swagger/doc.json"),
	))
}
