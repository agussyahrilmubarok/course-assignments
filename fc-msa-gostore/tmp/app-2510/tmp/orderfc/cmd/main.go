package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/orderfc/internal/client"
	"example.com/orderfc/internal/consumer"
	"example.com/orderfc/internal/handler"
	"example.com/orderfc/internal/producer"
	"example.com/orderfc/internal/service"
	"example.com/orderfc/internal/store"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	app_middleware "example.com/pkg/middleware"
)

var (
	APP_NAME    = "order-service"
	APP_VERSION = "1.0.0"
	APP_PORT    = "8083"
	APP_ENV     = "development"
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

	db := mongoClient.Database("orderdb")
	logger.Info("connected to mongo successfully")

	productGrpcClientHost := "localhost:7072"
	if APP_ENV == "production" {
		productGrpcClientHost = "product-service:7072"
	}
	productGrpcClient, err := client.NewProductGrpcClient(productGrpcClientHost, logger)
	if err != nil {
		logger.Fatal("failed to connect to product grpc server", zap.Error(err))
	}
	defer productGrpcClient.Close()

	kafkaBrokers := []string{"localhost:9092", "localhost:9094", "localhost:9096"}
	if APP_ENV == "production" {
		kafkaBrokers = []string{"kafka1:9092", "kafka2:9094", "kafka3:9096"}
	}
	kafkaProducer := producer.NewKafkaProducer(kafkaBrokers, logger)

	// Initialize services
	orderStore := store.NewOrderStore(db, "orders", logger)
	orderService := service.NewOrderService(orderStore, productGrpcClient, kafkaProducer, logger)
	orderHandler := handler.NewOrderHandler(orderService, logger)

	// Start the application (e.g., start HTTP server, etc.)
	r := gin.Default()

	r.Use(app_middleware.RequestLogger(logger))

	v1 := r.Group("/api/v1")
	v1.POST("/orders/checkout", orderHandler.CheckoutOrder)
	v1.POST("/orders/cancel", orderHandler.CancelOrder)
	v1.GET("/orders/history", orderHandler.GerOrderHistory)

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

	paymentSuccessConsumer := consumer.NewKafkaPaymentSuccessConsumer(kafkaBrokers, orderStore, logger)
	paymentFailedConsumer := consumer.NewKafkaPaymentFailedConsumer(kafkaBrokers, orderStore, kafkaProducer, logger)

	ctxConsumer, cancelConsumer := context.WithCancel(context.Background())

	go paymentSuccessConsumer.Start(ctxConsumer)
	go paymentFailedConsumer.Start(ctxConsumer)

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

	cancelConsumer()
	paymentSuccessConsumer.Stop(shutdownCtx)
	paymentFailedConsumer.Stop(shutdownCtx)

	logger.Info("server exited gracefully")
}
