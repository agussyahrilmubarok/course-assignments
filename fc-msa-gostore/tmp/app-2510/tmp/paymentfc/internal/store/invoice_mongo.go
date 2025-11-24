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

type IInvoiceStore interface {
	Create(ctx context.Context, invoice *Invoice) (*Invoice, error)
	FindAllByUserID(ctx context.Context) ([]Invoice, error)
	FindByID(ctx context.Context, invoiceID string) (*Invoice, error)
	UpdateStatus(ctx context.Context, invoiceID string, newStatus InvoiceStatus) (*Invoice, error)
	DeleteByID(ctx context.Context, invoiceID string) error
}

type invoiceStore struct {
	collection *mongo.Collection
	log        *zap.Logger
}

func NewInvoiceStore(db *mongo.Database, collectionName string, log *zap.Logger) IInvoiceStore {
	return &invoiceStore{
		collection: db.Collection(collectionName),
		log:        log,
	}
}

func (s *invoiceStore) Create(ctx context.Context, invoice *Invoice) (*Invoice, error) {
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()
	_, err := s.collection.InsertOne(ctx, invoice)
	if err != nil {
		s.log.Error("failed to insert invoice", zap.Error(err))
		return nil, err
	}
	return invoice, nil
}

func (s *invoiceStore) FindAllByUserID(ctx context.Context) ([]Invoice, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, errors.New("user_id not found in context")
	}

	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		s.log.Error("failed to find invoices by user ID", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var invoices []Invoice
	if err := cursor.All(ctx, &invoices); err != nil {
		s.log.Error("failed to decode invoice list", zap.Error(err))
		return nil, err
	}

	return invoices, nil
}

func (s *invoiceStore) FindByID(ctx context.Context, invoiceID string) (*Invoice, error) {
	filter := bson.M{"id": invoiceID}

	var invoice Invoice
	err := s.collection.FindOne(ctx, filter).Decode(&invoice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to find invoice by ID", zap.Error(err))
		return nil, err
	}

	return &invoice, nil
}

func (s *invoiceStore) UpdateStatus(ctx context.Context, invoiceID string, newStatus InvoiceStatus) (*Invoice, error) {
	filter := bson.M{"id": invoiceID}
	update := bson.M{
		"$set": bson.M{
			"status":  newStatus,
			"updated": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Invoice
	err := s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		s.log.Error("failed to update invoice status", zap.Error(err))
		return nil, err
	}

	return &updated, nil
}

func (s *invoiceStore) DeleteByID(ctx context.Context, invoiceID string) error {
	filter := bson.M{"id": invoiceID}

	res, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		s.log.Error("failed to delete invoice", zap.Error(err))
		return err
	}

	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
