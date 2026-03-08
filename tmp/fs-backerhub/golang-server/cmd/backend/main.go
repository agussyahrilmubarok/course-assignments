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

	"example.com.backend/internal/config"
	"example.com.backend/internal/controller"
	"example.com.backend/internal/domain"
	restV1 "example.com.backend/internal/rest/v1"
	restV2 "example.com.backend/internal/rest/v2"
	"example.com.backend/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	backendMiddleware "example.com.backend/internal/middleware"
)

func main() {
	cfgPath := flag.String("config", "configs/config.json", "Config filepath")
	flag.Parse()

	cfg, err := config.NewConfig(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Filepath, cfg.Logger.GelfAddr); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.GetLogger().Sync()

	log := logger.GetLogger()

	db, err := config.NewPostgres(&cfg.Postgres)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to postgres: %v\n", err)
		os.Exit(1)
	}

	if backendLevel := cfg.Backend.Level; backendLevel != "production" {
		if err := db.AutoMigrate(
			&domain.User{},
			&domain.Campaign{},
			&domain.CampaignImage{},
			&domain.Transaction{},
		); err != nil {
			fmt.Fprintf(os.Stderr, "failed to migrate postgres: %v\n", err)
			os.Exit(1)
		}
	}

	r := gin.New()
	r.Use(backendMiddleware.RequestIDMiddleware())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowClients,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	controller.Register(r, cfg, db)
	restV1.Register(r, cfg, db)
	restV2.Register(r)

	addr := fmt.Sprintf(":%v", cfg.Backend.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Sugar().Infof("%s is running on port %s", cfg.Backend.Name, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Sugar().Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Sugar().Infof("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Sugar().Error("server forced to shutdown", zap.Error(err))
	} else {
		log.Sugar().Errorf("server exited gracefully")
	}
}
