package repos

import (
	"context"

	"example.com.backend/internal/domain"
	"example.com.backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
	log := logger.GetLoggerFromContext(ctx)

	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).Find(&campaigns).Error; err != nil {
		log.Error("failed fetching all campaigns", zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched all campaigns", zap.Int("count", len(campaigns)))
	return campaigns, nil
}

func (r *campaignRepository) FindAllWithCampaignImages(ctx context.Context) ([]domain.Campaign, error) {
	log := logger.GetLoggerFromContext(ctx)

	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).Preload("CampaignImages").Find(&campaigns).Error; err != nil {
		log.Error("failed fetching campaigns with images", zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched campaigns with images", zap.Int("count", len(campaigns)))
	return campaigns, nil
}

func (r *campaignRepository) FindAllByUserID(ctx context.Context, userID string) ([]domain.Campaign, error) {
	log := logger.GetLoggerFromContext(ctx)

	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&campaigns).Error; err != nil {
		log.Error("failed fetching campaigns by user id", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched campaigns by user id",
		zap.String("user_id", userID),
		zap.Int("count", len(campaigns)),
	)
	return campaigns, nil
}

func (r *campaignRepository) FindAllWithCampaignImagesByUserID(ctx context.Context, userID string) ([]domain.Campaign, error) {
	log := logger.GetLoggerFromContext(ctx)

	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).
		Preload("CampaignImages").
		Where("user_id = ?", userID).
		Find(&campaigns).Error; err != nil {
		log.Error("failed fetching campaigns with images by user id", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched campaigns with images by user id",
		zap.String("user_id", userID),
		zap.Int("count", len(campaigns)),
	)
	return campaigns, nil
}

func (r *campaignRepository) FindByID(ctx context.Context, id string) (*domain.Campaign, error) {
	log := logger.GetLoggerFromContext(ctx)

	var campaign domain.Campaign
	if err := r.db.WithContext(ctx).First(&campaign, "id = ?", id).Error; err != nil {
		log.Error("failed fetching campaign by id", zap.String("campaign_id", id), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched campaign by id", zap.String("campaign_id", id))
	return &campaign, nil
}

func (r *campaignRepository) FindByIDWithCampaignImages(ctx context.Context, id string) (*domain.Campaign, error) {
	log := logger.GetLoggerFromContext(ctx)

	var campaign domain.Campaign
	if err := r.db.WithContext(ctx).
		Preload("CampaignImages").
		First(&campaign, "id = ?", id).Error; err != nil {
		log.Error("failed fetching campaign with images by id", zap.String("campaign_id", id), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched campaign with images by id", zap.String("campaign_id", id))
	return &campaign, nil
}

func (r *campaignRepository) FindTopCampaigns(ctx context.Context, limit int) ([]domain.Campaign, error) {
	log := logger.GetLoggerFromContext(ctx)

	var campaigns []domain.Campaign
	if err := r.db.WithContext(ctx).
		Preload("CampaignImages").
		Order("current_amount DESC").
		Limit(limit).
		Find(&campaigns).Error; err != nil {
		log.Error("failed fetching top campaigns", zap.Int("limit", limit), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched top campaigns",
		zap.Int("limit", limit),
		zap.Int("count", len(campaigns)),
	)
	return campaigns, nil
}

func (r *campaignRepository) Create(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Create(campaign).Error; err != nil {
		log.Error("failed creating campaign", zap.String("user_id", campaign.UserID), zap.Error(err))
		return nil, err
	}

	log.Info("successfully created campaign", zap.String("campaign_id", campaign.ID), zap.String("user_id", campaign.UserID))
	return campaign, nil
}

func (r *campaignRepository) Update(ctx context.Context, campaign *domain.Campaign) (*domain.Campaign, error) {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Save(campaign).Error; err != nil {
		log.Error("failed updating campaign", zap.String("campaign_id", campaign.ID), zap.Error(err))
		return nil, err
	}

	log.Info("successfully updated campaign", zap.String("campaign_id", campaign.ID))
	return campaign, nil
}

func (r *campaignRepository) DeleteByID(ctx context.Context, id string) error {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Delete(&domain.Campaign{}, "id = ?", id).Error; err != nil {
		log.Error("failed deleting campaign", zap.String("campaign_id", id), zap.Error(err))
		return err
	}

	log.Info("successfully deleted campaign", zap.String("campaign_id", id))
	return nil
}
