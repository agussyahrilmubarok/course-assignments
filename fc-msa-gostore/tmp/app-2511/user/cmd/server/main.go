package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/api/users/internal/user"
	"github.com/agussyahrilmubarok/gox/pkg/xlogger/xzerolog"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "example.com/api/users/cmd/server/docs"
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
	appHost := "localhost"
	appPort := 8081
	mongoUri := "mongodb://root:secret@localhost:27017"

	logger, err := xzerolog.NewLogger("./logs/app.log", "info")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	mongoDb := setupMongoDB(mongoUri)
	store := user.NewStore(mongoDb, logger)
	service := user.NewServie(store, logger)
	handler := user.NewHandler(service, logger)

	app := fiber.New()

	v1 := app.Group("/api/v1/users")
	{
		v1.Post("/auth/signup", handler.SignUp)
		v1.Post("/auth/signin", handler.SignIn)
	}

	app.Get("/swagger/*", swagger.FiberWrapHandler())
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	go func() {
		logger.Info().Msgf("server running at http://%s:%d", appHost, appPort)
		if err := app.Listen(fmt.Sprintf(":%d", appPort)); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info().Msg("shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctxShutdown); err != nil {
		logger.Fatal().Err(err).Msg("graceful shutdown failed")
	}

	logger.Info().Msg("server stopped gracefully")
}

func setupMongoDB(mongoUri string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect mongo: %v\n", err)
		os.Exit(1)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping connect mongo: %v\n", err)
		os.Exit(1)
	}

	db := mongoClient.Database("user_db")

	users := db.Collection("users")
	usersIndexModel := mongo.IndexModel{
		Keys: bson.M{"email": 1},
		Options: options.Index().
			SetUnique(true).
			SetName("unique_email"),
	}
	_, err = users.Indexes().CreateOne(ctx, usersIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
