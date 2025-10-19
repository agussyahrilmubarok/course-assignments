package repository

import (
	"context"

	"example.com/backend/internal/domain"
	"gorm.io/gorm"
)

//go:generate mockery --name=ICampaignRepository
type ICampaignRepository interface {
	FindAll(ctx context.Context) ([]domain.Campaign, error)
	FindAllWithCampaignImages(ctx context.Context) ([]domain.Campaign, error)
	FindAllByUserID(ctx context.Context, userID string) ([]domain.Campaign, error)
	FindAllWithCampaignImagesByUserID(ctx context.Context, userID string) ([]domain.Campaign, error)
	FindByID(ctx context.Context, id string) (*domain.Campaign, error)
	FindByIDWithCampaignImages(ctx context.Context, id string) (*domain.Campaign, error)
	FindTopCampaigns(ctx context.Context, limit int) ([]domain.Campaign, error)
	Create(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)
	Update(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error)
	DeleteByID(ctx context.Context, id string) error
}

type campaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) ICampaignRepository {
	return &campaignRepository{db: db}
}

func (r *campaignRepository) FindAll(ctx context.Context) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *campaignRepository) FindAllWithCampaignImages(ctx context.Context) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).Preload("CampaignImages").Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *campaignRepository) FindAllByUserID(ctx context.Context, userID string) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *campaignRepository) FindAllWithCampaignImagesByUserID(ctx context.Context, userID string) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).
		Preload("CampaignImages").
		Where("user_id = ?", userID).
		Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *campaignRepository) FindByID(ctx context.Context, id string) (*domain.Campaign, error) {
	var campaign domain.Campaign
	if err := r.db.WithContext(ctx).First(&campaign, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *campaignRepository) FindByIDWithCampaignImages(ctx context.Context, id string) (*domain.Campaign, error) {
	var campaign domain.Campaign
	if err := r.db.WithContext(ctx).
		Preload("CampaignImages").
		First(&campaign, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *campaignRepository) FindTopCampaigns(ctx context.Context, limit int) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).
		Preload("CampaignImages").
		Order("current_amount DESC").
		Limit(limit).
		Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (r *campaignRepository) Create(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	if err := r.db.WithContext(ctx).Create(campaign).Error; err != nil {
		return nil, err
	}
	return campaign, nil
}

func (r *campaignRepository) Update(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	if err := r.db.WithContext(ctx).Save(campaign).Error; err != nil {
		return nil, err
	}
	return campaign, nil
}

func (r *campaignRepository) DeleteByID(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Campaign{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
