package service

import (
	"context"

	"example.com/backend/internal/domain"
	"example.com/backend/internal/model"
	"example.com/backend/internal/repository"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=ICampaignService
type ICampaignService interface {
	FindAll(ctx context.Context) ([]model.CampaignDetailDTO, error)
	FindByID(ctx context.Context, id string) (*model.CampaignDetailDTO, error)
	Create(ctx context.Context, campaignDto model.CampaignDTO) error
	Update(ctx context.Context, campaignDto model.CampaignDTO) error
	DeleteByID(ctx context.Context, id string) error
	UploadImage(ctx context.Context, campaignImageDto model.CampaignImageDTO) error
}

type campaignService struct {
	campaignRepo      repository.ICampaignRepository
	campaignImageRepo repository.ICampaignImageRepository
	log               zerolog.Logger
}

func NewCampaignService(
	campaignRepo repository.ICampaignRepository,
	campaignImageRepo repository.ICampaignImageRepository,
	log zerolog.Logger,
) ICampaignService {
	return &campaignService{
		campaignRepo:      campaignRepo,
		campaignImageRepo: campaignImageRepo,
		log:               log,
	}
}

func (s *campaignService) FindAll(ctx context.Context) ([]model.CampaignDetailDTO, error) {
	campaigns, err := s.campaignRepo.FindAllWithCampaignImages(ctx)
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to retrieve campaigns")
		return nil, err
	}

	var campaignDtos []model.CampaignDetailDTO
	for _, campaign := range campaigns {
		var campaignDto model.CampaignDetailDTO
		campaignDto.FromCampaign(&campaign)
		campaignDtos = append(campaignDtos, campaignDto)
	}

	return campaignDtos, nil
}

func (s *campaignService) FindByID(ctx context.Context, id string) (*model.CampaignDetailDTO, error) {
	campaign, err := s.campaignRepo.FindByIDWithCampaignImages(ctx, id)
	if err != nil || campaign == nil {
		s.log.Error().Err(err).Msgf("failed to retrieve campaign id %s", id)
		return nil, err
	}

	var campaignDto model.CampaignDetailDTO
	campaignDto.FromCampaign(campaign)

	return &campaignDto, nil
}

func (s *campaignService) Create(ctx context.Context, campaignDto model.CampaignDTO) error {
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
		s.log.Error().Err(err).Msgf("failed to create new campaign for user %s", campaignDto.UserID)
		return err
	}

	return nil
}

func (s *campaignService) Update(ctx context.Context, campaignDto model.CampaignDTO) error {
	campaignDto.GenerateSlug()

	campaign, err := s.campaignRepo.FindByIDWithCampaignImages(ctx, campaignDto.ID)
	if err != nil || campaign == nil {
		s.log.Error().Err(err).Msgf("failed to retrieve campaign id %s", campaignDto.ID)
		return err
	}

	campaign.Title = campaignDto.Title
	campaign.ShortDescription = campaignDto.ShortDescription
	campaign.GoalAmount = campaignDto.GoalAmount
	campaign.CurrentAmount = campaignDto.CurrentAmount
	campaign.BackerCount = campaignDto.BackerCount
	campaign.Perks = campaignDto.Perks

	_, err = s.campaignRepo.Update(ctx, campaign)
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to update campaign")
		return err
	}

	return nil
}

func (s *campaignService) DeleteByID(ctx context.Context, id string) error {
	if err := s.campaignRepo.DeleteByID(ctx, id); err != nil {
		s.log.Error().Err(err).Msgf("failed to delete campaign id %s", id)
		return err
	}

	return nil
}

func (s *campaignService) UploadImage(ctx context.Context, campaignImageDto model.CampaignImageDTO) error {
	campaignImages, _ := s.campaignImageRepo.FindAllByCampaignID(ctx, campaignImageDto.CampaignID)
	for _, image := range campaignImages {
		image.IsPrimary = false
		_, err := s.campaignImageRepo.Update(ctx, &image)
		s.log.Error().Err(err).Msgf("failed to update campaign image for campaign %s", campaignImageDto.CampaignID)
	}

	campaignImage := &domain.CampaignImage{
		ImageName:  campaignImageDto.ImageName,
		IsPrimary:  true,
		CampaignID: campaignImageDto.CampaignID,
	}

	_, err := s.campaignImageRepo.Create(ctx, campaignImage)
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to save campaign image for campaign %s", campaignImageDto.CampaignID)
		return err
	}

	return nil
}
