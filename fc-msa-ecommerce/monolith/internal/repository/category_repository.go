package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=ICategoryRepository
type ICategoryRepository interface {
	FindAll(ctx context.Context) ([]domain.Category, error)
	FindByID(ctx context.Context, id uint) (*domain.Category, error)
	Save(ctx context.Context, category *domain.Category) (*domain.Category, error)
	DeleteByID(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string) bool
}

type categoryRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewCategoryRepository(db *gorm.DB, logger zerolog.Logger) ICategoryRepository {
	return &categoryRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.DB.WithContext(ctx).
		Preload("Products").
		Find(&categories).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to fetch categories")
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id uint) (*domain.Category, error) {
	var category domain.Category
	if err := r.DB.WithContext(ctx).
		Preload("Products").
		First(&category, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("category_id", id).Msg("category not found")
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Save(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if err := r.DB.WithContext(ctx).Save(category).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to save category")
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&domain.Category{}, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("category_id", id).Msg("failed to delete category")
		return err
	}
	return nil
}

func (r *categoryRepository) ExistsByName(ctx context.Context, name string) bool {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&domain.Category{}).
		Where("name = ?", name).
		Count(&count).Error; err != nil {
		r.Logger.Error().Err(err).Str("name", name).Msg("failed to check category existence by name")
		return false
	}
	return count > 0
}
