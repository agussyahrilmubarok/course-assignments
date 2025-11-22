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

	"example.com/coupon-service/internal/api/dummy"
	v1 "example.com/coupon-service/internal/api/v1"
	"example.com/coupon-service/internal/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	cfgPath := flag.String("config", "config.yml", "Config filepath")
	flag.Parse()

	cfg, err := config.NewConfig(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logging, err := config.NewLogging(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize config: %v\n", err)
		os.Exit(1)
	}
	defer logging.Logger.Sync()
	logger := logging.Logger

	ctx := context.Background()

	pg, err := config.NewPostgres(ctx, cfg)
	if err != nil {
		logger.Fatal("failed to connect postgres", zap.Error(err))
	}
	defer pg.Close()

	rdb, err := config.NewRedis(ctx, cfg)
	if err != nil {
		logger.Fatal("failed to connect postgres", zap.Error(err))
	}
	defer rdb.Close()

	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	dummyHandler := dummy.NewHandler(pg, rdb, logger)
	e.GET("/init-dummy-v1", dummyHandler.InitDummyV1)
	e.GET("/clean-dummy-v1", dummyHandler.CleanDummyV1)
	e.GET("/check-quantity-v1/:policy_code", dummyHandler.CheckQuantityV1)
	e.GET("/init-dummy-v2", dummyHandler.InitDummyV2)
	e.GET("/clean-dummy-v2", dummyHandler.CleanDummyV2)

	apiV1 := e.Group("/api/v1")
	v1.RegisterAPIV1(apiV1, pg, logger)

	serverAddr := ":8080"
	go func() {
		logger.Info("starting echo server", zap.String("addr", serverAddr))
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("shutting down server due to error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("shutting down server...", zap.String("signal", sig.String()))

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctxShutDown); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}
