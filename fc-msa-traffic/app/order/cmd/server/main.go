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

	"example.com/order/internal/order"
	"example.com/order/pkg/discovery"
	"example.com/order/pkg/discovery/consul"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "example.com/order/cmd/server/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Order Service API
// @version 1.0
// @description This is an order service API.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8083
// @BasePath /api/v1/orders
func main() {
	configFlag := flag.String("config", "configs/config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := order.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := order.NewZerolog(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	dbShards, err := order.NewDBShard(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set up hybrid shards: %v\n", err)
		os.Exit(1)
	}

	instanceID := discovery.GenerateInstanceID(cfg.App.Name)
	consulRegistry, err := consul.NewRegistry(cfg.Consul.Address)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to register consul discovery")
		os.Exit(1)
	}

	ctx := context.Background()
	if err := consulRegistry.Register(ctx, instanceID, cfg.App.Name, fmt.Sprintf("%v:%d", cfg.App.Host, cfg.App.Port)); err != nil {
		logger.Fatal().Err(err).Msg("Failed to register consul discovery")
		os.Exit(1)
	}

	go func() {
		for {
			if err := consulRegistry.ReportHealthyState(instanceID, cfg.App.Name); err != nil {
				logger.Info().Msg("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer consulRegistry.Deregister(ctx, instanceID, cfg.App.Name)

	store := order.NewStore(dbShards, order.NewShardRouter(2), logger)
	client := order.NewClient(cfg, logger)
	service := order.NewService(cfg, store, client, logger)
	handler := order.NewHandler(store, service, logger)

	e := echo.New()
	e.HideBanner = false
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &order.CustomValidator{Validator: validator.New()}

	v1 := e.Group("/api/v1/orders")
	{
		v1.GET("/healthz", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"status":  "ok",
				"service": "order-service",
			})
		})
		v1.GET("/swagger/*", echoSwagger.WrapHandler)

		v1.POST("/flash", handler.CreateFlashOrder)
		v1.POST("/cancel", handler.CancelOrder)
		v1.GET("/:id", handler.GetOrder)
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
