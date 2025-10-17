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

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "example.com/account/cmd/server/docs"
	"example.com/account/internal/account"
	"example.com/account/pkg/discovery"
	"example.com/account/pkg/discovery/consul"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Account Service API
// @version 1.0
// @description This is an account service API.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1/accounts

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configFlag := flag.String("config", "configs/account.yaml", "Path to config file")
	flag.Parse()

	cfg, err := account.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := account.NewZerolog(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	db, err := account.NewPostgres(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	// Auto migrate (development only)
	// if cfg.App.Env != "production" {
	// 	if err := db.AutoMigrate(&account.User{}); err != nil {
	// 		logger.Fatal().Err(err).Msg("AutoMigrate failed")
	// 		os.Exit(1)
	// 	}
	// 	logger.Info().Msg("AutoMigrate executed")
	// }

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

	store := account.NewStore(db, logger)
	service := account.NewService(cfg, logger)
	handler := account.NewHandler(store, service, logger)
	customMiddleware := account.NewCustomMiddleware(service)

	e := echo.New()
	e.HideBanner = false
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &account.CustomValidator{Validator: validator.New()}

	v1 := e.Group("/api/v1/accounts")
	{
		v1.GET("/healthz", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"status":  "ok",
				"service": "account-service",
			})
		})
		v1.GET("/swagger/*", echoSwagger.WrapHandler)

		v1.POST("/sign-up", handler.SignUp)
		v1.POST("/sign-in", handler.SignIn, middleware.RateLimiterWithConfig(customMiddleware.RateLimiterConfig()))
		v1.POST("/validate", handler.Validate)
		v1.GET("/me", handler.GetMe, customMiddleware.Auth())
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
