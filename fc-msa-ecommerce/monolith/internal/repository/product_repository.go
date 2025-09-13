package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IProductRepository
type IProductRepository interface {
	FindAll(ctx context.Context) ([]domain.Product, error)
	FindAllByCategoryID(ctx context.Context, categoryId uint) ([]domain.Product, error)
	FindByID(ctx context.Context, id uint) (*domain.Product, error)
	Save(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteByID(ctx context.Context, id uint) error
}

type productRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewProductRepository(db *gorm.DB, logger zerolog.Logger) IProductRepository {
	return &productRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *productRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.DB.WithContext(ctx).
		Preload("Category").
		Preload("Tags").
		Find(&products).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to fetch products")
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindAllByCategoryID(ctx context.Context, categoryId uint) ([]domain.Product, error) {
	var products []domain.Product
	if err := r.DB.WithContext(ctx).
		Where("category_id = ?", categoryId).
		Preload("Category").
		Preload("Tags").
		Find(&products).Error; err != nil {
		r.Logger.Error().Err(err).
			Uint("category_id", categoryId).
			Msg("failed to fetch products by category")
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.DB.WithContext(ctx).
		Preload("Category").
		Preload("Tags").
		First(&product, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("product_id", id).Msg("product not found")
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Save(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if err := r.DB.WithContext(ctx).Save(product).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to save product")
		return nil, err
	}
	return product, nil
}

func (r *productRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&domain.Product{}, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("product_id", id).Msg("failed to delete product")
		return err
	}
	return nil
}
