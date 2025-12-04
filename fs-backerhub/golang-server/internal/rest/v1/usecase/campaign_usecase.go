package usecaseV1

import (
	"context"
	"fmt"
	"strings"
	"time"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/repos"
	"example.com.backend/internal/service"
	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/logger"
	"github.com/gosimple/slug"
	"go.uber.org/zap"

	payloadV1 "example.com.backend/internal/rest/v1/payload"
)

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
	campaignRepo      repos.ICampaignRepository
	campaignImageRepo repos.ICampaignImageRepository
	uploadService     service.IUploadService
}

func NewCampaignUseCaseV1(
	campaignRepo repos.ICampaignRepository,
	campaignImageRepo repos.ICampaignImageRepository,
	uploadService service.IUploadService,
) ICampaignUseCaseV1 {
	return &campaignUseCaseV1{
		campaignRepo:      campaignRepo,
		campaignImageRepo: campaignImageRepo,
		uploadService:     uploadService,
	}
}

func (uc *campaignUseCaseV1) FindAll(ctx context.Context) ([]payloadV1.CampaignResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaigns, err := uc.campaignRepo.FindAllWithCampaignImages(ctx)
	if err != nil {
		log.Error("failed fetching campaigns", zap.Error(err))
		return nil, exception.NewBadRequest("Campaigns are not found", err)
	}

	var resps []payloadV1.CampaignResponse
	for _, campaign := range campaigns {
		var resp payloadV1.CampaignResponse
		resp.FromCampaign(&campaign)
		resps = append(resps, resp)
	}

	log.Info("successfully fetched campaigns", zap.Int("count", len(resps)))
	return resps, nil
}

func (uc *campaignUseCaseV1) FindByID(ctx context.Context, campaignID string) (*payloadV1.CampaignDetailResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaign, err := uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaignID)
	if err != nil {
		log.Error("failed fetching campaign by id", zap.String("campaign_id", campaignID), zap.Error(err))
		return nil, exception.NewInternal("Failed to get campaign", err)
	}
	if campaign == nil {
		log.Warn("campaign not found by id", zap.String("campaign_id", campaignID))
		return nil, exception.NewNotFound("Campaign is not found", nil)
	}

	var resp payloadV1.CampaignDetailResponse
	resp.FromCampaign(campaign)
	log.Info("successfully fetched campaign by id", zap.String("campaign_id", campaignID))
	return &resp, nil
}

func (uc *campaignUseCaseV1) FindAllTop(ctx context.Context) ([]payloadV1.CampaignResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaigns, err := uc.campaignRepo.FindTopCampaigns(ctx, 6)
	if err != nil {
		log.Error("failed fetching top campaigns", zap.Error(err))
		return nil, exception.NewBadRequest("Top campaigns not found", err)
	}

	resps := make([]payloadV1.CampaignResponse, 0, len(campaigns))
	for i := range campaigns {
		var resp payloadV1.CampaignResponse
		resp.FromCampaign(&campaigns[i])
		resps = append(resps, resp)
	}

	log.Info("successfully fetched top campaigns", zap.Int("count", len(resps)))
	return resps, nil
}

func (uc *campaignUseCaseV1) FindAllByUser(ctx context.Context, userID string) ([]payloadV1.CampaignResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaigns, err := uc.campaignRepo.FindAllWithCampaignImagesByUserID(ctx, userID)
	if err != nil {
		log.Error("failed fetching campaigns by user", zap.String("user_id", userID), zap.Error(err))
		return nil, exception.NewBadRequest("Campaigns are not found", err)
	}

	var resps []payloadV1.CampaignResponse
	for _, campaign := range campaigns {
		var resp payloadV1.CampaignResponse
		resp.FromCampaign(&campaign)
		resps = append(resps, resp)
	}

	log.Info("successfully fetched campaigns by user", zap.String("user_id", userID), zap.Int("count", len(resps)))
	return resps, nil
}

func (uc *campaignUseCaseV1) FindByIDUser(ctx context.Context, campaignID, userID string) (*payloadV1.CampaignDetailResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaign, err := uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaignID)
	if err != nil {
		log.Error("failed fetching campaign by id", zap.String("campaign_id", campaignID), zap.Error(err))
		return nil, exception.NewInternal("Failed to get campaign", err)
	}
	if campaign == nil {
		log.Warn("campaign not found by id", zap.String("campaign_id", campaignID))
		return nil, exception.NewNotFound("Campaign is not found", nil)
	}

	if campaign.UserID != userID {
		log.Warn("unauthorized access to campaign", zap.String("campaign_id", campaignID), zap.String("user_id", userID))
		return nil, exception.NewUnauthorized("You do not have permission", nil)
	}

	var resp payloadV1.CampaignDetailResponse
	resp.FromCampaign(campaign)
	log.Info("successfully fetched campaign by user", zap.String("campaign_id", campaignID), zap.String("user_id", userID))
	return &resp, nil
}

func (uc *campaignUseCaseV1) CreateByUser(ctx context.Context, param payloadV1.CampaignRequest, userID string) (*payloadV1.CampaignResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

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
		log.Error("failed creating campaign", zap.String("user_id", userID), zap.Error(err))
		return nil, exception.NewInternal("Create new campaign failed", err)
	}

	campaign, err = uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaign.ID)
	if err != nil || campaign == nil {
		log.Error("failed fetching newly created campaign", zap.String("campaign_id", campaign.ID), zap.Error(err))
		return nil, exception.NewInternal("Create new campaign failed", err)
	}

	var campaignResp payloadV1.CampaignResponse
	campaignResp.FromCampaign(campaign)
	log.Info("successfully created campaign", zap.String("campaign_id", campaign.ID), zap.String("user_id", userID))
	return &campaignResp, nil
}

func (uc *campaignUseCaseV1) UpdateByIDByUser(ctx context.Context, param payloadV1.CampaignRequest, campaignID, userID string) (*payloadV1.CampaignResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaign, err := uc.campaignRepo.FindByID(ctx, campaignID)
	if err != nil {
		log.Error("failed fetching campaign for update", zap.String("campaign_id", campaignID), zap.Error(err))
		return nil, exception.NewInternal("Update campaign failed", err)
	}
	if campaign == nil {
		log.Warn("campaign not found for update", zap.String("campaign_id", campaignID))
		return nil, exception.NewNotFound("Campaign is not found", nil)
	}

	if campaign.UserID != userID {
		log.Warn("unauthorized update attempt", zap.String("campaign_id", campaignID), zap.String("user_id", userID))
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
		log.Error("failed updating campaign", zap.String("campaign_id", campaignID), zap.Error(err))
		return nil, exception.NewInternal("Update campaign failed", err)
	}

	campaign, err = uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaign.ID)
	if err != nil || campaign == nil {
		log.Error("failed fetching updated campaign", zap.String("campaign_id", campaign.ID), zap.Error(err))
		return nil, exception.NewInternal("Update campaign failed", err)
	}

	var campaignResp payloadV1.CampaignResponse
	campaignResp.FromCampaign(campaign)
	log.Info("successfully updated campaign", zap.String("campaign_id", campaign.ID), zap.String("user_id", userID))
	return &campaignResp, nil
}

func (uc *campaignUseCaseV1) UploadImageByUser(ctx context.Context, param payloadV1.CampaignImageRequest) error {
	log := logger.GetLoggerFromContext(ctx)

	campaign, err := uc.campaignRepo.FindByID(ctx, param.CampaignID)
	if err != nil {
		log.Error("failed fetching campaign for image upload", zap.String("campaign_id", param.CampaignID), zap.Error(err))
		return exception.NewInternal("Failed to get campaign", err)
	}
	if campaign == nil {
		log.Warn("campaign not found for image upload", zap.String("campaign_id", param.CampaignID))
		return exception.NewNotFound("Campaign is not found", nil)
	}

	if campaign.UserID != param.UserID {
		log.Warn("unauthorized image upload attempt", zap.String("campaign_id", param.CampaignID), zap.String("user_id", param.UserID))
		return exception.NewUnauthorized("You do not have permission", nil)
	}

	campaignUploadPath := "public/uploads/campaigns"
	imagePath, err := uc.uploadService.SaveLocal(ctx, campaignUploadPath, param.CampaignImage, param.CampaignID)
	if err != nil {
		log.Error("failed uploading campaign image", zap.String("campaign_id", param.CampaignID), zap.Error(err))
		return exception.NewInternal("Upload campaign image fail", err)
	}

	if param.IsPrimary {
		if err := uc.campaignImageRepo.MarkAllImagesAsNonPrimary(ctx, param.CampaignID); err != nil {
			log.Error("failed marking campaign images as non-primary", zap.String("campaign_id", param.CampaignID), zap.Error(err))
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
		uc.uploadService.RemoveLocal(ctx, fmt.Sprintf("%v/", campaignUploadPath), trimmedPath)
		log.Error("failed saving campaign image to DB", zap.String("campaign_id", param.CampaignID), zap.Error(err))
		return exception.NewInternal("Upload campaign image fail", err)
	}

	log.Info("successfully uploaded campaign image", zap.String("campaign_id", param.CampaignID), zap.String("image_name", trimmedPath))
	return nil
}

func (uc *campaignUseCaseV1) DeleteByIDByUser(ctx context.Context, campaignID, userID string) error {
	log := logger.GetLoggerFromContext(ctx)

	campaign, err := uc.campaignRepo.FindByIDWithCampaignImages(ctx, campaignID)
	if err != nil {
		log.Error("failed fetching campaign for deletion", zap.String("campaign_id", campaignID), zap.Error(err))
		return exception.NewInternal("Failed to get campaign", err)
	}
	if campaign == nil {
		log.Warn("campaign not found for deletion", zap.String("campaign_id", campaignID))
		return exception.NewNotFound("Campaign is not found", nil)
	}

	if campaign.UserID != userID {
		log.Warn("unauthorized delete attempt", zap.String("campaign_id", campaignID), zap.String("user_id", userID))
		return exception.NewUnauthorized("You do not have permission", nil)
	}

	campaignUploadPath := "public/uploads/campaigns"
	for _, campaignImage := range campaign.CampaignImages {
		if err := uc.uploadService.RemoveLocal(ctx, fmt.Sprintf("%v/", campaignUploadPath), campaignImage.ImageName); err != nil {
			log.Error("failed removing campaign image", zap.String("campaign_image_id", campaignImage.ID), zap.Error(err))
			continue
		}
	}

	if err = uc.campaignRepo.DeleteByID(ctx, campaignID); err != nil {
		log.Error("failed deleting campaign", zap.String("campaign_id", campaignID), zap.Error(err))
		return exception.NewInternal("Delete campaign fail", err)
	}

	log.Info("successfully deleted campaign", zap.String("campaign_id", campaignID), zap.String("user_id", userID))
	return nil
}
