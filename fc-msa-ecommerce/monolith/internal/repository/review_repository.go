package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IReviewRepository
type IReviewRepository interface {
	FindAll(ctx context.Context) ([]domain.Review, error)
	FindByID(ctx context.Context, id uint) (*domain.Review, error)
	FindByProductID(ctx context.Context, productID uint) ([]domain.Review, error)
	Save(ctx context.Context, review *domain.Review) (*domain.Review, error)
	Update(ctx context.Context, review *domain.Review) (*domain.Review, error)
	DeleteByID(ctx context.Context, id uint) error
}

type reviewRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewReviewRepository(db *gorm.DB, logger zerolog.Logger) IReviewRepository {
	return &reviewRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *reviewRepository) FindAll(ctx context.Context) ([]domain.Review, error) {
	var reviews []domain.Review
	if err := r.DB.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Find(&reviews).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to fetch reviews")
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) FindByID(ctx context.Context, id uint) (*domain.Review, error) {
	var review domain.Review
	if err := r.DB.WithContext(ctx).
		Preload("User").
		Preload("Product").
		First(&review, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("review_id", id).Msg("review not found")
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) FindByProductID(ctx context.Context, productID uint) ([]domain.Review, error) {
	var reviews []domain.Review
	if err := r.DB.WithContext(ctx).
		Preload("User").
		Where("product_id = ?", productID).
		Find(&reviews).Error; err != nil {
		r.Logger.Error().Err(err).Uint("product_id", productID).Msg("failed to fetch reviews by product")
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) Save(ctx context.Context, review *domain.Review) (*domain.Review, error) {
	if err := r.DB.WithContext(ctx).Create(review).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to save review")
		return nil, err
	}
	return review, nil
}

func (r *reviewRepository) Update(ctx context.Context, review *domain.Review) (*domain.Review, error) {
	if err := r.DB.WithContext(ctx).Save(review).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to update review")
		return nil, err
	}
	return review, nil
}

func (r *reviewRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&domain.Review{}, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("review_id", id).Msg("failed to delete review")
		return err
	}
	return nil
}
