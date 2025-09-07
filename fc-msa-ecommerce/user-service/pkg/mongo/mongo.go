package mongo

import (
	"context"
	"time"

	"example.com/user-service/pkg/config"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func Connect(cfg *config.Config, log zerolog.Logger) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(cfg.Mongo.DBName)

	return &MongoClient{
		Client: client,
		DB:     db,
	}, nil
}

func (m *MongoClient) Disconnect(ctx context.Context, log zerolog.Logger) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		log.Error().Err(err).Msg("failed to disconnect MongoDB")
		return err
	}

	return nil
}
