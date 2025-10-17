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

	"example.com/catalog/internal/catalog"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "example.com/catalog/cmd/server/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Catalog Service API
// @version 1.0
// @description This is an catalog service API.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8082
// @BasePath /api/v1/catalogs
func main() {
	configFlag := flag.String("config", "configs/catalog.yaml", "Path to config file")
	flag.Parse()

	cfg, err := catalog.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := catalog.NewZerolog(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	db, err := catalog.NewPostgres(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	if err := db.AutoMigrate(&catalog.Product{}); err != nil {
		logger.Fatal().Err(err).Msg("AutoMigrate failed")
		os.Exit(1)
	}

	rdb, err := catalog.NewRedis(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to cache")
		os.Exit(1)
	}

	store := catalog.NewStore(db, rdb, logger)
	handler := catalog.NewHandler(store, logger)

	e := echo.New()
	e.HideBanner = false
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &catalog.CustomValidator{Validator: validator.New()}

	v1 := e.Group("/api/v1/catalogs")
	{
		v1.GET("/healthz", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"status":  "ok",
				"service": "catalog-service",
			})
		})
		v1.GET("/swagger/*", echoSwagger.WrapHandler)

		v1.GET("/products", handler.GetProducts)
		v1.GET("/products/:id", handler.GetProduct)
		v1.POST("/products/stocks/reverse", handler.ReverseProductStock)
		v1.POST("/products/stocks/release", handler.ReleaseProductStock)
		v1.GET("/products/stocks/:id", handler.GetProductStock)
	}

	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		logger.Info().Msgf("Server running at %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit 
	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Graceful shutdown failed")
		os.Exit(1)
	}

	logger.Info().Msg("Server stopped gracefully")
}
