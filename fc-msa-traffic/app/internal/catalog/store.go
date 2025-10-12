package catalog

import (
	"context"

	"gorm.io/gorm"
)

//go:generate mockery --name=ICatalogStore
type ICatalogStore interface {
	FindAll(ctx context.Context) ([]Product, error)
	FindByID(ctx context.Context, productID string) (*Product, error)
	FindByName(ctx context.Context, name string) (*Product, error)
	Save(ctx context.Context, product *Product) error
	DeleteByID(ctx context.Context, productID string) error
}

type catalogStore struct {
	db *gorm.DB
}

func NewCatalogStore(db *gorm.DB) ICatalogStore {
	return &catalogStore{db: db}
}

func (s *catalogStore) FindAll(ctx context.Context) ([]Product, error) {
	var products []Product
	if err := s.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (s *catalogStore) FindByID(ctx context.Context, productID string) (*Product, error) {
	var product Product
	if err := s.db.WithContext(ctx).First(&product, "id = ?", productID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *catalogStore) FindByName(ctx context.Context, name string) (*Product, error) {
	var product Product
	if err := s.db.WithContext(ctx).First(&product, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *catalogStore) Save(ctx context.Context, product *Product) error {
	return s.db.WithContext(ctx).Save(product).Error
}

func (s *catalogStore) DeleteByID(ctx context.Context, productID string) error {
	return s.db.WithContext(ctx).Delete(&Product{}, "id = ?", productID).Error
}
