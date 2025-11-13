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

//go:generate mockery --name=IOrderStore
type IOrderStore interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	FindAll(ctx context.Context) ([]Order, error)
	FindAllByUserID(ctx context.Context, userID string) ([]Order, error)
	FindByID(ctx context.Context, orderID string) (*Order, error)
	UpdateStatus(ctx context.Context, orderID string, newStatus OrderStatus) (*Order, error)
	DeleteByID(ctx context.Context, orderID string) error
}

type orderStore struct {
	collection *mongo.Collection
	log        *zap.Logger
}

func NewOrderStore(db *mongo.Database, collectionName string, log *zap.Logger) IOrderStore {
	return &orderStore{
		collection: db.Collection(collectionName),
		log:        log,
	}
}

func (s *orderStore) Create(ctx context.Context, order *Order) (*Order, error) {
	now := time.Now()
	if order.ID == "" {
		return nil, errors.New("order ID is required")
	}
	order.CreatedAt = now
	order.UpdatedAt = now

	_, err := s.collection.InsertOne(ctx, order)
	if err != nil {
		s.log.Error("failed to insert order", zap.Error(err))
		return nil, err
	}
	return order, nil
}

func (s *orderStore) FindAll(ctx context.Context) ([]Order, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error("failed to find all orders", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []Order
	for cursor.Next(ctx) {
		var order Order
		if err := cursor.Decode(&order); err != nil {
			s.log.Error("failed to decode order", zap.Error(err))
			continue
		}
		orders = append(orders, order)
	}
	if err := cursor.Err(); err != nil {
		s.log.Error("cursor error in FindAll", zap.Error(err))
		return nil, err
	}
	return orders, nil
}

func (s *orderStore) FindAllByUserID(ctx context.Context, userID string) ([]Order, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		s.log.Error("failed to find orders by user ID", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []Order
	for cursor.Next(ctx) {
		var order Order
		if err := cursor.Decode(&order); err != nil {
			s.log.Error("failed to decode order", zap.Error(err))
			continue
		}
		orders = append(orders, order)
	}

	if err := cursor.Err(); err != nil {
		s.log.Error("cursor error in FindAllByUserID", zap.Error(err))
		return nil, err
	}

	return orders, nil
}

func (s *orderStore) FindByID(ctx context.Context, orderID string) (*Order, error) {
	var order Order
	err := s.collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to find order by ID", zap.Error(err))
		return nil, err
	}
	return &order, nil
}

func (s *orderStore) UpdateStatus(ctx context.Context, orderID string, newStatus OrderStatus) (*Order, error) {
	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Order
	err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": orderID}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to update order status", zap.Error(err))
		return nil, err
	}
	return &updated, nil
}

func (s *orderStore) DeleteByID(ctx context.Context, orderID string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": orderID})
	if err != nil {
		s.log.Error("failed to delete order", zap.Error(err))
		return err
	}
	return nil
}
