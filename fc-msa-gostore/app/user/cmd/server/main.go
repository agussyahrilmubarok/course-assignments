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

	"example.com/user/internal/user"
	"example.com/user/pkg/config"
	"example.com/user/pkg/logger"
	"example.com/user/pkg/mongodb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	_ "example.com/user/cmd/server/docs"
	swagger "github.com/swaggo/fiber-swagger"
)

// @title User Application API
// @version 1.0
// @description This is the API documentation for the User Application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/example
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api/v1/users
func main() {
	configFlag := flag.String("config", "configs/config.json", "Path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(*configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.GetLogger(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	mongodbFactory := mongodb.NewMongoFactory(cfg)
	mongoClient1, err := mongodbFactory.GetClient(cfg.MongoDB.URI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect mongodb: %v\n", err)
		os.Exit(1)
	}
	mongoDb1 := mongoClient1.Database(cfg.MongoDB.DbName)

	store := user.NewStore(mongoDb1, log)
	service := user.NewService(store, cfg, log)
	handler := user.NewHandler(service, log)

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(requestid.New())

	v1 := app.Group("/api/v1/users")

	v1.Post("/auth/signup", handler.SignUp)
	v1.Post("/auth/signin", handler.SignIn)
	v1.Get("/:id", handler.FindByID)

	app.Get("/swagger/*", swagger.FiberWrapHandler())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "user-app",
		})
	})

	go func() {
		log.Info().Msgf("Server running at http://%s:%d", cfg.App.Host, cfg.App.Port)
		if err := app.Listen(fmt.Sprintf(":%d", cfg.App.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Info().Msg("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		log.Fatal().Err(err).Msg("Graceful shutdown failed")
	}

	if err := mongodbFactory.CloseAll(); err != nil {
		log.Warn().Err(err).Msg("Failed to close MongoDB connections cleanly")
	}

	log.Info().Msg("Server stopped gracefully")
}
