package product

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
	FindAll(ctx context.Context) ([]Product, error)
	FindByID(ctx context.Context, productID string) (*Product, error)
	Save(ctx context.Context, product *Product) error
	DeleteByID(ctx context.Context, productID string) error
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
		collection: db.Collection("products"),
		log:        log,
	}
}

func (s *store) FindAll(ctx context.Context) ([]Product, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to fetch products from database")
		return nil, err
	}
	defer cur.Close(ctx)

	var products []Product
	if err := cur.All(ctx, &products); err != nil {
		s.log.Error().Err(err).Msg("Failed to decode product list")
		return nil, err
	}

	s.log.Info().Int("count", len(products)).Msg("Successfully fetched all products")
	return products, nil
}

func (s *store) FindByID(ctx context.Context, productID string) (*Product, error) {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		s.log.Warn().Str("product_id", productID).Msg("Invalid product ID format")
		return nil, errors.New("invalid product ID format")
	}

	var product Product
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.log.Warn().Str("product_id", productID).Msg("Product not found")
			return nil, nil
		}
		s.log.Error().Err(err).Str("product_id", productID).Msg("Failed to find product by ID")
		return nil, err
	}

	s.log.Info().Str("product_id", productID).Msg("Product found by ID")
	return &product, nil
}

func (s *store) Save(ctx context.Context, product *Product) error {
	now := time.Now()
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}
	product.UpdatedAt = now

	var objID primitive.ObjectID
	var err error

	if product.ID == "" {
		objID = primitive.NewObjectID()
		product.ID = objID.Hex()
	} else {
		objID, err = primitive.ObjectIDFromHex(product.ID)
		if err != nil {
			s.log.Warn().Str("product_id", product.ID).Msg("Invalid product ID format during save")
			return errors.New("invalid product ID format")
		}
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": product}
	opts := options.Update().SetUpsert(true)

	_, err = s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		s.log.Error().Err(err).Str("product_id", product.ID).Msg("Failed to save product")
		return err
	}

	s.log.Info().Str("product_id", product.ID).Msg("Product saved successfully")
	return nil
}

func (s *store) DeleteByID(ctx context.Context, productID string) error {
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		s.log.Warn().Str("product_id", productID).Msg("Invalid product ID format during delete")
		return errors.New("invalid product ID format")
	}

	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		s.log.Error().Err(err).Str("product_id", productID).Msg("Failed to delete product")
		return err
	}

	if res.DeletedCount == 0 {
		s.log.Warn().Str("product_id", productID).Msg("No product deleted (product not found)")
		return nil
	}

	s.log.Info().Str("product_id", productID).Msg("Product deleted successfully")
	return nil
}
