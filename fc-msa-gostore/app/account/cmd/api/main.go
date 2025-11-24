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

	"example.com/account-service/internal/account"
	"example.com/account-service/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfgFlag := flag.String("config", "config.yml", "Config filepath")
	flag.Parse()

	cfg, err := config.NewConfig(*cfgFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration file: %v\n", err)
		os.Exit(1)
	}

	logging, err := config.NewLogging(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logging server: %v\n", err)
		os.Exit(1)
	}
	logger := logging.Logger

	mongodb, err := config.NewMongoDB(context.Background(), cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect mongodb server: %v\n", err)
		os.Exit(1)
	}

	repository := account.NewRepository(mongodb.DB, logger)
	service := account.NewService(repository, logger)
	handler := account.NewHandler(service, logger)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	apiV1 := router.Group("/api/v1/accounts")
	apiV1.POST("/sign-up", handler.SignUp)
	apiV1.POST("/sign-in", handler.SignIn)
	apiV1.GET("/:id", handler.FindByID)

	srvPort := fmt.Sprintf(":%v", cfg.Server.Port)
	srv := &http.Server{
		Addr:    srvPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server forced to shutdown: %s\n", err)
	}

	logger.Println("server exiting cleanly")
}
