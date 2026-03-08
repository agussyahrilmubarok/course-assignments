package main

import (
	"context"
	"finly-go/configs"
	"finly-go/internal/server"
	"finly-go/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if os.Getenv("APP_PORT") == "" || os.Getenv("MONGODB_URI") == "" || os.Getenv("AUTH_SECRET") == "" {
		if err := godotenv.Load(); err != nil {
			logger.Log.Warn("No .env file found, fallback to system environment")
		} else {
			logger.Log.Info("Loaded environment variables from .env")
		}
	}

	app, err := server.NewApp()
	if err != nil {
		logger.Log.Fatal("Error initializing the app", zap.Error(err))
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		// Wait for signal
		sigReceived := <-c
		logger.Log.Info("Signal received, shutting down gracefully", zap.String("signal", sigReceived.String()))

		// Context for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the server
		if err := app.GetApp().ShutdownWithContext(ctx); err != nil {
			logger.Log.Fatal("Error during server shutdown", zap.Error(err))
		}

		// Close resources (e.g., database connections)
		app.Close()
		logger.Log.Info("Server shut down successfully")
	}()

	// Start the server
	port := configs.GetAppPort()
	logger.Log.Info("Starting server", zap.String("port", port))
	if err := app.GetApp().Listen(port); err != nil {
		logger.Log.Fatal("Error starting server", zap.Error(err))
	}
}
