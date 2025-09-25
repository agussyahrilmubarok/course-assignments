package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=ITagRepository
type ITagRepository interface {
	FindAll(ctx context.Context) ([]domain.Tag, error)
	FindByID(ctx context.Context, id uint) (*domain.Tag, error)
	Save(ctx context.Context, Tag *domain.Tag) (*domain.Tag, error)
	DeleteByID(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string) bool
}

type tagRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewTagRepository(db *gorm.DB, logger zerolog.Logger) ITagRepository {
	return &tagRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *tagRepository) FindAll(ctx context.Context) ([]domain.Tag, error) {
	var tags []domain.Tag
	if err := r.DB.WithContext(ctx).
		Preload("Products").
		Find(&tags).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to fetch tags")
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) FindByID(ctx context.Context, id uint) (*domain.Tag, error) {
	var tag domain.Tag
	if err := r.DB.WithContext(ctx).
		Preload("Products").
		First(&tag, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("tag_id", id).Msg("tag not found")
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Save(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	if err := r.DB.WithContext(ctx).Save(tag).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to save tag")
		return nil, err
	}
	return tag, nil
}

func (r *tagRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&domain.Tag{}, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("tag_id", id).Msg("failed to delete tag")
		return err
	}
	return nil
}

func (r *tagRepository) ExistsByName(ctx context.Context, name string) bool {
	var count int64
	if err := r.DB.WithContext(ctx).Model(&domain.Tag{}).
		Where("name = ?", name).
		Count(&count).Error; err != nil {
		r.Logger.Error().Err(err).Str("name", name).Msg("failed to check tag existence by name")
		return false
	}
	return count > 0
}
