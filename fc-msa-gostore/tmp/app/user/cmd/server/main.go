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

	"example.com/user/internal/middleware"
	"example.com/user/internal/user"
	"example.com/user/pkg/mongodb"
	"github.com/agussyahrilmubarok/gox/xconfig/xenv"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:8081
// @BasePath /api/v1/users
func main() {
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("MONGO_URI")
		os.Unsetenv("MONGO_DB")
	}()

	env, err := xenv.NewConfig()
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

	jwtMiddleware := middleware.JWTMiddleware(cfg)

	app := fiber.New()

	app.Use(recover.New())

	v1 := app.Group("/api/v1/users")

	v1.Post("/auth/signup", handler.SignUp)
	v1.Post("/auth/signin", handler.SignIn)
	v1.Get("/:id", handler.FindByID)
	v1.Get("/me", jwtMiddleware, handler.FindByMe)

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

func loadEnv() {
	appName := flag.String("app_name", "", "Application name")
	appPort := flag.String("app_port", "", "Application port")
	mongoURI := flag.String("mongo_uri", "", "MongoDB URI")
	mongoDB := flag.String("mongo_db", "", "MongoDB database name")
	flag.Parse()

	os.Setenv("APP_NAME", chooseValue(*appName, "APP_NAME", "USER_SERVICE"))
	os.Setenv("APP_PORT", chooseValue(*appPort, "APP_PORT", "8081"))
	os.Setenv("MONGO_URI", chooseValue(*mongoURI, "MONGO_URI", "mongodb://root:secret@localhost:27017"))
	os.Setenv("MONGO_DB", chooseValue(*mongoDB, "MONGO_DB", "user_db"))
}

func chooseValue(flagVal, envKey, defaultVal string) string {
	if flagVal != "" {
		return flagVal
	}
	if envVal, ok := os.LookupEnv(envKey); ok {
		return envVal
	}
	return defaultVal
}
