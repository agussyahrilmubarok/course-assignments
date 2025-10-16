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

	"example.com/internal/account"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load config
	configFlag := flag.String("config", "configs/account.yaml", "Path to config file")
	flag.Parse()

	cfg, err := account.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Init logger
	logger, err := account.NewZerolog()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	// Init DB
	db, err := account.NewPostgres(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	// Auto migrate (development only)
	if cfg.App.Env != "production" {
		if err := db.AutoMigrate(&account.User{}); err != nil {
			logger.Fatal().Err(err).Msg("AutoMigrate failed")
			os.Exit(1)
		}
		logger.Info().Msg("AutoMigrate executed")
	}

	// Init Echo server
	e := echo.New()
	e.HideBanner = false
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &account.CustomValidator{Validator: validator.New()}

	// Register routes
	store := account.NewStore(db, logger)
	service := account.NewService()
	handler := account.NewHandler(store, service, logger)

	v1 := e.Group("/api/v1/accounts")
	{
		v1.GET("/healthz", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"status":  "ok",
				"service": "account-service",
			})
		})

		v1.POST("/sign-up", handler.SignUp)
		v1.POST("/sign-in", handler.SignIn)
		v1.POST("/validate", handler.Validate)
		v1.GET("/me", handler.GetMe)
	}

	// 7. Start server in goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		logger.Info().Msgf("Server running at %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
			os.Exit(1)
		}
	}()

	// 8. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // Wait for termination signal
	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Graceful shutdown failed")
		os.Exit(1)
	}

	logger.Info().Msg("Server stopped gracefully")
}
