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
	"traffic-control/account/internal/config"
	"traffic-control/account/internal/handler"
	"traffic-control/account/internal/repository"
	"traffic-control/account/internal/service"
	"traffic-control/account/pkg/logger"

	"github.com/labstack/echo/v4"
)

func main() {
	configFlag := flag.String("config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	cfg := config.LoadEnv(*configFlag)

	log := logger.NewLogger(cfg.Logging.Level, true)
	log.Info().Str("service", cfg.App.Name).Msg("Starting service initialization...")

	db, err := config.NewPostgres(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize PostgreSQL")
	}
	defer db.Close()

	e := echo.New()
	e.HideBanner = true

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": cfg.App.Name,
		})
	})

	userRepo := repository.NewUserRepository(db, log)
	userService := service.NewUserService(userRepo, log)
	userHandler := handler.NewUserHandler(userService, log)
	userHandler.RegisterRoutes(e)

	serverAddr := fmt.Sprintf(":%d", cfg.App.Port)
	go func() {
		log.Info().Str("addr", serverAddr).Msg("🚀 Starting HTTP server...")
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server startup failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Warn().Msg("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Forced to shutdown the server")
	}

	log.Info().Msg("Server exited cleanly")
}
