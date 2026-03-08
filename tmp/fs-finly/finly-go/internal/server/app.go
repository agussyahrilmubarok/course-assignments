package server

import (
	"context"
	"finly-go/internal/sessions"
	"finly-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type App struct {
	fiberApp    *fiber.App
	mongoClient *mongo.Client
}

func NewApp() (*App, error) {
	client, err := ConnectMongoDB()
	if err != nil {
		logger.Log.Error("Error connecting to MongoDB", zap.Error(err))
		return nil, err
	}

	app := &App{
		mongoClient: client,
		fiberApp: fiber.New(fiber.Config{
			Views:       html.New("./views", ".html"),
			ViewsLayout: "layouts/main",
		}),
	}

	app.fiberApp.Static("/", "./assets")

	sessions.AddSession()

	WebRoutes(app.fiberApp)

	return app, nil
}

func (a *App) GetApp() *fiber.App {
	return a.fiberApp
}

func (a *App) Close() {
	if a.mongoClient != nil {
		logger.Log.Info("Disconnecting from MongoDB...")

		if err := a.mongoClient.Disconnect(context.Background()); err != nil {
			logger.Log.Error("Error during MongoDB disconnection", zap.Error(err))
		} else {
			logger.Log.Info("Successfully disconnected from MongoDB")
		}
	}
}
