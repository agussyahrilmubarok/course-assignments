package usecaseV1

import (
	"context"
	"fmt"
	"strings"
	"time"

	payloadV1 "example.com/backend/api/v1/payload"
	"example.com/backend/internal/domain"
	"example.com/backend/internal/exception"
	"example.com/backend/internal/repository"
	"example.com/backend/internal/service"
	"github.com/gosimple/slug"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=ICampaignUseCaseV1
type ICampaignUseCaseV1 interface {
	FindAll(ctx context.Context) ([]payloadV1.CampaignResponse, error)
	FindAllTop(ctx context.Context) ([]payloadV1.CampaignResponse, error)
	FindByID(ctx context.Context, campaignID string) (*payloadV1.CampaignDetailResponse, error)
	FindAllByUser(ctx context.Context, userID string) ([]payloadV1.CampaignResponse, error)
	FindByIDUser(ctx context.Context, campaignID, userID string) (*payloadV1.CampaignDetailResponse, error)
	CreateByUser(ctx context.Context, param payloadV1.CampaignRequest, userID string) (*payloadV1.CampaignResponse, error)
	UpdateByIDByUser(ctx context.Context, param payloadV1.CampaignRequest, campaignID string, userID string) (*payloadV1.CampaignResponse, error)
	UploadImageByUser(ctx context.Context, param payloadV1.CampaignImageRequest) error
	DeleteByIDByUser(ctx context.Context, campaignID, userID string) error
}

type campaignUseCaseV1 struct {
	campaignRepo      repository.ICampaignRepository
	campaignImageRepo repository.ICampaignImageRepository
	uploadService     service.IUploadService
	log               zerolog.Logger
}

func NewCampaignUseCaseV1(
	campaignRepo repository.ICampaignRepository,
	campaignImageRepo repository.ICampaignImageRepository,
	uploadService service.IUploadService,
	log zerolog.Logger,
) ICampaignUseCaseV1 {
	return &campaignUseCaseV1{
		campaignRepo:      campaignRepo,
		campaignImageRepo: campaignImageRepo,
		uploadService:     uploadService,
		log:               log,
	}
}

func (uc *campaignUseCaseV1) FindAll(ctx context.Context) ([]payloadV1.CampaignResponse, error) {
	campaigns, err := uc.campaignRepo.FindAllWithCampaignImages(ctx)
	if err != nil {
		uc.log.Warn().Msg("Campaigns are not found")
		return nil, exception.NewBadRequest("Campaigns are not found", err)
	}

	var resps []payloadV1.CampaignResponse
	for _, campaign := range campaigns {
		var resp payloadV1.CampaignResponse
		resp.FromCampaign(&campaign)
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *campaignUseCaseV1) FindByID(ctx context.Context, campaignID string) (*payloadV1.CampaignDetailResponse, error) {
	campaign, err := uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaignID)
	if err != nil || campaign == nil {
		uc.log.Warn().Msg("campaign id not found")
		return nil, exception.NewNotFound("Campaign is not found", err)
	}

	var resp payloadV1.CampaignDetailResponse
	resp.FromCampaign(campaign)

	return &resp, nil
}

func (uc *campaignUseCaseV1) FindAllTop(ctx context.Context) ([]payloadV1.CampaignResponse, error) {
	campaigns, err := uc.campaignRepo.FindTopCampaigns(ctx, 6)
	if err != nil {
		uc.log.Warn().Err(err).Msg("failed to fetch top campaigns")
		return nil, exception.NewBadRequest("campaigns are not found", err)
	}

	resps := make([]payloadV1.CampaignResponse, 0, len(campaigns))

	for i := range campaigns {
		var resp payloadV1.CampaignResponse
		resp.FromCampaign(&campaigns[i])
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *campaignUseCaseV1) FindAllByUser(ctx context.Context, userID string) ([]payloadV1.CampaignResponse, error) {
	campaigns, err := uc.campaignRepo.FindAllWithCampaignImagesByUserID(ctx, userID)
	if err != nil {
		uc.log.Warn().Msg("Campaigns are not found")
		return nil, exception.NewBadRequest("Campaigns are not found", err)
	}

	var resps []payloadV1.CampaignResponse
	for _, campaign := range campaigns {
		var resp payloadV1.CampaignResponse
		resp.FromCampaign(&campaign)
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *campaignUseCaseV1) FindByIDUser(ctx context.Context, campaignID string, userID string) (*payloadV1.CampaignDetailResponse, error) {
	campaign, err := uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaignID)
	if err != nil || campaign == nil {
		uc.log.Warn().Msg("campaign id not found")
		return nil, exception.NewNotFound("Campaign is not found", err)
	}

	if campaign.UserID != userID {
		uc.log.Warn().Msg("campaign is unauthorized")
		return nil, exception.NewUnauthorized("You do not have permission", nil)
	}

	var resp payloadV1.CampaignDetailResponse
	resp.FromCampaign(campaign)

	return &resp, nil
}

func (uc *campaignUseCaseV1) CreateByUser(ctx context.Context, param payloadV1.CampaignRequest, userID string) (*payloadV1.CampaignResponse, error) {
	campaign := &domain.Campaign{
		Title:            param.Title,
		ShortDescription: param.ShortDescription,
		Description:      param.Description,
		GoalAmount:       param.GoalAmount,
		CurrentAmount:    0,
		BackerCount:      0,
		Perks:            param.Perks,
		Slug:             slug.Make(fmt.Sprintf("%s-%s-%d", param.Title, userID, time.Now().Unix())),
		UserID:           userID,
	}

	campaign, err := uc.campaignRepo.Create(ctx, campaign)
	if err != nil || campaign == nil {
		uc.log.Error().Err(err).Msg("failed when creating campaign")
		return nil, exception.NewInternal("Create new campaign failed", err)
	}

	campaign, err = uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaign.ID)
	if err != nil || campaign == nil {
		uc.log.Error().Err(err).Msg("failed when creating campaign")
		return nil, exception.NewInternal("Create new campaign failed", err)
	}

	var campaignResp payloadV1.CampaignResponse
	campaignResp.FromCampaign(campaign)

	return &campaignResp, nil
}

func (uc *campaignUseCaseV1) UpdateByIDByUser(ctx context.Context, param payloadV1.CampaignRequest, campaignID string, userID string) (*payloadV1.CampaignResponse, error) {
	campaign, err := uc.campaignRepo.FindByID(ctx, campaignID)
	if err != nil || campaign == nil {
		uc.log.Warn().Msg("campaign id not found")
		return nil, exception.NewNotFound("Campaign is not found", err)
	}

	if campaign.UserID != userID {
		uc.log.Warn().Msg("campaign is unauthorized")
		return nil, exception.NewUnauthorized("You do not have permission", nil)
	}

	campaign.Title = param.Title
	campaign.ShortDescription = param.ShortDescription
	campaign.Description = param.Description
	campaign.GoalAmount = param.GoalAmount
	campaign.Perks = param.Perks
	campaign.Slug = slug.Make(fmt.Sprintf("%s-%s-%d", param.Title, userID, time.Now().Unix()))

	campaign, err = uc.campaignRepo.Update(ctx, campaign)
	if err != nil || campaign == nil {
		uc.log.Error().Err(err).Msg("failed when updating campaign")
		return nil, exception.NewInternal("Update campaign failed", err)
	}

	campaign, err = uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaign.ID)
	if err != nil || campaign == nil {
		uc.log.Error().Err(err).Msg("failed when creating campaign")
		return nil, exception.NewInternal("Create new campaign failed", err)
	}

	var campaignResp payloadV1.CampaignResponse
	campaignResp.FromCampaign(campaign)

	return &campaignResp, nil
}

func (uc *campaignUseCaseV1) UploadImageByUser(ctx context.Context, param payloadV1.CampaignImageRequest) error {
	campaign, err := uc.campaignRepo.FindByID(ctx, param.CampaignID)
	if err != nil || campaign == nil {
		uc.log.Warn().Msg("campaign id not found")
		return exception.NewNotFound("Campaign is not found", err)
	}

	if campaign.UserID != param.UserID {
		uc.log.Warn().Msg("campaign is unauthorized")
		return exception.NewUnauthorized("You do not have permission", nil)
	}

	campaignUploadPath := "public/uploads/campaigns"
	imagePath, err := uc.uploadService.SaveLocal(campaignUploadPath, param.CampaignImage, param.CampaignID)
	if err != nil {
		uc.log.Error().Err(err).Msg("failed when uploading image campaign")
		return exception.NewInternal("Upload campaign image fail", err)
	}

	if param.IsPrimary {
		if err := uc.campaignImageRepo.MarkAllImagesAsNonPrimary(ctx, param.CampaignID); err != nil {
			uc.log.Error().Err(err).Msg("failed to set campaign images as non primary")
		}
	}

	trimmedPath := strings.TrimPrefix(imagePath, fmt.Sprintf("%v/", campaignUploadPath))
	campaignImage := &domain.CampaignImage{
		ImageName:  trimmedPath,
		IsPrimary:  param.IsPrimary,
		CampaignID: param.CampaignID,
	}
	_, err = uc.campaignImageRepo.Create(ctx, campaignImage)
	if err != nil {
		uc.uploadService.RemoveLocal(fmt.Sprintf("%v/", campaignUploadPath), trimmedPath)
		uc.log.Error().Err(err).Msg("failed when saving campaign image")
		return exception.NewInternal("Upload campaign image fail", err)
	}

	return nil
}

func (uc *campaignUseCaseV1) DeleteByIDByUser(ctx context.Context, campaignID string, userID string) error {
	campaign, err := uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaignID)
	if err != nil || campaign == nil {
		uc.log.Warn().Msg("campaign id not found")
		return exception.NewNotFound("Campaign is not found", err)
	}

	if campaign.UserID != userID {
		uc.log.Warn().Msg("campaign is unauthorized")
		return exception.NewUnauthorized("You do not have permission", nil)
	}

	campaignUploadPath := "public/uploads/campaigns"
	for _, campaignImage := range campaign.CampaignImages {
		if err := uc.uploadService.RemoveLocal(fmt.Sprintf("%v/", campaignUploadPath), campaignImage.ImageName); err != nil {
			uc.log.Error().Err(err).Msgf("failed to remove campaign image %s", campaignImage.ID)
			continue
		}
	}

	if err = uc.campaignRepo.DeleteByID(ctx, campaignID); err != nil {
		uc.log.Error().Err(err).Msg("failed when deleting campaign")
		return exception.NewInternal("Delete campaign fail", err)
	}

	return nil
}
