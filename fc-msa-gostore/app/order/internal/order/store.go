package order

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockery --name=IStore
type IStore interface {
	FindAll(ctx context.Context) ([]Order, error)
	FindByID(ctx context.Context, orderID string) (*Order, error)
	FindByUserID(ctx context.Context, userID string) ([]Order, error)
	Save(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, orderID string, status OrderStatus) error
	DeleteByID(ctx context.Context, orderID string) error
}

type store struct {
	collection *mongo.Collection
	log        zerolog.Logger
}

func NewStore(db *mongo.Database, log zerolog.Logger) IStore {
	if db == nil {
		log.Fatal().Msg("Database connection is nil")
	}

	return &store{
		collection: db.Collection("orders"),
		log:        log,
	}
}

func (s *store) FindAll(ctx context.Context) ([]Order, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to fetch orders from database")
		return nil, err
	}
	defer cur.Close(ctx)

	var orders []Order
	if err := cur.All(ctx, &orders); err != nil {
		s.log.Error().Err(err).Msg("Failed to decode order list")
		return nil, err
	}

	s.log.Info().Int("count", len(orders)).Msg("Successfully fetched all orders")
	return orders, nil
}

func (s *store) FindByID(ctx context.Context, orderID string) (*Order, error) {
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		s.log.Warn().Str("order_id", orderID).Msg("Invalid order ID format")
		return nil, errors.New("invalid order ID format")
	}

	var order Order
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.log.Warn().Str("order_id", orderID).Msg("Order not found")
			return nil, nil
		}
		s.log.Error().Err(err).Str("order_id", orderID).Msg("Failed to find order by ID")
		return nil, err
	}

	s.log.Info().Str("order_id", orderID).Msg("Order found by ID")
	return &order, nil
}

func (s *store) FindByUserID(ctx context.Context, userID string) ([]Order, error) {
	cur, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		s.log.Error().Err(err).Str("user_id", userID).Msg("Failed to fetch orders by user ID")
		return nil, err
	}
	defer cur.Close(ctx)

	var orders []Order
	if err := cur.All(ctx, &orders); err != nil {
		s.log.Error().Err(err).Str("user_id", userID).Msg("Failed to decode orders for user")
		return nil, err
	}

	s.log.Info().Str("user_id", userID).Int("count", len(orders)).Msg("Successfully fetched orders by user ID")
	return orders, nil
}

func (s *store) Save(ctx context.Context, order *Order) error {
	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	order.UpdatedAt = now

	var objID primitive.ObjectID
	var err error

	if order.ID == "" {
		objID = primitive.NewObjectID()
		order.ID = objID.Hex()
	} else {
		objID, err = primitive.ObjectIDFromHex(order.ID)
		if err != nil {
			s.log.Warn().Str("order_id", order.ID).Msg("Invalid order ID format during save")
			return errors.New("invalid order ID format")
		}
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": order}
	opts := options.Update().SetUpsert(true)

	_, err = s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		s.log.Error().Err(err).Str("order_id", order.ID).Msg("Failed to save order")
		return err
	}

	s.log.Info().Str("order_id", order.ID).Msg("Order saved successfully")
	return nil
}

func (s *store) UpdateStatus(ctx context.Context, orderID string, status OrderStatus) error {
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		s.log.Warn().Str("order_id", orderID).Msg("Invalid order ID format during update status")
		return errors.New("invalid order ID format")
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	res, err := s.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		s.log.Error().Err(err).Str("order_id", orderID).Msg("Failed to update order status")
		return err
	}

	if res.MatchedCount == 0 {
		s.log.Warn().Str("order_id", orderID).Msg("Order not found for status update")
		return nil
	}

	s.log.Info().Str("order_id", orderID).Str("status", string(status)).Msg("Order status updated successfully")
	return nil
}

func (s *store) DeleteByID(ctx context.Context, orderID string) error {
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		s.log.Warn().Str("order_id", orderID).Msg("Invalid order ID format during delete")
		return errors.New("invalid order ID format")
	}

	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		s.log.Error().Err(err).Str("order_id", orderID).Msg("Failed to delete order")
		return err
	}

	if res.DeletedCount == 0 {
		s.log.Warn().Str("order_id", orderID).Msg("No order deleted (order not found)")
		return nil
	}

	s.log.Info().Str("order_id", orderID).Msg("Order deleted successfully")
	return nil
}
