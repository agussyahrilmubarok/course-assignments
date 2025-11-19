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

	"github.com/agussyahrilmubarok/gox/pkg/xconfig/xviper"
	"github.com/agussyahrilmubarok/gox/pkg/xdiscovery"
	"github.com/agussyahrilmubarok/gox/pkg/xdiscovery/xconsul"
	"github.com/agussyahrilmubarok/gox/pkg/xgorm"
	"github.com/agussyahrilmubarok/gox/pkg/xlogger/xzerolog"

	"example.com/account/internal/account"
	"example.com/account/internal/config"
	"example.com/account/internal/server"
)

// @title Account Service API
// @version 1.0
// @description This is an account service API.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1/accounts

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configFlag := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	vCfg, err := xviper.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	var cfg *config.Config
	if err := vCfg.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := xzerolog.NewLogger(cfg.Logger.Filepath, cfg.Logger.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	db, err := xgorm.NewGorm("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DbName,
		cfg.Postgres.SslMode,
	), &xgorm.Options{
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
		os.Exit(1)
	}

	if err := db.AutoMigrate(&account.User{}); err != nil {
		logger.Fatal().Err(err).Msg("auto migrate failed")
		os.Exit(1)
	}

	instanceID := xdiscovery.GenerateInstanceID(cfg.App.Name)
	consulRegistry, err := xconsul.NewRegistry(cfg.Consul.Address)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to register consul discovery")
		os.Exit(1)
	}

	ctx := context.Background()
	if err := consulRegistry.Register(ctx, instanceID, cfg.App.Name, fmt.Sprintf("%v:%d", cfg.Http.Host, cfg.Http.Port)); err != nil {
		logger.Fatal().Err(err).Msg("failed to register consul discovery")
		os.Exit(1)
	}

	go func() {
		for {
			if err := consulRegistry.ReportHealthyState(instanceID, cfg.App.Name); err != nil {
				logger.Info().Msg("failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer consulRegistry.Deregister(ctx, instanceID, cfg.App.Name)

	httpServer := server.NewHttpServer(cfg, logger, db)

	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("failed to start server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		logger.Fatal().Err(err).Msg("graceful shutdown failed")
		os.Exit(1)
	}

	logger.Info().Msg("server stopped gracefully")
}
