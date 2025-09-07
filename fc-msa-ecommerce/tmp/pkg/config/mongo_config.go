package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI    string `json:"uri" mapstructure:"uri"`
	DBName string `json:"dbName" mapstructure:"dbName"`
	client *mongo.Client
}

func (mc *MongoConfig) Connect() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mc.URI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	mc.client = client

	return client.Database(mc.DBName), nil
}

func (mc *MongoConfig) Disconnect(ctx context.Context) error {
	if mc.client != nil {
		if err := mc.client.Disconnect(ctx); err != nil {
			return err
		}
	}
	return nil
}
