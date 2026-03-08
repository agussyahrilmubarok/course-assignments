package handlerV1

import (
	"errors"
	"net/http"

	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/logger"
	"example.com.backend/pkg/response"
	"example.com.backend/pkg/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	payloadV1 "example.com.backend/internal/rest/v1/payload"
	usecaseV1 "example.com.backend/internal/rest/v1/usecase"
)

type CampaignHandlerV1 struct {
	campaignUseCase usecaseV1.ICampaignUseCaseV1
}

func NewCampaignHandlerV1(
	campaignUseCase usecaseV1.ICampaignUseCaseV1,
) *CampaignHandlerV1 {
	return &CampaignHandlerV1{
		campaignUseCase: campaignUseCase,
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
func (h *CampaignHandlerV1) FindAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	campaigns, err := h.campaignUseCase.FindAll(ctx)
	if err != nil {
		log.Error("failed to get campaigns", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "failed to get campaigns", err.Error())
		return
	}

	log.Info("campaigns retrieved", zap.Int("count", len(campaigns)))
	response.Success(c, http.StatusOK, "success", campaigns)
}

// FindAllTop godoc
// @Summary      Get Top Campaigns
// @Description  Get top 6 campaigns ordered by current amount
// @Tags         Campaigns
// @Produce      json
// @Success      200  {array}  payloadV1.CampaignResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /campaigns/top [get]
func (h *CampaignHandlerV1) FindAllTop(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	campaigns, err := h.campaignUseCase.FindAllTop(ctx)
	if err != nil {
		log.Error("failed to get top campaigns", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "failed to get top campaigns", err.Error())
		return
	}

	log.Info("top campaigns retrieved", zap.Int("count", len(campaigns)))
	response.Success(c, http.StatusOK, "success", campaigns)
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
func (h *CampaignHandlerV1) FindByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	id := c.Param("id")

	campaign, err := h.campaignUseCase.FindByID(ctx, id)
	if err != nil {
		log.Warn("campaign not found", zap.String("campaign_id", id), zap.Error(err))
		response.Error(c, http.StatusNotFound, "campaign not found", err.Error())
		return
	}

	log.Info("campaign retrieved", zap.String("campaign_id", id))
	response.Success(c, http.StatusOK, "success", campaign)
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
func (h *CampaignHandlerV1) FindAllByUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing user id in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid user id type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	campaigns, err := h.campaignUseCase.FindAllByUser(ctx, userID)
	if err != nil {
		log.Error("failed to get campaigns", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "failed to get campaigns", err.Error())
		return
	}

	log.Info("campaigns retrieved", zap.String("user_id", userID), zap.Int("count", len(campaigns)))
	response.Success(c, http.StatusOK, "success", campaigns)
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
func (h *CampaignHandlerV1) FindByIDByUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	id := c.Param("id")

	campaign, err := h.campaignUseCase.FindByID(ctx, id)
	if err != nil {
		log.Warn("campaign not found", zap.String("campaign_id", id), zap.Error(err))
		response.Error(c, http.StatusNotFound, "campaign not found", err.Error())
		return
	}

	log.Info("campaign retrieved", zap.String("campaign_id", id))
	response.Success(c, http.StatusOK, "success", campaign)
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
func (h *CampaignHandlerV1) CreateByUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing user id in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid user id type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	var req payloadV1.CampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid campaign request", zap.Error(err))
		msg := validation.ExtractValidationError(err)
		response.Error(c, http.StatusBadRequest, msg, err.Error())
		return
	}

	campaign, err := h.campaignUseCase.CreateByUser(ctx, req, userID)
	if err != nil {
		log.Error("failed to create campaign", zap.Error(err), zap.String("user_id", userID))
		response.Error(c, http.StatusInternalServerError, "failed to create campaign", err.Error())
		return
	}

	log.Info("campaign created", zap.String("campaign_id", campaign.ID), zap.String("user_id", userID))
	response.Success(c, http.StatusCreated, "Campaign created", campaign)
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
func (h *CampaignHandlerV1) UpdateByIDByUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing user id in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid user id type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	id := c.Param("id")

	var req payloadV1.CampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid campaign request", zap.String("campaign_id", id), zap.Error(err))
		msg := validation.ExtractValidationError(err)
		response.Error(c, http.StatusBadRequest, msg, err.Error())
		return
	}

	campaign, err := h.campaignUseCase.UpdateByIDByUser(ctx, req, id, userID)
	if err != nil {
		log.Error("failed to update campaign", zap.Error(err), zap.String("campaign_id", id), zap.String("user_id", userID))
		response.Error(c, http.StatusInternalServerError, "failed to update campaign", err.Error())
		return
	}

	log.Info("campaign updated", zap.String("campaign_id", campaign.ID), zap.String("user_id", userID))
	response.Success(c, http.StatusOK, "Campaign updated", campaign)
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
func (h *CampaignHandlerV1) UploadImageByUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing user id in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid user id type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	campaignID := c.Param("id")

	file, err := c.FormFile("campaign_image")
	if err != nil {
		log.Warn("invalid file upload", zap.String("campaign_id", campaignID), zap.Error(err))
		response.Error(c, http.StatusBadRequest, "invalid file", err.Error())
		return
	}

	isPrimary := c.PostForm("is_primary") == "true"

	param := payloadV1.CampaignImageRequest{
		CampaignID:    campaignID,
		UserID:        userID,
		CampaignImage: file,
		IsPrimary:     isPrimary,
	}

	if err := h.campaignUseCase.UploadImageByUser(ctx, param); err != nil {
		log.Error("failed to upload image", zap.Error(err), zap.String("campaign_id", campaignID), zap.String("user_id", userID))
		response.Error(c, http.StatusInternalServerError, "failed to upload image", err.Error())
		return
	}

	log.Info("campaign image uploaded", zap.String("campaign_id", campaignID), zap.String("user_id", userID), zap.Bool("is_primary", isPrimary))
	response.Success(c, http.StatusOK, "Image campaign uploaded", gin.H{"is_uploaded": true})
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
func (h *CampaignHandlerV1) DeleteByUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing user id in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid user id type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	id := c.Param("id")

	if err := h.campaignUseCase.DeleteByIDByUser(ctx, id, userID); err != nil {
		log.Error("failed to delete campaign", zap.Error(err), zap.String("campaign_id", id), zap.String("user_id", userID))
		response.Error(c, http.StatusInternalServerError, "Failed to delete campaign", err.Error())
		return
	}

	log.Info("campaign deleted", zap.String("campaign_id", id), zap.String("user_id", userID))
	response.Success(c, http.StatusOK, "Campaign deleted", gin.H{"deleted": true})
}
