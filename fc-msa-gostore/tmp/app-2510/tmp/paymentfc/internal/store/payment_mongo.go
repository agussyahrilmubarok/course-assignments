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

type IPaymentStore interface {
	Create(ctx context.Context, payment *Payment) (*Payment, error)
	FindAll(ctx context.Context) ([]Payment, error)
	FindAllByStatus(ctx context.Context, status PaymentStatus) ([]Payment, error)
	FindByID(ctx context.Context, paymentID string) (*Payment, error)
	FindByRefCode(ctx context.Context, refCode string) (*Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*Payment, error)
	UpdateByID(ctx context.Context, paymentID string, payment *Payment) (*Payment, error)
	UpdateStatus(ctx context.Context, paymentID string, newStatus PaymentStatus) (*Payment, error)
	DeleteByID(ctx context.Context, paymentID string) error
}

type paymentStore struct {
	collection *mongo.Collection
	log        *zap.Logger
}

func NewPaymentStore(db *mongo.Database, collectionName string, log *zap.Logger) IPaymentStore {
	return &paymentStore{
		collection: db.Collection(collectionName),
		log:        log,
	}
}

func (s *paymentStore) FindAll(ctx context.Context) ([]Payment, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error("failed to find all payments", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []Payment
	for cursor.Next(ctx) {
		var payment Payment
		if err := cursor.Decode(&payment); err != nil {
			s.log.Error("failed to decode payment", zap.Error(err))
			continue
		}
		payments = append(payments, payment)
	}
	if err := cursor.Err(); err != nil {
		s.log.Error("cursor error", zap.Error(err))
		return nil, err
	}
	return payments, nil
}

func (s *paymentStore) FindAllByStatus(ctx context.Context, status PaymentStatus) ([]Payment, error) {
	filter := bson.M{"status": status}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		s.log.Error("failed to find payments by status", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []Payment
	for cursor.Next(ctx) {
		var payment Payment
		if err := cursor.Decode(&payment); err != nil {
			s.log.Error("failed to decode payment", zap.Error(err))
			continue
		}
		payments = append(payments, payment)
	}

	if err := cursor.Err(); err != nil {
		s.log.Error("cursor error while iterating payments by status", zap.Error(err))
		return nil, err
	}

	return payments, nil
}

func (s *paymentStore) Create(ctx context.Context, payment *Payment) (*Payment, error) {
	now := time.Now()
	if payment.ID == "" {
		return nil, errors.New("payment ID is required")
	}
	payment.CreatedAt = now
	payment.UpdatedAt = now

	_, err := s.collection.InsertOne(ctx, payment)
	if err != nil {
		s.log.Error("failed to insert payment", zap.Error(err))
		return nil, err
	}
	return payment, nil
}

func (s *paymentStore) FindByID(ctx context.Context, paymentID string) (*Payment, error) {
	var payment Payment
	err := s.collection.FindOne(ctx, bson.M{"_id": paymentID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to find payment by ID", zap.Error(err))
		return nil, err
	}
	return &payment, nil
}

func (s *paymentStore) FindByRefCode(ctx context.Context, refCode string) (*Payment, error) {
	var payment Payment
	err := s.collection.FindOne(ctx, bson.M{"ref_code": refCode}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to find payment by ref code", zap.Error(err))
		return nil, err
	}
	return &payment, nil
}

func (s *paymentStore) FindByOrderID(ctx context.Context, orderID string) (*Payment, error) {
	var payment Payment
	err := s.collection.FindOne(ctx, bson.M{"order_id": orderID}).Decode(&payment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to find payment by order id", zap.Error(err))
		return nil, err
	}
	return &payment, nil
}

func (s *paymentStore) UpdateByID(ctx context.Context, paymentID string, payment *Payment) (*Payment, error) {
	payment.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"ref_code":    payment.RefCode,
			"user_id":     payment.UserID,
			"order_id":    payment.OrderID,
			"amount":      payment.Amount,
			"status":      payment.Status,
			"invoice_url": payment.InvoiceURL,
			"paid_at":     payment.PaidAt,
			"expired_at":  payment.ExpiredAt,
			"updated_at":  payment.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Payment
	err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": paymentID}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to update payment by ID", zap.Error(err))
		return nil, err
	}

	return &updated, nil
}

func (s *paymentStore) UpdateStatus(ctx context.Context, paymentID string, newStatus PaymentStatus) (*Payment, error) {
	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Payment
	err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": paymentID}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to update payment status", zap.Error(err))
		return nil, err
	}
	return &updated, nil
}

func (s *paymentStore) DeleteByID(ctx context.Context, paymentID string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": paymentID})
	if err != nil {
		s.log.Error("failed to delete payment", zap.Error(err))
		return err
	}
	return nil
}
