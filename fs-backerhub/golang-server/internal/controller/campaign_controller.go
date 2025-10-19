package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"example.com/backend/internal/domain"
	"example.com/backend/internal/model"
	"example.com/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	CAMPAIGN_UPLOAD_PATH = "public/uploads/campaigns"
)

type campaignController struct {
	baseController
	campaignService service.ICampaignService
	userService     service.IUserService
	uploadServie    service.IUploadService
	log             zerolog.Logger
}

func NewCampaignController(
	campaignService service.ICampaignService,
	userService service.IUserService,
	uploadServie service.IUploadService,
	log zerolog.Logger,
) *campaignController {
	return &campaignController{
		campaignService: campaignService,
		userService:     userService,
		uploadServie:    uploadServie,
		log:             log,
	}
}

func (h *campaignController) Index(c *gin.Context) {
	data := gin.H{
		"title": "Campaigns",
	}

	ctx := c.Request.Context()

	h.showAllCampaigns(c, ctx, data)
}

func (h *campaignController) Add(c *gin.Context) {
	data := gin.H{
		"title": "New Campaign",
	}
	ctx := c.Request.Context()
	h.showNewCampaignWithData(c, ctx, data)
}

func (h *campaignController) Create(c *gin.Context) {
	data := gin.H{
		"title": "New Campaign",
	}

	ctx := c.Request.Context()
	var input struct {
		Title            string `form:"title" binding:"required"`
		ShortDescription string `form:"short_description" binding:"required"`
		Description      string `form:"description" binding:"required"`
		GoalAmount       int    `form:"goal_amount" binding:"required,gt=0"`
		Perks            string `form:"perks" binding:"required"`
		UserID           string `form:"user_id" binding:"required,gt=0"`
	}
	if err := c.ShouldBind(&input); err != nil {
		h.log.Warn().Msgf("invalid bind add campaign form")
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

	err := h.campaignService.Create(ctx, campaignDto)
	if err != nil {
		h.log.Error().Err(err).Msgf("campaign createion failed")
		data["form"] = input
		data["error"] = "An error occurred while processing the request."
		h.showNewCampaignWithData(c, ctx, data)
		return
	}

	data["title"] = "Campaign"
	data["form"] = nil
	data["success"] = fmt.Sprintf("New campaign created successfully: %s", input.Title)

	h.showAllCampaigns(c, ctx, data)
}

func (h *campaignController) Show(c *gin.Context) {
	data := gin.H{
		"title": "Campaign",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if campaignDto == nil || err != nil {
		h.log.Error().Err(err).Msgf("campaign is not found %s", idStr)
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	data["campaign"] = campaignDto

	h.renderHTML(c, http.StatusOK, "campaign_show.html", data)
}

func (h *campaignController) Edit(c *gin.Context) {
	data := gin.H{
		"title": "Edit Campaign",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		h.log.Error().Err(err).Msgf("campaign is not found %s", idStr)
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	data["campaign"] = campaignDto

	h.renderHTML(c, http.StatusOK, "campaign_edit.html", data)
}

func (h *campaignController) Update(c *gin.Context) {
	data := gin.H{
		"title": "Edit Campaign",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		h.log.Error().Err(err).Msgf("campaign is not found %s", idStr)
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	var input struct {
		ID               string
		Title            string `form:"title" binding:"required"`
		ShortDescription string `form:"short_description" binding:"required"`
		Description      string `form:"description" binding:"required"`
		GoalAmount       int    `form:"goal_amount" binding:"required"`
		CurrentAmount    int    `form:"current_amount"`
		BackerCount      int    `form:"backer_count"`
		Perks            string `form:"perks" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		h.log.Warn().Msgf("invalid bind edit campaign form")
		data["campaign"] = campaignDto
		data["error"] = "Invalid input."
		h.renderHTML(c, http.StatusBadRequest, "campaign_edit.html", data)
		return
	}

	var campaign model.CampaignDTO
	campaign.ID = campaignDto.ID
	campaign.Title = input.Title
	campaign.ShortDescription = input.ShortDescription
	campaign.GoalAmount = float64(input.GoalAmount)
	campaign.CurrentAmount = float64(input.CurrentAmount)
	campaign.BackerCount = int64(input.BackerCount)
	campaign.Perks = input.Perks

	err = h.campaignService.Update(ctx, campaign)
	if err != nil {
		h.log.Error().Err(err).Msgf("failed to update campaign")
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
	data := gin.H{
		"title": "Upload Campaign Image",
	}

	ctx := c.Request.Context()

	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		h.log.Error().Err(err).Msgf("campaign is not found %s", idStr)
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	data["campaign"] = campaignDto

	h.renderHTML(c, http.StatusOK, "campaign_edit_image.html", data)
}

func (h *campaignController) Upload(c *gin.Context) {
	data := gin.H{
		"title": "Upload Campaign Image",
	}

	ctx := c.Request.Context()
	idStr := c.Param("id")
	campaignDto, err := h.campaignService.FindByID(ctx, idStr)
	if err != nil || campaignDto == nil {
		h.log.Error().Err(err).Msgf("campaign is not found %s", idStr)
		data["error"] = "Campaign not found."
		h.showAllCampaigns(c, ctx, data)
		return
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		h.log.Warn().Msgf("failed to get image file")
		data["error"] = "Failed to get image file."
		h.renderHTML(c, http.StatusBadRequest, "campaign_edit_image.html", data)
		return
	}

	imagePath, err := h.uploadServie.SaveLocal(CAMPAIGN_UPLOAD_PATH, imageFile, idStr)
	if err != nil {
		h.log.Err(err).Msgf("failed to upload image file %s", idStr)
		data["error"] = "Failed to upload image file."
		h.renderHTML(c, http.StatusInternalServerError, "campaign_edit_image.html", data)
		return
	}

	trimmedPath := strings.TrimPrefix(imagePath, fmt.Sprintf("%v/", CAMPAIGN_UPLOAD_PATH))
	err = h.campaignService.UploadImage(ctx, model.CampaignImageDTO{
		ImageName:  trimmedPath,
		CampaignID: campaignDto.ID,
	})
	if err != nil {
		h.log.Error().Err(err).Msgf("failed to upload image campaign")
		h.uploadServie.RemoveLocal(fmt.Sprintf("%v/", CAMPAIGN_UPLOAD_PATH), trimmedPath)
		data["error"] = "Failed to upload image file."
		h.renderHTML(c, http.StatusInternalServerError, "campaign_edit_image.html", data)
		return
	}

	data["title"] = "Campaigns"
	data["success"] = "Campaign image uploaded successfully"

	h.showAllCampaigns(c, ctx, data)
}

func (h *campaignController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	data := gin.H{
		"title": "Delete User",
	}

	idStr := c.Param("id")
	campaignDto, _ := h.campaignService.FindByID(ctx, idStr)
	for _, campaignImage := range campaignDto.CampaignImages {
		if err := h.uploadServie.RemoveLocal(fmt.Sprintf("%v/", CAMPAIGN_UPLOAD_PATH), campaignImage.ImageName); err != nil {
			h.log.Error().Err(err).Msgf("failed to remove campaign image for id %s", campaignImage.ID)
			continue
		}
	}

	err := h.campaignService.DeleteByID(ctx, idStr)
	if err != nil {
		h.log.Error().Err(err).Msgf("campaign is not found %s", idStr)
		data["error"] = "Failed to delete a campaign."
		h.showAllCampaigns(c, ctx, data)
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
