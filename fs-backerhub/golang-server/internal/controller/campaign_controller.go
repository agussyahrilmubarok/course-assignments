package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	"example.com.backend/internal/service"
	"example.com.backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	CAMPAIGN_UPLOAD_PATH = "public/uploads/campaigns"
)

type campaignController struct {
	baseController
	campaignService service.ICampaignService
	userService     service.IUserService
	uploadService   service.IUploadService
}

func NewCampaignController(
	campaignService service.ICampaignService,
	userService service.IUserService,
	uploadService service.IUploadService,
) *campaignController {
	return &campaignController{
		campaignService: campaignService,
		userService:     userService,
		uploadService:   uploadService,
	}
}

func (h *campaignController) Index(c *gin.Context) {
	data := gin.H{"title": "Campaigns"}

	ctx := c.Request.Context()

	h.showAllCampaigns(c, ctx, data)
}

func (h *campaignController) Add(c *gin.Context) {
	data := gin.H{"title": "New Campaign"}

	ctx := c.Request.Context()

	h.showNewCampaignWithData(c, ctx, data)
}

func (h *campaignController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "New Campaign"}

	var input struct {
		Title            string `form:"title" binding:"required"`
		ShortDescription string `form:"short_description" binding:"required"`
		Description      string `form:"description" binding:"required"`
		GoalAmount       int    `form:"goal_amount" binding:"required,gt=0"`
		Perks            string `form:"perks" binding:"required"`
		UserID           string `form:"user_id" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		log.Warn("failed to bind create campaign input", zap.Error(err))
		data["form"] = input
		data["error"] = "Invalid input."
		h.showNewCampaignWithData(c, ctx, data)
		return
	}

	campaignDto := model.CampaignDTO{
		Title:            input.Title,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       float64(input.GoalAmount),
		Perks:            input.Perks,
		UserID:           input.UserID,
	}

	if err := h.campaignService.Create(ctx, campaignDto); err != nil {
		log.Error("failed to create campaign", zap.Error(err), zap.String("title", input.Title))
		data["form"] = input
		data["error"] = "Failed to create campaign."
		h.showNewCampaignWithData(c, ctx, data)
		return
	}

	data["title"] = "Campaign"
	data["form"] = nil
	data["success"] = fmt.Sprintf("New campaign created successfully: %s", input.Title)

	h.showAllCampaigns(c, ctx, data)
}

func (h *campaignController) Show(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Campaign"}

	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		log.Error("campaign not found", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	data["campaign"] = campaignDto

	h.renderHTML(c, http.StatusOK, "campaign_show.html", data)
}

func (h *campaignController) Edit(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Edit Campaign"}

	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		log.Error("campaign not found", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	data["campaign"] = campaignDto

	h.renderHTML(c, http.StatusOK, "campaign_edit.html", data)
}

func (h *campaignController) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Edit Campaign"}

	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		log.Error("campaign not found", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	var input struct {
		Title            string `form:"title" binding:"required"`
		ShortDescription string `form:"short_description" binding:"required"`
		Description      string `form:"description" binding:"required"`
		GoalAmount       int    `form:"goal_amount" binding:"required"`
		CurrentAmount    int    `form:"current_amount"`
		BackerCount      int    `form:"backer_count"`
		Perks            string `form:"perks" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		log.Warn("invalid bind update campaign input", zap.Error(err))
		data["campaign"] = campaignDto
		data["error"] = "Invalid input."
		h.renderHTML(c, http.StatusBadRequest, "campaign_edit.html", data)
		return
	}

	campaign := model.CampaignDTO{
		ID:               campaignDto.ID,
		Title:            input.Title,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       float64(input.GoalAmount),
		CurrentAmount:    float64(input.CurrentAmount),
		BackerCount:      int64(input.BackerCount),
		Perks:            input.Perks,
	}

	if err := h.campaignService.Update(ctx, campaign); err != nil {
		log.Error("failed to update campaign", zap.Error(err), zap.String("campaign_id", campaign.ID))
		data["campaign"] = input
		data["error"] = "Failed to update campaign."
		h.renderHTML(c, http.StatusInternalServerError, "campaign_edit.html", data)
		return
	}

	data["title"] = "Campaigns"
	data["success"] = fmt.Sprintf("Campaign updated successfully: %s", input.Title)

	h.showAllCampaigns(c, ctx, data)
}

func (h *campaignController) Image(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Upload Campaign Image"}

	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		log.Error("campaign not found for image upload", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	data["campaign"] = campaignDto

	h.renderHTML(c, http.StatusOK, "campaign_edit_image.html", data)
}

func (h *campaignController) Upload(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Upload Campaign Image"}

	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		log.Error("campaign not found for image upload", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		log.Error("failed to retrieve campaign image", zap.Error(err))
		data["error"] = "Failed to get image file."
		h.renderHTML(c, http.StatusBadRequest, "campaign_edit_image.html", data)
		return
	}

	imagePath, err := h.uploadService.SaveLocal(ctx, CAMPAIGN_UPLOAD_PATH, imageFile, idStr)
	if err != nil {
		log.Error("failed to upload campaign image", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Failed to upload image."
		h.renderHTML(c, http.StatusInternalServerError, "campaign_edit_image.html", data)
		return
	}

	trimmedPath := strings.TrimPrefix(imagePath, fmt.Sprintf("%v/", CAMPAIGN_UPLOAD_PATH))

	if err := h.campaignService.UploadImage(ctx, model.CampaignImageDTO{
		ImageName:  trimmedPath,
		CampaignID: campaignDto.ID,
	}); err != nil {
		h.uploadService.RemoveLocal(ctx, fmt.Sprintf("%v/", CAMPAIGN_UPLOAD_PATH), trimmedPath)
		log.Error("failed to save campaign image", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Failed to save campaign image."
		h.renderHTML(c, http.StatusInternalServerError, "campaign_edit_image.html", data)
		return
	}

	data["title"] = "Campaigns"
	data["success"] = "Campaign image uploaded successfully"

	h.showAllCampaigns(c, ctx, data)
}

func (h *campaignController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Delete Campaign"}

	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		log.Error("campaign not found for deletion", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	for _, img := range campaignDto.CampaignImages {
		if img.ImageName != "default.jpeg" {
			if err := h.uploadService.RemoveLocal(ctx, fmt.Sprintf("%v/", CAMPAIGN_UPLOAD_PATH), img.ImageName); err != nil {
				log.Warn("failed to remove campaign image", zap.Error(err), zap.String("image", img.ImageName))
			}
		}
	}

	if err := h.campaignService.DeleteByID(ctx, idStr); err != nil {
		log.Error("failed to delete campaign", zap.Error(err), zap.String("campaign_id", idStr))
		data["error"] = "Failed to delete campaign."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	data["title"] = "Campaigns"
	data["success"] = "Campaign has been deleted successfully"

	h.showAllCampaigns(c, ctx, data)
}

func (h *campaignController) showAllCampaigns(c *gin.Context, ctx context.Context, data gin.H) {
	campaignDtos, err := h.campaignService.FindAll(ctx)
	if err != nil {
		data["campaigns"] = []model.CampaignDetailDTO{}
		h.renderHTML(c, http.StatusBadRequest, "campaign_index.html", data)
		return
	}

	data["campaigns"] = campaignDtos

	h.renderHTML(c, http.StatusOK, "campaign_index.html", data)
}

func (h *campaignController) showNewCampaignWithData(c *gin.Context, ctx context.Context, data gin.H) {
	userDtos, err := h.userService.FindAll(ctx)
	if err != nil {
		data["users"] = []model.UserDTO{}
		h.renderHTML(c, http.StatusOK, "campaign_add.html", data)
		return
	}

	var onlyUsers []model.UserDTO
	for _, userDto := range userDtos {
		if userDto.Role == domain.RoleUser {
			onlyUsers = append(onlyUsers, userDto)
		}
	}

	data["users"] = onlyUsers

	h.renderHTML(c, http.StatusOK, "campaign_add.html", data)
}
