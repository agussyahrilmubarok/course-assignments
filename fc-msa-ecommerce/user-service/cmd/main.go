package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"example.com/user-service/pkg/config"
	"example.com/user-service/pkg/logger"
	"example.com/user-service/pkg/mongo"
	"example.com/user-service/pkg/server"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	log := logger.NewZerolog(cfg)

	mongoClient, err := mongo.Connect(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect MongoDB")
	}

	app := server.NewFiber(cfg)

	addr := cfg.App.Host + ":" + strconv.Itoa(cfg.App.Port)

	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()
	log.Info().Msgf("server started on %s", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}

	if err := mongoClient.Disconnect(ctx, log); err != nil {
		log.Error().Err(err).Msg("mongo disconnect error")
	}

	log.Info().Msg("server stopped gracefully")
}
