package main

import (
	"context"
	"errors"
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
	"example.com/coupon-service/internal/instrument/logging"
	"example.com/coupon-service/internal/instrument/metrics"
	"example.com/coupon-service/internal/instrument/tracing"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	v1 "example.com/coupon-service/internal/api/v1"
	v2 "example.com/coupon-service/internal/api/v2"
	v3 "example.com/coupon-service/internal/api/v3"
	v4 "example.com/coupon-service/internal/api/v4"
)

func main() {
	cfgPath := flag.String("config", "config.yml", "Config filepath")
	flag.Parse()

	cfg, err := config.NewConfig(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logging.InitLogging(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logging.GetLogger().Sync()
	log := logging.GetLogger()

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

	traceExporter := tracing.NewZipkinExporter(cfg.Zipkin.Url)
	shutdownTrace := tracing.InitTraceProvider(ctx, cfg.Server.Name, traceExporter)
	tracing.NewTracer(cfg.Server.Name)
	defer shutdownTrace(ctx)

	e := echo.New()
	e.Use(middleware.TraceIDMiddleware())

	dummyHandler := dummy.NewHandler(e, pg, rdb)
	dummyHandler.RegisterDummyAPI()

	api := e.Group("/api")
	v1.RegisterAPIV1(api, pg)
	v2.RegisterAPIV2(api, pg)
	v3.RegisterAPIV3(api, pg, rdb)
	v4.RegisterAPIV4(api, cfg, pg, rdb)

	// START HTTP SERVER
	serverAddr := fmt.Sprintf(":%v", cfg.Server.Port)
	go func() {
		log.Info("starting echo server", zap.String("addr", serverAddr))
		if err := e.Start(serverAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("shutting down server due to error", zap.Error(err))
		}
	}()

	// START KAFKA CONSUMER
	consumer := v4.NewKafkaConsumer(cfg, pg, rdb)
	go func() {
		log.Info("starting kafka consumer...")
		if err := consumer.Start(ctx); err != nil {
			log.Error("kafka consumer stopped with error", zap.Error(err))
		}
	}()

	// START METRIC SERVER
	metricAddr := fmt.Sprintf(":%v", cfg.Metric.Port)
	go func() {
		log.Info("starting metric server", zap.String("addr", metricAddr))
		metricsServer := metrics.NewMetricServer(cfg)
		if err := metricsServer.Start(metricAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("shutting down metric server due to error", zap.Error(err))
		}
	}()

	// WAIT SIGNAL
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("received shutdown signal", zap.String("signal", sig.String()))

	// SHUTDOWN KAFKA CONSUMER
	log.Info("closing kafka consumer...")
	consumer.Close()
	log.Info("kafka consumer closed")

	// SHUTDOWN HTTP SERVER
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctxShutDown); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	log.Info("application exited gracefully")
}
