package main

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/handler"
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
	categoryRepository := repository.NewCategoryRepository(db, log.Logger)
	categoryService := service.NewCategoryService(categoryRepository, log.Logger)
	categoryHandler := handler.NewCategoryHandler(categoryService, log.Logger)

	// Setup Gin
	router := gin.Default()

	// Register Swagger docs
	router.GET("/swagger/*any", gin.WrapH(httpSwagger.WrapHandler))

	// Register sample health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// TODO: Register API routes here (product, user, order, etc.)
	categoryRoute := router.Group("/api/categories")
	{
		categoryRoute.GET("", categoryHandler.GetAll)
		categoryRoute.GET("/:id", categoryHandler.GetByID)
		categoryRoute.POST("", categoryHandler.Create)
		categoryRoute.PUT("/:id", categoryHandler.Update)
		categoryRoute.DELETE("/:id", categoryHandler.Delete)
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
