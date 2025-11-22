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

	"example.com/coupon/internal/server"
	"example.com/coupon/pkg/config"
	"example.com/coupon/pkg/instrument"
)

func main() {
	configFlag := flag.String("config", "configs/config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := config.NewZerolog(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	traceExporter := instrument.NewOTLPExporter(ctx, fmt.Sprintf("%s:%v", cfg.OTEL.Host, cfg.OTEL.Port))
	shutdownTrace := instrument.InitTraceProvider(ctx, cfg.App.Name, traceExporter)
	defer shutdownTrace(ctx)

	logger.Info().Msg("starting server...")

	httpServer := server.NewHttpServer(cfg, logger)
	go func() {
		if err := httpServer.Run(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server failed: %v\n", err)
			os.Exit(1)
		}
	}()

	metricsServer := instrument.NewMetricServer(cfg)
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "metrics server failed: %v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err := server.CouponKafkaConsumerV4.ConsumeCouponIssueRequest(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "kafka consumer failed: %v\n", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info().Msg("shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := metricsServer.Shutdown(ctxShutdown); err != nil {
		logger.Fatal().Err(err).Msg("graceful shutdown failed")
	}

	<-ctxShutdown.Done()
	logger.Info().Msg("server stopped gracefully")
}
