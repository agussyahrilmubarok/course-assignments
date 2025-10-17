package handlerV1

import (
	"errors"
	"net/http"

	payloadV1 "example.com/backend/api/v1/payload"
	usecaseV1 "example.com/backend/api/v1/usecase"
	"example.com/backend/internal/exception"
	"example.com/backend/pkg/response"
	"example.com/backend/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CampaignHandlerV1 struct {
	campaignUseCase usecaseV1.ICampaignUseCaseV1
	log             zerolog.Logger
}

func NewCampaignHandlerV1(
	campaignUseCase usecaseV1.ICampaignUseCaseV1,
	log zerolog.Logger,
) *CampaignHandlerV1 {
	return &CampaignHandlerV1{
		campaignUseCase: campaignUseCase,
		log:             log,
	}
}

// FindAll godoc
// @Summary      Get All Campaigns
// @Description  Get list of campaigns
// @Tags         Campaigns
// @Produce      json
// @Success      200  {array}  payloadV1.CampaignResponse
// @Failure      400  {object} response.ErrorResponse
// @Router       /campaigns [get]
func (h *CampaignHandlerV1) FindAll(ctx *gin.Context) {
	campaigns, err := h.campaignUseCase.FindAll(ctx.Request.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("failed to get campaigns")
		response.Error(ctx, http.StatusBadRequest, "failed to get campaigns", err.Error())
		return
	}

	h.log.Info().Int("count", len(campaigns)).Msg("campaigns retrieved")
	response.Success(ctx, http.StatusOK, "success", campaigns)
}

// FindAllTop godoc
// @Summary      Get Top Campaigns
// @Description  Get top 6 campaigns ordered by current amount
// @Tags         Campaigns
// @Produce      json
// @Success      200  {array}  payloadV1.CampaignResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /campaigns/top [get]
func (h *CampaignHandlerV1) FindAllTop(ctx *gin.Context) {
	campaigns, err := h.campaignUseCase.FindAllTop(ctx.Request.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("failed to get top campaigns")
		response.Error(ctx, http.StatusInternalServerError, "failed to get top campaigns", err.Error())
		return
	}

	h.log.Info().Int("count", len(campaigns)).Msg("top campaigns retrieved")
	response.Success(ctx, http.StatusOK, "success", campaigns)
}

// FindByID godoc
// @Summary      Get Campaign by ID
// @Description  Get campaign detail by ID
// @Tags         Campaigns
// @Produce      json
// @Param        id   path      string  true  "Campaign ID"
// @Success      200  {object}  payloadV1.CampaignDetailResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /campaigns/{id} [get]
func (h *CampaignHandlerV1) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	campaign, err := h.campaignUseCase.FindByID(ctx.Request.Context(), id)
	if err != nil {
		h.log.Warn().Str("campaign_id", id).Err(err).Msg("campaign not found")
		response.Error(ctx, http.StatusNotFound, "campaign not found", err.Error())
		return
	}

	h.log.Info().Str("campaign_id", id).Msg("campaign retrieved")
	response.Success(ctx, http.StatusOK, "success", campaign)
}

// FindAllByUser godoc
// @Summary      Get All Campaigns By User
// @Description  Get list of campaigns by user
// @Tags         Campaigns
// @Produce      json
// @Success      200  {array}  payloadV1.CampaignResponse
// @Failure      400  {object} response.ErrorResponse
// @Router       /campaigns/me [get]
// @Security     BearerAuth
func (h *CampaignHandlerV1) FindAllByUser(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	campaigns, err := h.campaignUseCase.FindAllByUser(ctx.Request.Context(), userID)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to get campaigns")
		response.Error(ctx, http.StatusBadRequest, "failed to get campaigns", err.Error())
		return
	}

	h.log.Info().Int("count", len(campaigns)).Msg("campaigns retrieved")
	response.Success(ctx, http.StatusOK, "success", campaigns)
}

// FindByIDByUser godoc
// @Summary      Get Campaign by ID by User
// @Description  Get campaign detail by ID by User
// @Tags         Campaigns
// @Produce      json
// @Param        id   path      string  true  "Campaign ID"
// @Success      200  {object}  payloadV1.CampaignDetailResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /campaigns/{id}/me [get]
func (h *CampaignHandlerV1) FindByIDByUser(ctx *gin.Context) {
	id := ctx.Param("id")

	campaign, err := h.campaignUseCase.FindByID(ctx.Request.Context(), id)
	if err != nil {
		h.log.Warn().Str("campaign_id", id).Err(err).Msg("campaign not found")
		response.Error(ctx, http.StatusNotFound, "campaign not found", err.Error())
		return
	}

	h.log.Info().Str("campaign_id", id).Msg("campaign retrieved")
	response.Success(ctx, http.StatusOK, "success", campaign)
}

// CreateByUser godoc
// @Summary      Create Campaign
// @Description  Create a new campaign by authenticated user
// @Tags         Campaigns
// @Accept       json
// @Produce      json
// @Param        body  body      payloadV1.CampaignRequest  true  "Campaign Request"
// @Success      201   {object}  payloadV1.CampaignResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      401   {object}  response.ErrorResponse
// @Router       /campaigns [post]
// @Security     BearerAuth
func (h *CampaignHandlerV1) CreateByUser(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	var req payloadV1.CampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Warn().Err(err).Msg("invalid campaign request")
		msg := validation.ExtractValidationError(err)
		response.Error(ctx, http.StatusBadRequest, msg, err.Error())
		return
	}

	campaign, err := h.campaignUseCase.CreateByUser(ctx.Request.Context(), req, userID)
	if err != nil {
		h.log.Error().Err(err).Str("user_id", userID).Msg("failed to create campaign")
		response.Error(ctx, http.StatusInternalServerError, "failed to create campaign", err.Error())
		return
	}

	h.log.Info().Str("campaign_id", campaign.ID).Str("user_id", userID).Msg("campaign created")
	response.Success(ctx, http.StatusCreated, "Campaign created", campaign)
}

// UpdateByIDByUser godoc
// @Summary      Update Campaign
// @Description  Update campaign by ID (only owner can update)
// @Tags         Campaigns
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "Campaign ID"
// @Param        body  body      payloadV1.CampaignRequest  true  "Campaign Request"
// @Success      200   {object}  payloadV1.CampaignResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      401   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /campaigns/{id} [put]
// @Security     BearerAuth
func (h *CampaignHandlerV1) UpdateByIDByUser(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	id := ctx.Param("id")

	var req payloadV1.CampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Warn().Err(err).Str("campaign_id", id).Msg("invalid campaign request")
		msg := validation.ExtractValidationError(err)
		response.Error(ctx, http.StatusBadRequest, msg, err.Error())
		return
	}

	campaign, err := h.campaignUseCase.UpdateByIDByUser(ctx.Request.Context(), req, id, userID)
	if err != nil {
		h.log.Error().Err(err).Str("campaign_id", id).Str("user_id", userID).Msg("failed to update campaign")
		response.Error(ctx, http.StatusInternalServerError, "failed to update campaign", err.Error())
		return
	}

	h.log.Info().Str("campaign_id", campaign.ID).Str("user_id", userID).Msg("campaign updated")
	response.Success(ctx, http.StatusOK, "Campaign updated", campaign)
}

// UploadImageByUser godoc
// @Summary      Upload Campaign Image
// @Description  Upload image for campaign (set primary if needed)
// @Tags         Campaigns
// @Accept       multipart/form-data
// @Produce      json
// @Param        id           path      string  true   "Campaign ID"
// @Param        is_primary   formData  bool    true   "Set as primary image"
// @Param        campaign_image formData file true "Campaign Image"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /campaigns/{id}/images [post]
// @Security     BearerAuth
func (h *CampaignHandlerV1) UploadImageByUser(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	campaignID := ctx.Param("id")

	file, err := ctx.FormFile("campaign_image")
	if err != nil {
		h.log.Warn().Err(err).Str("campaign_id", campaignID).Msg("invalid file upload")
		response.Error(ctx, http.StatusBadRequest, "invalid file", err.Error())
		return
	}

	isPrimary := ctx.PostForm("is_primary") == "true"

	param := payloadV1.CampaignImageRequest{
		CampaignID:    campaignID,
		UserID:        userID,
		CampaignImage: file,
		IsPrimary:     isPrimary,
	}

	if err := h.campaignUseCase.UploadImageByUser(ctx.Request.Context(), param); err != nil {
		h.log.Error().Err(err).Str("campaign_id", campaignID).Str("user_id", userID).Msg("failed to upload image")
		response.Error(ctx, http.StatusInternalServerError, "failed to upload image", err.Error())
		return
	}

	h.log.Info().Str("campaign_id", campaignID).Str("user_id", userID).Bool("is_primary", isPrimary).Msg("campaign image uploaded")
	response.Success(ctx, http.StatusOK, "Image campaign uploaded", gin.H{"is_uploaded": true})
}

// DeleteByUser godoc
// @Summary      Delete Campaign
// @Description  Delete campaign by ID (only owner can delete)
// @Tags         Campaigns
// @Produce      json
// @Param        id   path      string  true  "Campaign ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /campaigns/{id} [delete]
// @Security     BearerAuth
func (h *CampaignHandlerV1) DeleteByUser(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		h.log.Error().Err(ex.Err).Msg(ex.Message)
		response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	id := ctx.Param("id")

	if err := h.campaignUseCase.DeleteByIDByUser(ctx.Request.Context(), id, userID); err != nil {
		h.log.Error().Err(err).Str("campaign_id", id).Str("user_id", userID).Msg("failed to delete campaign")
		response.Error(ctx, http.StatusInternalServerError, "failed to delete campaign", err.Error())
		return
	}

	h.log.Info().Str("campaign_id", id).Str("user_id", userID).Msg("campaign deleted")
	response.Success(ctx, http.StatusOK, "Campaign deleted", gin.H{"deleted": true})
}
