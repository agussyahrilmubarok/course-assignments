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
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := config.NewZerolog(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	traceExporter := instrument.NewOTLPExporter(ctx, fmt.Sprintf("%s:%v", cfg.OTEL.Host, cfg.OTEL.Port))
	shutdownTrace := instrument.InitTraceProvider(ctx, cfg.App.Name, traceExporter)
	defer shutdownTrace(ctx)

	srv := server.NewGinRouter(cfg, logger)
	go func() {
		logger.Info().Msg("Starting server on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed")
		}
	}()

	metricsServer := instrument.NewMetricServer(cfg)
	go func() {
		logger.Info().Msg("Starting server on " + metricsServer.Addr)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Info().Msg("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Fatal().Err(err).Msg("Graceful shutdown failed")
	}

	if err := metricsServer.Shutdown(ctxShutdown); err != nil {
		logger.Fatal().Err(err).Msg("Graceful shutdown failed")
	}

	logger.Info().Msg("Server stopped gracefully")
}
