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

	"example.com/order/internal/order"
	"github.com/agussyahrilmubarok/gox/pkg/xconfig/xviper"
	"github.com/agussyahrilmubarok/gox/pkg/xdiscovery"
	"github.com/agussyahrilmubarok/gox/pkg/xdiscovery/xconsul"
	"github.com/agussyahrilmubarok/gox/pkg/xgorm"
	"github.com/agussyahrilmubarok/gox/pkg/xlogger/xzerolog"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"

	_ "example.com/order/cmd/server/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Order Service API
// @version 1.0
// @description This is an order service API.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8083
// @BasePath /api/v1/orders
func main() {
	configFlag := flag.String("config", "configs/config.yaml", "Path to config file")
	flag.Parse()

	vCfg, err := xviper.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	var cfg *order.Config
	if err := vCfg.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := xzerolog.NewLogger(cfg.Logger.Filepath, cfg.Logger.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	dbShards, err := openDBShard(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set up hybrid shards: %v\n", err)
		os.Exit(1)
	}

	instanceID := xdiscovery.GenerateInstanceID(cfg.App.Name)
	consulRegistry, err := xconsul.NewRegistry(cfg.Consul.Address)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to register consul discovery")
		os.Exit(1)
	}

	ctx := context.Background()
	if err := consulRegistry.Register(ctx, instanceID, cfg.App.Name, fmt.Sprintf("%v:%d", cfg.App.Host, cfg.App.Port)); err != nil {
		logger.Fatal().Err(err).Msg("Failed to register consul discovery")
		os.Exit(1)
	}

	go func() {
		for {
			if err := consulRegistry.ReportHealthyState(instanceID, cfg.App.Name); err != nil {
				logger.Info().Msg("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer consulRegistry.Deregister(ctx, instanceID, cfg.App.Name)

	store := order.NewStore(dbShards, order.NewShardRouter(2), logger)
	client := order.NewClient(cfg, logger)
	service := order.NewService(cfg, store, client, logger)
	handler := order.NewHandler(store, service, logger)

	e := echo.New()
	e.HideBanner = false
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &order.CustomValidator{Validator: validator.New()}

	v1 := e.Group("/api/v1/orders")
	{
		v1.GET("/healthz", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"status":  "ok",
				"service": "order-service",
			})
		})
		v1.GET("/swagger/*", echoSwagger.WrapHandler)

		v1.POST("/flash", handler.CreateFlashOrder)
		v1.POST("/cancel", handler.CancelOrder)
		v1.GET("/:id", handler.GetOrder)
	}

	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		logger.Info().Msgf("Server running at %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Graceful shutdown failed")
		os.Exit(1)
	}

	logger.Info().Msg("Server stopped gracefully")
}

func openDBShard(cfg *order.Config) ([]*gorm.DB, error) {
	var shards []*gorm.DB

	postgresDb, err := xgorm.NewGorm("postgres", fmt.Sprintf(
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
		return nil, err
	}
	shards = append(shards, postgresDb)

	mysqlDb, err := xgorm.NewGorm("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DbName,
	), &xgorm.Options{
		MaxOpenConns:    cfg.MySQL.MaxOpenConns,
		MaxIdleConns:    cfg.MySQL.MaxIdleConns,
		ConnMaxLifetime: cfg.MySQL.ConnMaxLifetime,
	})
	if err != nil {
		return nil, err
	}
	shards = append(shards, mysqlDb)

	for _, db := range shards {
		if err := db.AutoMigrate(&order.Order{}, &order.OrderItem{}); err != nil {
			return nil, err
		}
	}

	return shards, nil
}
