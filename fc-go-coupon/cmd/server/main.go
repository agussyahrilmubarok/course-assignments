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

	"example.com/internal/server"
	"example.com/pkg/config"
	"example.com/pkg/logger"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Path for config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic("cannot found config file")
	}

	log := logger.NewZerolog(cfg)

	r := server.NewRouter()

	serverAddr := fmt.Sprintf(":%v", cfg.Server.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		log.Printf("Server running on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("server forced to shutdown: %v", err)
	}

	log.Info().Msg("server stopped cleanly")
}
