package service

import (
	"context"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	"example.com.backend/internal/repos"
	"example.com.backend/pkg/logger"
	"go.uber.org/zap"
)

type ICampaignService interface {
	FindAll(ctx context.Context) ([]model.CampaignDetailDTO, error)
	FindByID(ctx context.Context, id string) (*model.CampaignDetailDTO, error)
	Create(ctx context.Context, campaignDto model.CampaignDTO) error
	Update(ctx context.Context, campaignDto model.CampaignDTO) error
	DeleteByID(ctx context.Context, id string) error
	UploadImage(ctx context.Context, campaignImageDto model.CampaignImageDTO) error
}

type campaignService struct {
	campaignRepo      repos.ICampaignRepository
	campaignImageRepo repos.ICampaignImageRepository
}

func NewCampaignService(
	campaignRepo repos.ICampaignRepository,
	campaignImageRepo repos.ICampaignImageRepository,
) ICampaignService {
	return &campaignService{
		campaignRepo:      campaignRepo,
		campaignImageRepo: campaignImageRepo,
	}
}

func (s *campaignService) FindAll(ctx context.Context) ([]model.CampaignDetailDTO, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaigns, err := s.campaignRepo.FindAllWithCampaignImages(ctx)
	if err != nil {
		log.Error("failed retrieving all campaigns", zap.Error(err))
		return nil, err
	}

	var campaignDtos []model.CampaignDetailDTO
	for _, campaign := range campaigns {
		var dto model.CampaignDetailDTO
		dto.FromCampaign(&campaign)
		campaignDtos = append(campaignDtos, dto)
	}

	log.Info("successfully retrieved campaigns", zap.Int("count", len(campaigns)))
	return campaignDtos, nil
}

func (s *campaignService) FindByID(ctx context.Context, id string) (*model.CampaignDetailDTO, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaign, err := s.campaignRepo.FindByIDWithCampaignImages(ctx, id)
	if err != nil {
		log.Error("failed retrieving campaign by id", zap.String("campaign_id", id), zap.Error(err))
		return nil, err
	}
	if campaign == nil {
		log.Warn("campaign not found", zap.String("campaign_id", id))
		return nil, nil
	}

	var dto model.CampaignDetailDTO
	dto.FromCampaign(campaign)

	log.Info("successfully retrieved campaign", zap.String("campaign_id", id))
	return &dto, nil
}

func (s *campaignService) Create(ctx context.Context, campaignDto model.CampaignDTO) error {
	log := logger.GetLoggerFromContext(ctx)

	campaignDto.GenerateSlug()

	campaign := &domain.Campaign{
		Title:            campaignDto.Title,
		ShortDescription: campaignDto.ShortDescription,
		Description:      campaignDto.Description,
		GoalAmount:       campaignDto.GoalAmount,
		CurrentAmount:    0,
		BackerCount:      0,
		Perks:            campaignDto.Perks,
		Slug:             campaignDto.Slug,
		UserID:           campaignDto.UserID,
	}

	_, err := s.campaignRepo.Create(ctx, campaign)
	if err != nil {
		log.Error("failed creating campaign",
			zap.String("user_id", campaignDto.UserID),
			zap.Error(err),
		)
		return err
	}

	log.Info("successfully created campaign",
		zap.String("campaign_id", campaign.ID),
		zap.String("user_id", campaign.UserID),
	)
	return nil
}

func (s *campaignService) Update(ctx context.Context, campaignDto model.CampaignDTO) error {
	log := logger.GetLoggerFromContext(ctx)

	campaignDto.GenerateSlug()

	campaign, err := s.campaignRepo.FindByIDWithCampaignImages(ctx, campaignDto.ID)
	if err != nil {
		log.Error("failed retrieving campaign for update",
			zap.String("campaign_id", campaignDto.ID),
			zap.Error(err),
		)
		return err
	}
	if campaign == nil {
		log.Warn("campaign not found for update", zap.String("campaign_id", campaignDto.ID))
		return nil
	}

	campaign.Title = campaignDto.Title
	campaign.ShortDescription = campaignDto.ShortDescription
	campaign.GoalAmount = campaignDto.GoalAmount
	campaign.CurrentAmount = campaignDto.CurrentAmount
	campaign.BackerCount = campaignDto.BackerCount
	campaign.Perks = campaignDto.Perks

	_, err = s.campaignRepo.Update(ctx, campaign)
	if err != nil {
		log.Error("failed updating campaign", zap.String("campaign_id", campaignDto.ID), zap.Error(err))
		return err
	}

	log.Info("successfully updated campaign", zap.String("campaign_id", campaignDto.ID))
	return nil
}

func (s *campaignService) DeleteByID(ctx context.Context, id string) error {
	log := logger.GetLoggerFromContext(ctx)

	err := s.campaignRepo.DeleteByID(ctx, id)
	if err != nil {
		log.Error("failed deleting campaign", zap.String("campaign_id", id), zap.Error(err))
		return err
	}

	log.Info("successfully deleted campaign", zap.String("campaign_id", id))
	return nil
}

func (s *campaignService) UploadImage(ctx context.Context, campaignImageDto model.CampaignImageDTO) error {
	log := logger.GetLoggerFromContext(ctx)

	// mark all existing images as non-primary
	err := s.campaignImageRepo.MarkAllImagesAsNonPrimary(ctx, campaignImageDto.CampaignID)
	if err != nil {
		log.Error("failed marking images as non-primary",
			zap.String("campaign_id", campaignImageDto.CampaignID),
			zap.Error(err),
		)
		return err
	}

	campaignImage := &domain.CampaignImage{
		ImageName:  campaignImageDto.ImageName,
		IsPrimary:  true,
		CampaignID: campaignImageDto.CampaignID,
	}

	_, err = s.campaignImageRepo.Create(ctx, campaignImage)
	if err != nil {
		log.Error("failed saving campaign image",
			zap.String("campaign_id", campaignImageDto.CampaignID),
			zap.Error(err),
		)
		return err
	}

	log.Info("successfully uploaded campaign image",
		zap.String("campaign_id", campaignImageDto.CampaignID),
		zap.String("campaign_image_name", campaignImageDto.ImageName),
	)
	return nil
}
