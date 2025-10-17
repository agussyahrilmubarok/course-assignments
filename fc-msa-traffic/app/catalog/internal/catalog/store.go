package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IStore
type IStore interface {
	FindProducts(ctx context.Context) ([]Product, error)
	FindProductByID(ctx context.Context, productID string) (*Product, error)
	FindProductByName(ctx context.Context, name string) (*Product, error)
	SaveProduct(ctx context.Context, product *Product) error
	DeleteProductByID(ctx context.Context, productID string) error
	ReverseProductByID(ctx context.Context, productID string, quantity int) error
	ReleaseProductByID(ctx context.Context, productID string, quantity int) error
}

var (
	PRODUCT_STOCK_KEY = "product:stock:%s"
	PRODUCT_STOCK_TTL = 60 * time.Minute
	PRODUCT_KEY       = "product:price:%s"
	PRODUCT_TTL       = 60 * time.Minute
)

type store struct {
	db  *gorm.DB
	rdb *redis.Client
	log zerolog.Logger
}

func NewStore(
	db *gorm.DB,
	rdb *redis.Client,
	logger zerolog.Logger,
) IStore {
	return &store{
		db:  db,
		rdb: rdb,
		log: logger,
	}
}

func (s *store) FindProducts(ctx context.Context) ([]Product, error) {
	var products []Product
	if err := s.db.WithContext(ctx).Find(&products).Error; err != nil {
		s.log.Error().Err(err).Msg("Failed to find all products")
		return nil, err
	}

	s.log.Info().Int("products_count", len(products)).Msg("Find all products successfully")
	return products, nil
}

func (s *store) FindProductByID(ctx context.Context, productID string) (*Product, error) {
	var product Product
	if err := s.db.WithContext(ctx).First(&product, "id = ?", productID).Error; err != nil {
		s.log.Error().Err(err).Str("product_id", productID).Msg("Failed to find product by ID")
		return nil, err
	}

	s.log.Info().Str("product_id", productID).Msg("Find product by id successfully")
	return &product, nil
}

func (s *store) FindProductByName(ctx context.Context, name string) (*Product, error) {
	var product Product
	if err := s.db.WithContext(ctx).First(&product, "name = ?", name).Error; err != nil {
		s.log.Error().Err(err).Str("product_name", name).Msg("Failed to find product by name")
		return nil, err
	}

	s.log.Info().Str("product_name", name).Msg("Find product by name successfully")
	return &product, nil
}

func (s *store) SaveProduct(ctx context.Context, product *Product) error {
	if err := s.db.WithContext(ctx).Save(product).Error; err != nil {
		s.log.Error().Err(err).Str("product_name", product.Name).Msg("Failed to save product")
		return err
	}

	s.log.Info().Str("product_id", product.ID).Msg("Save product successfully")
	return nil
}

func (s *store) DeleteProductByID(ctx context.Context, productID string) error {
	if err := s.db.WithContext(ctx).Delete(&Product{}, "id = ?", productID).Error; err != nil {
		s.log.Error().Err(err).Str("product_id", productID).Msg("Failed to delete product")
		return err
	}

	s.log.Info().Str("product_id", productID).Msg("Delete product successfully")
	return nil
}

func (s *store) ReverseProductByID(ctx context.Context, productID string, quantity int) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product Product

		if err := tx.First(&product, "id = ?", productID).Error; err != nil {
			s.log.Error().Err(err).Str("product_id", productID).Msg("Failed to find product for stock reversal")
			return err
		}

		if product.Stock < quantity {
			s.log.Warn().Int("current_stock", product.Stock).Int("requested_quantity", quantity).Str("product_id", productID).Msg("insufficient stock")
			return fmt.Errorf("insufficient stock for product ID %s", productID)
		}

		product.Stock -= quantity
		if err := tx.Save(&product).Error; err != nil {
			s.log.Error().Err(err).Str("product_id", productID).Msg("Failed to update product stock after reversal")
			return err
		}

		key := fmt.Sprintf(PRODUCT_STOCK_KEY, productID)
		err := s.rdb.Set(ctx, key, product.Stock, PRODUCT_STOCK_TTL).Err()
		if err != nil {
			s.log.Warn().Err(err).Str("product_id", productID).Msg("Failed to update Redis stock cache")
		}

		s.log.Info().Str("product_id", productID).Int("new_stock", product.Stock).Int("decreased_by", quantity).Msg("Successfully reversed product stock")
		return nil
	})

}

func (s *store) ReleaseProductByID(ctx context.Context, productID string, quantity int) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product Product

		if err := tx.First(&product, "id = ?", productID).Error; err != nil {
			s.log.Error().Err(err).Str("product_id", productID).Msg("Failed to find product for stock release")
			return err
		}

		product.Stock += quantity
		if err := tx.Save(&product).Error; err != nil {
			s.log.Error().Err(err).Str("product_id", productID).Msg("Failed to update product stock after release")
			return err
		}

		key := fmt.Sprintf(PRODUCT_STOCK_KEY, productID)
		err := s.rdb.Set(ctx, key, product.Stock, PRODUCT_STOCK_TTL).Err()
		if err != nil {
			s.log.Warn().Err(err).Str("product_id", productID).Msg("Failed to update Redis stock cache")
		}

		s.log.Info().Str("product_id", productID).Int("new_stock", product.Stock).Int("increased_by", quantity).Msg("Successfully released product stock")
		return nil
	})
}
