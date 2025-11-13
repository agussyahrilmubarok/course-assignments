package mongodb

import (
	"context"
	"sync"
	"time"

	"example.com/user/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFactory struct {
	mu      sync.Mutex
	clients map[string]*mongo.Client
	cfg     *config.Config
}

func NewMongoFactory(cfg *config.Config) *MongoFactory {
	return &MongoFactory{
		cfg:     cfg,
		clients: make(map[string]*mongo.Client),
	}
}

func (f *MongoFactory) GetClient(uri string) (*mongo.Client, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if client, ok := f.clients[uri]; ok {
		return client, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	f.clients[uri] = client
	return client, nil
}

func (f *MongoFactory) GetDatabase(dbName string) (*mongo.Database, error) {
	client, err := f.GetClient(f.cfg.MongoDB.URI)
	if err != nil {
		return nil, err
	}

	return client.Database(dbName), nil
}

func (f *MongoFactory) CloseAll() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for uri, client := range f.clients {
		_ = client.Disconnect(context.Background())
		delete(f.clients, uri)
	}

	return nil
}
