package v1

import (
	handlerV1 "example.com/backend/api/v1/handler"
	usecaseV1 "example.com/backend/api/v1/usecase"
	"example.com/backend/internal/repository"
	"example.com/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RegisterAuthRoutes(
	r *gin.RouterGroup,
	userRepo repository.IUserRepository,
	jwtService service.IJwtService,
	log zerolog.Logger,
) {
	authUseCase := usecaseV1.NewAuthUseCaseV1(userRepo, jwtService, log)
	authHandler := handlerV1.NewAuthHandlerV1(authUseCase, log)

	v1 := r.Group("/auth")
	{
		v1.POST("/sign-up", authHandler.SignUp)
		v1.POST("/sign-in", authHandler.SignIn)
	}
}

func RegisterProfileRoutes(
	r *gin.RouterGroup,
	userRepo repository.IUserRepository,
	jwtAuthMiddleware gin.HandlerFunc,
	log zerolog.Logger,
) {
	userUseCase := usecaseV1.NewUserUseCaseV1(userRepo, log)
	profileHandler := handlerV1.NewProfileHandlerV1(userUseCase, log)

	v1 := r.Group("/profiles")
	{
		v1.GET("/me", jwtAuthMiddleware, profileHandler.GetMe)
	}
}

func RegisterCampaignRoutes(
	r *gin.RouterGroup,
	campaignRepo repository.ICampaignRepository,
	campaignImageRepo repository.ICampaignImageRepository,
	uploadService service.IUploadService,
	jwtAuthMiddleware gin.HandlerFunc,
	log zerolog.Logger,
) {
	campaignUseCase := usecaseV1.NewCampaignUseCaseV1(campaignRepo, campaignImageRepo, uploadService, log)
	campaignHandler := handlerV1.NewCampaignHandlerV1(campaignUseCase, log)

	v1 := r.Group("/campaigns")
	{
		v1.GET("", campaignHandler.FindAll)
		v1.GET("/top", campaignHandler.FindAllTop)
		v1.GET("/:id", campaignHandler.FindByID)

		v1.GET("/me", jwtAuthMiddleware, campaignHandler.FindAllByUser)
		v1.GET("/:id/me", jwtAuthMiddleware, campaignHandler.FindByIDByUser)
		v1.POST("", jwtAuthMiddleware, campaignHandler.CreateByUser)
		v1.PUT("/:id", jwtAuthMiddleware, campaignHandler.UpdateByIDByUser)
		v1.POST("/:id/images", jwtAuthMiddleware, campaignHandler.UploadImageByUser)
		v1.DELETE("/:id", jwtAuthMiddleware, campaignHandler.DeleteByUser)
	}
}

func RegisterTransactionRoutes(
	r *gin.RouterGroup,
	transactionRepo repository.ITransactionRepository,
	midtransService service.IMidtransService,
	userRepo repository.IUserRepository,
	campaignRepo repository.ICampaignRepository,
	jwtAuthMiddleware gin.HandlerFunc,
	log zerolog.Logger,
) {
	transactionUseCase := usecaseV1.NewTransactionUseCaseV1(transactionRepo, midtransService, userRepo, campaignRepo, log)
	transactionHandler := handlerV1.NewTransactionHanderV1(transactionUseCase, log)

	v1 := r.Group("/transactions")
	{
		v1.GET("/me", jwtAuthMiddleware, transactionHandler.FindAllByUser)
		v1.GET("/campaign/:id", jwtAuthMiddleware, transactionHandler.FindAllByCampaign)
		v1.GET("/:id", jwtAuthMiddleware, transactionHandler.FindByID)
		v1.POST("/donation", jwtAuthMiddleware, transactionHandler.Donation)
	}
}

func RegisterPaymentRoutes(
	r *gin.RouterGroup,
	transactionRepo repository.ITransactionRepository,
	midtransService service.IMidtransService,
	userRepo repository.IUserRepository,
	campaignRepo repository.ICampaignRepository,
	log zerolog.Logger,
) {
	transactionUseCase := usecaseV1.NewTransactionUseCaseV1(transactionRepo, midtransService, userRepo, campaignRepo, log)
	paymentHandler := handlerV1.NewPaymentHanderV1(transactionUseCase, log)

	v1 := r.Group("/payments")
	{
		v1.POST("/midtrans/notification", paymentHandler.MidtransPaymentNotification)
	}
}
