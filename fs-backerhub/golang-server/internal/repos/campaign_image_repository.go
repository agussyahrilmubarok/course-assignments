package repos

import (
	"context"

	"example.com.backend/internal/domain"
	"example.com.backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
	log := logger.GetLoggerFromContext(ctx)

	var images []domain.CampaignImage
	if err := r.db.WithContext(ctx).Find(&images).Error; err != nil {
		log.Error("failed fetching all campaign images", zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched all campaign images", zap.Int("count", len(images)))
	return images, nil
}

func (r *campaignImageRepository) FindAllByCampaignID(ctx context.Context, campaignID string) ([]domain.CampaignImage, error) {
	log := logger.GetLoggerFromContext(ctx)

	var images []domain.CampaignImage
	if err := r.db.WithContext(ctx).
		Where("campaign_id = ?", campaignID).
		Find(&images).Error; err != nil {

		log.Error("failed fetching campaign images by campaign id",
			zap.String("campaign_id", campaignID),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("successfully fetched campaign images by campaign id",
		zap.String("campaign_id", campaignID),
		zap.Int("count", len(images)),
	)
	return images, nil
}

func (r *campaignImageRepository) FindByID(ctx context.Context, id string) (*domain.CampaignImage, error) {
	log := logger.GetLoggerFromContext(ctx)

	var image domain.CampaignImage
	if err := r.db.WithContext(ctx).First(&image, "id = ?", id).Error; err != nil {
		log.Error("failed fetching campaign image by id", zap.String("campaign_image_id", id), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched campaign image by id", zap.String("campaign_image_id", id))
	return &image, nil
}

func (r *campaignImageRepository) Create(ctx context.Context, image *domain.CampaignImage) (*domain.CampaignImage, error) {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Create(image).Error; err != nil {
		log.Error("failed creating campaign image",
			zap.String("campaign_id", image.CampaignID),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("successfully created campaign image",
		zap.String("campaign_image_id", image.ID),
		zap.String("campaign_id", image.CampaignID),
	)
	return image, nil
}

func (r *campaignImageRepository) Update(ctx context.Context, image *domain.CampaignImage) (*domain.CampaignImage, error) {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Save(image).Error; err != nil {
		log.Error("failed updating campaign image",
			zap.String("id", image.ID),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("successfully updated campaign image", zap.String("campaign_image_id", image.ID))
	return image, nil
}

func (r *campaignImageRepository) DeleteByID(ctx context.Context, id string) error {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Delete(&domain.CampaignImage{}, "id = ?", id).Error; err != nil {
		log.Error("failed deleting campaign image", zap.String("campaign_image_id", id), zap.Error(err))
		return err
	}

	log.Info("successfully deleted campaign image", zap.String("campaign_image_id", id))
	return nil
}

func (r *campaignImageRepository) MarkAllImagesAsNonPrimary(ctx context.Context, campaignID string) error {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).
		Model(&domain.CampaignImage{}).
		Where("campaign_id = ?", campaignID).
		Update("is_primary", false).Error; err != nil {

		log.Error("failed marking all images as non-primary",
			zap.String("campaign_id", campaignID),
			zap.Error(err),
		)
		return err
	}

	log.Info("successfully marked all images as non-primary",
		zap.String("campaign_id", campaignID),
	)
	return nil
}
