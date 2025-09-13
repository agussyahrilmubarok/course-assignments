package main

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/handler"
	"ecommerce/internal/middleware"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "ecommerce/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Ecommerce API
// @version 1.0
// @description This is an Ecommerce API server.
// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Setup Zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Connect to PostgreSQL
	dsn := "host=localhost user=postgres password=yourpassword dbname=ecommerce port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Run AutoMigrate
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.Tag{},
		&domain.Product{},
		&domain.Cart{},
		&domain.CartItem{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
		&domain.Address{},
		&domain.Review{},
	); err != nil {
		log.Fatal().Err(err).Msg("auto migration failed")
	}

	log.Info().Msg("database connected and migrated successfully")

	// TODO: Register dependency injection
	userRepository := repository.NewUserRepository(db, log.Logger)
	categoryRepository := repository.NewCategoryRepository(db, log.Logger)
	tagRepository := repository.NewTagRepository(db, log.Logger)

	authService := service.NewAuthService(userRepository, log.Logger)
	userService := service.NewUserService(userRepository, log.Logger)
	categoryService := service.NewCategoryService(categoryRepository, log.Logger)
	tagService := service.NewTagService(tagRepository, log.Logger)

	authHandler := handler.NewAuthHandler(authService, log.Logger)
	userHandler := handler.NewUserHandler(userService, log.Logger)
	categoryHandler := handler.NewCategoryHandler(categoryService, log.Logger)
	tagHandler := handler.NewTagHandler(tagService, log.Logger)

	// Setup Gin
	router := gin.Default()

	// Register Swagger docs
	router.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))

	// Register sample health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// TODO: Register API routes here (product, user, order, etc.)
	authRoute := router.Group("/api/auth")
	{
		authRoute.POST("/register", authHandler.Register)
		authRoute.POST("/login", authHandler.Login)
		authRoute.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
	}

	userRoute := router.Group("/api/users", middleware.AuthRequired())
	{
		userRoute.GET("/me", userHandler.GetCurrentUser)
		userRoute.PUT("/me", userHandler.UpdateCurrentUser)
		userRoute.GET("", middleware.AdminOnly(), userHandler.GetAll)
		userRoute.DELETE("/:id", middleware.AdminOnly(), userHandler.Delete)
	}

	categoryRoute := router.Group("/api/categories")
	{
		categoryRoute.GET("", categoryHandler.GetAll)
		categoryRoute.GET("/:id", categoryHandler.GetByID)
	}
	categoryRoute.Use(middleware.AuthRequired(), middleware.AdminOnly())
	{
		categoryRoute.POST("", categoryHandler.Create)
		categoryRoute.PUT("/:id", categoryHandler.Update)
		categoryRoute.DELETE("/:id", categoryHandler.Delete)
	}

	tagRoute := router.Group("/api/tags")
	{
		tagRoute.GET("", tagHandler.GetAll)
		tagRoute.GET("/:id", tagHandler.GetByID)
	}
	tagRoute.Use(middleware.AuthRequired(), middleware.AdminOnly())
	{
		tagRoute.POST("", tagHandler.Create)
		tagRoute.PUT("/:id", tagHandler.Update)
		tagRoute.DELETE("/:id", tagHandler.Delete)
	}

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Run server in goroutine
	go func() {
		log.Info().Msg("server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server...")

	// Graceful shutdown context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	} else {
		log.Info().Msg("server exited properly")
	}
}
