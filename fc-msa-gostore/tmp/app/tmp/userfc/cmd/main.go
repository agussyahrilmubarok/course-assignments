package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/userfc/internal/handler"
	"example.com/userfc/internal/service"
	"example.com/userfc/internal/store"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	app_middleware "example.com/pkg/middleware"
	pb "example.com/pkg/proto/v1"
)

var (
	APP_NAME    = "user-service"
	APP_VERSION = "1.0.0"
	APP_PORT    = "8081"
	APP_ENV     = "development"
	GRPC_PORT   = "7071"
)

func main() {
	// Initialize logger
	logger, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}.Build(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	defer logger.Sync()

	// Initialize mongo
	mongoUri := "mongodb://root:secret@localhost:27017"
	if APP_ENV == "production" {
		mongoUri = "mongodb://root:secret@mongo:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		logger.Fatal("failed to connect to mongo", zap.Error(err))
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		logger.Fatal("mongo ping failed", zap.Error(err))
	}

	db := mongoClient.Database("userdb")
	logger.Info("connected to mongo successfully")

	// Initialize services
	userStore := store.NewUserStore(db, "users", logger)
	userService := service.NewUserService(userStore, logger)
	userHandler := handler.NewUserHandler(userService, logger)

	// Start the application (e.g., start HTTP server, etc.)
	r := gin.Default()

	r.Use(app_middleware.RequestLogger(logger))

	v1 := r.Group("/api/v1")
	v1.POST("/sign-up", userHandler.SignUp)
	v1.POST("/sign-in", userHandler.SignIn)
	v1.GET("/users/:id", userHandler.FindUserByID)

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	srv := &http.Server{
		Addr:    ":" + APP_PORT,
		Handler: r,
	}

	go func() {
		logger.Info("starting server", zap.String("port", APP_PORT))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}

	}()

	go func() {
		lis, err := net.Listen("tcp", ":"+GRPC_PORT)
		if err != nil {
			logger.Fatal("failed to listen on grpc port", zap.Error(err))
		}

		grpcServer := grpc.NewServer()
		grpcUserService := service.NewUserGrpcService(userStore, logger)
		pb.RegisterUserServiceServer(grpcServer, grpcUserService)

		logger.Info("grpc server started", zap.String("port", GRPC_PORT))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("grpc server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	if err := mongoClient.Disconnect(shutdownCtx); err != nil {
		logger.Error("error disconnecting mongo", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}
