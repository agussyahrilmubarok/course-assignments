package main

import (
	"context"
	"ecommerce/user-service/handler"
	"ecommerce/user-service/service"
	"ecommerce/user-service/store"
	"ecommerce/user-service/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	mongoDb, err := connectMongo("mongodb://root:secret@localhost:27017", "userdb")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect mongodb")
	}

	jwtSecretKey := "secret"
	jwtDuration := time.Hour * 10

	userMongoStore := store.NewUserMongoStore(mongoDb, log)
	jwtService := service.NewJWTService(jwtSecretKey, jwtDuration)
	authUseCase := usecase.NewAuthUseCase(userMongoStore, jwtService, log)
	authHandler := handler.NewAuthHandler(authUseCase, log)

	authHandler.RegisterRoutes(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info().Msg("Shutting down gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoDb.Client().Disconnect(ctx); err != nil {
			log.Error().Err(err).Msg("failed to disconnect MongoDB")
		}

		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("failed to shutdown Fiber")
		}
	}()

	// start server
	if err := app.Listen(":8081"); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}

func connectMongo(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(dbName), nil
}
