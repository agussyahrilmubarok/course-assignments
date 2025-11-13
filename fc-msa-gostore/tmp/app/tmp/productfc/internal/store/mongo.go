package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

//go:generate mockery --name=IProductStore
type IProductStore interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	FindAll(ctx context.Context) ([]Product, error)
	FindByID(ctx context.Context, productID string) (*Product, error)
	UpdateByID(ctx context.Context, productID string, update *Product) (*Product, error)
	DeleteByID(ctx context.Context, productID string) error
}

type productStore struct {
	collection *mongo.Collection
	log        *zap.Logger
}

func NewProductStore(db *mongo.Database, collectionName string, log *zap.Logger) IProductStore {
	return &productStore{
		collection: db.Collection(collectionName),
		log:        log,
	}
}

func (s *productStore) Create(ctx context.Context, product *Product) (*Product, error) {
	now := time.Now()
	if product.ID == "" {
		return nil, errors.New("product ID is required")
	}
	product.CreatedAt = now
	product.UpdatedAt = now

	_, err := s.collection.InsertOne(ctx, product)
	if err != nil {
		s.log.Error("failed to insert product", zap.Error(err))
		return nil, err
	}
	return product, nil
}

func (s *productStore) FindAll(ctx context.Context) ([]Product, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error("failed to find all products", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			s.log.Error("failed to decode product", zap.Error(err))
			continue
		}
		products = append(products, product)
	}
	if err := cursor.Err(); err != nil {
		s.log.Error("cursor error", zap.Error(err))
		return nil, err
	}
	return products, nil
}

func (s *productStore) FindByID(ctx context.Context, productID string) (*Product, error) {
	var product Product
	err := s.collection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to find product by ID", zap.Error(err))
		return nil, err
	}
	return &product, nil
}

func (s *productStore) UpdateByID(ctx context.Context, productID string, update *Product) (*Product, error) {
	update.UpdatedAt = time.Now()

	changes := bson.M{
		"$set": bson.M{
			"name":       update.Name,
			"price":      update.Price,
			"stock":      update.Stock,
			"updated_at": update.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Product
	err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": productID}, changes, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to update product", zap.Error(err))
		return nil, err
	}
	return &updated, nil
}

func (s *productStore) DeleteByID(ctx context.Context, productID string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": productID})
	if err != nil {
		s.log.Error("failed to delete product", zap.Error(err))
		return err
	}
	return nil
}
