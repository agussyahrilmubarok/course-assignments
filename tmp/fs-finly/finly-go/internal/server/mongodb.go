package server

import (
	"context"
	"finly-go/configs"
	"finly-go/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var mongoDb *mongo.Database

func ConnectMongoDB() (*mongo.Client, error) {
	mongoURI := configs.GetMongoURI()
	mongoDB := configs.GetMongoDBName()
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Log.Error("Failed to connect to MongoDB", zap.String("uri", mongoURI), zap.Error(err))
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Log.Error("Failed to ping MongoDB", zap.Error(err))
		return nil, err
	}

	mongoDb = client.Database(mongoDB)

	logger.Log.Info("Successfully connected to MongoDB", zap.String("database", mongoDB))
	return client, nil
}

type Collections struct {
	UserCollection     *mongo.Collection
	CustomerCollection *mongo.Collection
	InvoiceCollection  *mongo.Collection
}

func GetCollections() *Collections {
	return &Collections{
		UserCollection:     mongoDb.Collection("users"),
		CustomerCollection: mongoDb.Collection("customers"),
		InvoiceCollection:  mongoDb.Collection("invoices"),
	}
}
