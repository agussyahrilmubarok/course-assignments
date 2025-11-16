package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/paymentfc/internal/client"
	"example.com/paymentfc/internal/consumer"
	"example.com/paymentfc/internal/handler"
	"example.com/paymentfc/internal/producer"
	"example.com/paymentfc/internal/scheduler"
	"example.com/paymentfc/internal/service"
	"example.com/paymentfc/internal/store"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	app_middleware "example.com/pkg/middleware"
)

var (
	APP_NAME    = "payment-service"
	APP_VERSION = "1.0.0"
	APP_PORT    = "8084"
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

	db := mongoClient.Database("paymentdb")
	logger.Info("connected to mongo successfully")

	userGrpcClientHost := "localhost:7071"
	if APP_ENV == "production" {
		userGrpcClientHost = "user-service:7071"
	}
	userGrpcClient, err := client.NewUserGrpcClient(userGrpcClientHost, logger)
	if err != nil {
		logger.Fatal("failed to connect to user grpc server", zap.Error(err))
	}
	defer userGrpcClient.Close()

	brokers := []string{"localhost:9092", "localhost:9094", "localhost:9096"}
	if APP_ENV == "production" {
		brokers = []string{"kafka1:9092", "kafka2:9094", "kafka3:9096"}
	}
	kafkaProducer := producer.NewKafkaProducer(brokers, logger)

	// Initialize services
	paymentStore := store.NewPaymentStore(db, "payments", logger)
	invoiceStore := store.NewInvoiceStore(db, "invoices", logger)
	paymentService := service.NewPaymentService("xnd_development_3D9HnCKDo5K4UiCugrU8SjbKy5Xx1JYoOs7iBTTt5rl3CX08p8I3TEbGQlB9K7", paymentStore, invoiceStore, userGrpcClient, *kafkaProducer, logger)
	paymentHandler := handler.NewPaymentHandler("WEBHOOK_TOKEN", paymentService, logger)

	// Initialize scheduler
	paymentSyncScheduler := scheduler.NewPaymentScheduler(paymentStore, logger)
	paymentSyncScheduler.Start()

	// Start the application (e.g., start HTTP server, etc.)
	r := gin.Default()

	r.Use(app_middleware.RequestLogger(logger))

	v1 := r.Group("/api/v1")
	v1.POST("/payments/xendit/webhook", paymentHandler.XenditWebhook)
	v1.POST("/payments/invoice/:order_id/pdf ", paymentHandler.DownloadInvoice)

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

	// Initialize Kafka consumer
	kafkaBrokers := []string{"localhost:9092", "localhost:9094", "localhost:9096"}
	if APP_ENV == "production" {
		kafkaBrokers = []string{"kafka1:9092", "kafka2:9094", "kafka3:9096"}
	}
	kafkaInitPaymentConsumer := consumer.NewKafkaInitPaymentConsumer(kafkaBrokers, paymentService, logger)

	ctxConsumer, cancelConsumer := context.WithCancel(context.Background())

	go kafkaInitPaymentConsumer.Start(ctxConsumer)

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
	kafkaInitPaymentConsumer.Stop(shutdownCtx)

	paymentSyncScheduler.Stop()

	logger.Info("server exited gracefully")
}
