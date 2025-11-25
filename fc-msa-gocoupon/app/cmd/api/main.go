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
	"example.com/coupon-service/internal/api/middleware"
	"example.com/coupon-service/internal/config"
	"example.com/coupon-service/internal/instrument"
	"example.com/coupon-service/internal/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	v1 "example.com/coupon-service/internal/api/v1"
	v2 "example.com/coupon-service/internal/api/v2"
)

func main() {
	cfgPath := flag.String("config", "config.yml", "Config filepath")
	flag.Parse()

	cfg, err := config.NewConfig(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.InitLogger(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.GetLogger().Sync()
	log := logger.GetLogger()

	ctx := context.Background()

	pg, err := config.NewPostgres(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to postgres: %v\n", err)
		os.Exit(1)
	}
	defer pg.Close()

	rdb, err := config.NewRedis(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to redis: %v\n", err)
		os.Exit(1)
	}
	defer rdb.Close()

	traceExporter := instrument.NewZipkinExporter(cfg.Zipkin.Url)
	shutdownTrace := instrument.InitTraceProvider(ctx, cfg.Server.Name, traceExporter)
	instrument.NewTracer(cfg.Server.Name)
	defer shutdownTrace(ctx)

	e := echo.New()

	e.Use(middleware.TraceIDMiddleware())

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	dummyHandler := dummy.NewHandler(pg, rdb)
	e.GET("/init-dummy-db", dummyHandler.InitDummyDB)
	e.GET("/clean-dummy-db", dummyHandler.CleanDummyDB)
	e.GET("/init-dummy-redis-db", dummyHandler.InitDummyRedisAndDB)
	e.GET("/clean-dummy-redis-db", dummyHandler.CleanDummyRedisAndDB)
	e.GET("/check-quantity/:policy_code", dummyHandler.CheckQuantity)

	api := e.Group("/api")
	v1.RegisterAPIV1(api, pg)
	v2.RegisterAPIV2(api, pg)

	serverAddr := ":8080"
	go func() {
		log.Info("starting echo server", zap.String("addr", serverAddr))
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down server due to error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("shutting down server...", zap.String("signal", sig.String()))

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctxShutDown); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited gracefully")
}
