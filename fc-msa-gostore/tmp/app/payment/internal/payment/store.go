package payment

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
	FindAll(ctx context.Context) ([]Payment, error)
	FindByID(ctx context.Context, paymentID string) (*Payment, error)
	FindByOrderID(ctx context.Context, orderID string) ([]Payment, error)
	Save(ctx context.Context, payment *Payment) error
	UpdateStatus(ctx context.Context, paymentID string, status PaymentStatus) error
	DeleteByID(ctx context.Context, paymentID string) error
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
		collection: db.Collection("payments"),
		log:        log,
	}
}

func (s *store) FindAll(ctx context.Context) ([]Payment, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to fetch payments from database")
		return nil, err
	}
	defer cur.Close(ctx)

	var payments []Payment
	if err := cur.All(ctx, &payments); err != nil {
		s.log.Error().Err(err).Msg("Failed to decode payment list")
		return nil, err
	}

	s.log.Info().Int("count", len(payments)).Msg("Successfully fetched all payments")
	return payments, nil
}

func (s *store) FindByID(ctx context.Context, paymentID string) (*Payment, error) {
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		s.log.Warn().Str("payment_id", paymentID).Msg("Invalid payment ID format")
		return nil, errors.New("invalid payment ID format")
	}

	var payment Payment
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&payment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.log.Warn().Str("payment_id", paymentID).Msg("Payment not found")
			return nil, nil
		}
		s.log.Error().Err(err).Str("payment_id", paymentID).Msg("Failed to find payment by ID")
		return nil, err
	}

	s.log.Info().Str("payment_id", paymentID).Msg("Payment found by ID")
	return &payment, nil
}

func (s *store) FindByOrderID(ctx context.Context, orderID string) ([]Payment, error) {
	cur, err := s.collection.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		s.log.Error().Err(err).Str("order_id", orderID).Msg("Failed to fetch payments by order ID")
		return nil, err
	}
	defer cur.Close(ctx)

	var payments []Payment
	if err := cur.All(ctx, &payments); err != nil {
		s.log.Error().Err(err).Str("order_id", orderID).Msg("Failed to decode payments for order")
		return nil, err
	}

	s.log.Info().Str("order_id", orderID).Int("count", len(payments)).Msg("Successfully fetched payments by order ID")
	return payments, nil
}

func (s *store) Save(ctx context.Context, payment *Payment) error {
	now := time.Now()
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	payment.UpdatedAt = now

	var objID primitive.ObjectID
	var err error

	if payment.ID == "" {
		objID = primitive.NewObjectID()
		payment.ID = objID.Hex()
	} else {
		objID, err = primitive.ObjectIDFromHex(payment.ID)
		if err != nil {
			s.log.Warn().Str("payment_id", payment.ID).Msg("Invalid payment ID format during save")
			return errors.New("invalid payment ID format")
		}
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": payment}
	opts := options.Update().SetUpsert(true)

	_, err = s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		s.log.Error().Err(err).Str("payment_id", payment.ID).Msg("Failed to save payment")
		return err
	}

	s.log.Info().Str("payment_id", payment.ID).Msg("Payment saved successfully")
	return nil
}

func (s *store) UpdateStatus(ctx context.Context, paymentID string, status PaymentStatus) error {
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		s.log.Warn().Str("payment_id", paymentID).Msg("Invalid payment ID format during update status")
		return errors.New("invalid payment ID format")
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	res, err := s.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		s.log.Error().Err(err).Str("payment_id", paymentID).Msg("Failed to update payment status")
		return err
	}

	if res.MatchedCount == 0 {
		s.log.Warn().Str("payment_id", paymentID).Msg("Payment not found for status update")
		return nil
	}

	s.log.Info().Str("payment_id", paymentID).Str("status", string(status)).Msg("Payment status updated successfully")
	return nil
}

func (s *store) DeleteByID(ctx context.Context, paymentID string) error {
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		s.log.Warn().Str("payment_id", paymentID).Msg("Invalid payment ID format during delete")
		return errors.New("invalid payment ID format")
	}

	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		s.log.Error().Err(err).Str("payment_id", paymentID).Msg("Failed to delete payment")
		return err
	}

	if res.DeletedCount == 0 {
		s.log.Warn().Str("payment_id", paymentID).Msg("No payment deleted (payment not found)")
		return nil
	}

	s.log.Info().Str("payment_id", paymentID).Msg("Payment deleted successfully")
	return nil
}
