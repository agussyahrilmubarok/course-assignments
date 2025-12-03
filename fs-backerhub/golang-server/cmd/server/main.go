package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiV1 "example.com/backend/api/v1"
	"example.com/backend/internal/controller"
	"example.com/backend/internal/domain"
	"example.com/backend/internal/helper"
	"example.com/backend/internal/middleware"
	"example.com/backend/internal/repository"
	"example.com/backend/internal/service"
	"example.com/backend/pkg/config"
	"example.com/backend/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "backerhub.json", "Path to config file")
	flag.Parse()

	// Load configuration file
	cfg := config.Load(*configPath)

	// Initialize logger
	log := logger.NewZeroLogger(cfg.Logging)

	// Connect Postgres Server
	db, err := config.NewGormPostgres(cfg.PostgresSQL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect postgres")
	}

	// Auto migrate if not production level
	if lvl := cfg.App.Level; lvl != "production" && lvl != "prod" {
		if err := db.AutoMigrate(
			&domain.User{},
			&domain.Campaign{},
			&domain.CampaignImage{},
			&domain.Transaction{},
		); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate postgres")
		}
	}

	// Define DI Composition
	userRepo := repository.NewUserRepository(db)
	campaignRepo := repository.NewCampaignRepository(db)
	campaignImageRepo := repository.NewCampaignImageRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	uploadService := service.NewUploadService(log)
	userService := service.NewUserService(userRepo, log)
	campaignService := service.NewCampaignService(campaignRepo, campaignImageRepo, log)
	transactionService := service.NewTransactionService(transactionRepo, log)
	jwtService := service.NewJwtService(cfg.App, log)
	midtransService := service.NewMidtransService(cfg.Midtrans, log)

	controller.NewBaseController(log)
	homeController := controller.NewHomeController()
	loginController := controller.NewLoginController(userService, log)
	dashboardController := controller.NewDashboardController()
	userController := controller.NewUserController(userService, uploadService, log)
	campaignController := controller.NewCampaignController(campaignService, userService, uploadService, log)
	transactionController := controller.NewTransactionController(transactionService, userService, campaignService, log)

	jwtAuthMiddleware := middleware.JwtAuth(jwtService)

	// Define Gin as engine server
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(sessions.Sessions(cfg.App.Name, cookie.NewStore([]byte(cfg.App.Name))))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.App.Clients,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Web HTML render
	r.HTMLRender = helper.LoadTemplate("./public/templates")
	r.Static("/assets", "./public/assets")
	r.Static("/uploads", "./public/uploads")
	r.StaticFile("/favicon.ico", "./public/assets/favicon.ico")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/", homeController.Index)
	r.GET("/login", loginController.Index)
	r.POST("/login", loginController.Login)

	adminDashboard := r.Group("/dashboard")
	{
		adminDashboard.GET("/", dashboardController.Index)

		users := adminDashboard.Group("/users")
		{
			users.GET("/", userController.Index)
			users.GET("/add", userController.Add)
			users.POST("/add", userController.Create)
			users.GET("/:id/edit", userController.Edit)
			users.POST("/:id/edit", userController.Update)
			users.GET("/:id/avatar", userController.Avatar)
			users.POST("/:id/avatar", userController.Upload)
			users.GET("/:id/delete", userController.Delete)
		}

		campaigns := adminDashboard.Group("/campaigns")
		{
			campaigns.GET("/", campaignController.Index)
			campaigns.GET("/add", campaignController.Add)
			campaigns.POST("/add", campaignController.Create)
			campaigns.GET("/:id/show", campaignController.Show)
			campaigns.GET("/:id/edit", campaignController.Edit)
			campaigns.POST("/:id/edit", campaignController.Update)
			campaigns.GET("/:id/image", campaignController.Image)
			campaigns.POST("/:id/image", campaignController.Upload)
			campaigns.GET("/:id/delete", campaignController.Delete)
		}

		transactions := adminDashboard.Group("/transactions")
		{
			transactions.GET("/", transactionController.Index)
			transactions.GET("/add", transactionController.Add)
			transactions.POST("/add", transactionController.Create)
			transactions.GET("/:id/show", transactionController.Show)
			transactions.GET("/:id/edit", transactionController.Edit)
			transactions.POST("/:id/edit", transactionController.Update)
			transactions.GET("/:id/delete", transactionController.Delete)
		}

		adminDashboard.GET("/logout", dashboardController.Logout)
	}

	// Register API v1
	v1 := r.Group("/api/v1")
	apiV1.RegisterAuthRoutes(v1, userRepo, jwtService, log)
	apiV1.RegisterProfileRoutes(v1, userRepo, jwtAuthMiddleware, log)
	apiV1.RegisterCampaignRoutes(v1, campaignRepo, campaignImageRepo, uploadService, jwtAuthMiddleware, log)
	apiV1.RegisterTransactionRoutes(v1, transactionRepo, midtransService, userRepo, campaignRepo, jwtAuthMiddleware, log)
	apiV1.RegisterPaymentRoutes(v1, transactionRepo, midtransService, userRepo, campaignRepo, log)

	addr := fmt.Sprintf(":%v", cfg.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Info().Msgf("%s is running on port %s", cfg.App.Name, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	} else {
		log.Info().Msg("server exited gracefully")
	}
}
