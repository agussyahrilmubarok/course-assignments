package repository

import (
	"context"

	"example.com/backend/internal/domain"
	"gorm.io/gorm"
)

//go:generate mockery --name=ICampaignImageRepository
type ICampaignImageRepository interface {
	FindAll(ctx context.Context) ([]domain.CampaignImage, error)
	FindAllByCampaignID(ctx context.Context, campaignID string) ([]domain.CampaignImage, error)
	FindByID(ctx context.Context, id string) (*domain.CampaignImage, error)
	Create(ctx context.Context, image *domain.CampaignImage) (*domain.CampaignImage, error)
	Update(ctx context.Context, image *domain.CampaignImage) (*domain.CampaignImage, error)
	DeleteByID(ctx context.Context, id string) error
	MarkAllImagesAsNonPrimary(ctx context.Context, campaignID string) error
}

type campaignImageRepository struct {
	db *gorm.DB
}

func NewCampaignImageRepository(db *gorm.DB) ICampaignImageRepository {
	return &campaignImageRepository{db: db}
}

func (r *campaignImageRepository) FindAll(ctx context.Context) ([]domain.CampaignImage, error) {
	var images []domain.CampaignImage
	if err := r.db.WithContext(ctx).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *campaignImageRepository) FindAllByCampaignID(ctx context.Context, campaignID string) ([]domain.CampaignImage, error) {
	var images []domain.CampaignImage
	if err := r.db.WithContext(ctx).
		Where("campaign_id = ?", campaignID).
		Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *campaignImageRepository) FindByID(ctx context.Context, id string) (*domain.CampaignImage, error) {
	var image domain.CampaignImage
	if err := r.db.WithContext(ctx).First(&image, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *campaignImageRepository) Create(ctx context.Context, image *domain.CampaignImage) (*domain.CampaignImage, error) {
	if err := r.db.WithContext(ctx).Create(image).Error; err != nil {
		return nil, err
	}
	return image, nil
}

func (r *campaignImageRepository) Update(ctx context.Context, image *domain.CampaignImage) (*domain.CampaignImage, error) {
	if err := r.db.WithContext(ctx).Save(image).Error; err != nil {
		return nil, err
	}
	return image, nil
}

func (r *campaignImageRepository) DeleteByID(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&domain.CampaignImage{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *campaignImageRepository) MarkAllImagesAsNonPrimary(ctx context.Context, campaignID string) error {
	if err := r.db.WithContext(ctx).Model(&domain.CampaignImage{}).Where("campaign_id = ?", campaignID).
		Update("is_primary", false).Error; err != nil {
		return err
	}

	return nil
}
