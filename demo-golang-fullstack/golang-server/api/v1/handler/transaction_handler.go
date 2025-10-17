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

type TransactionHandlerV1 struct {
	transactionUseCase usecaseV1.ITransactionUseCaseV1
	log                zerolog.Logger
}

func NewTransactionHanderV1(
	transactionUseCase usecaseV1.ITransactionUseCaseV1,
	log zerolog.Logger,
) *TransactionHandlerV1 {
	return &TransactionHandlerV1{
		transactionUseCase: transactionUseCase,
		log:                log,
	}
}

// FindAllByUser godoc
// @Summary      Get all user transactions
// @Description  Get all transactions for logged in user
// @Tags         Transaction
// @Produce      json
// @Success      200  {array}  payloadV1.TransactionResponse
// @Failure      401  {object} response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /transactions/me [get]
// @Security     BearerAuth
func (h *TransactionHandlerV1) FindAllByUser(ctx *gin.Context) {
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

	result, err := h.transactionUseCase.FindAllByUser(ctx.Request.Context(), userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			h.log.Error().Err(ex.Err).Msg(ex.Message)
			response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		h.log.Error().Err(err).Msg("failed to get transactions by user")
		response.Error(ctx, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Transactions found", result)
}

// FindAllByCampaign godoc
// @Summary      Get all transactions by campaign
// @Description  Get all transactions for a campaign
// @Tags         Transaction
// @Produce      json
// @Param        id   path      string  true  "Campaign ID"
// @Success      200  {array}  payloadV1.TransactionResponse
// @Failure      404  {object} response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /transactions/campaign/{id} [get]
// @Security     BearerAuth
func (h *TransactionHandlerV1) FindAllByCampaign(ctx *gin.Context) {
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

	result, err := h.transactionUseCase.FindAllByCampaign(ctx.Request.Context(), campaignID, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			h.log.Error().Err(ex.Err).Msg(ex.Message)
			response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		h.log.Error().Err(err).Msg("failed to get transactions by campaign")
		response.Error(ctx, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Transactions found", result)
}

// FindByID godoc
// @Summary      Get transaction detail
// @Description  Get transaction detail by ID
// @Tags         Transaction
// @Produce      json
// @Param        id   path      string  true  "Transaction ID"
// @Success      200  {object} payloadV1.TransactionResponse
// @Failure      404  {object} response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /transactions/{id} [get]
// @Security     BearerAuth
func (h *TransactionHandlerV1) FindByID(ctx *gin.Context) {
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

	txID := ctx.Param("id")

	result, err := h.transactionUseCase.FindByID(ctx.Request.Context(), txID, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			h.log.Error().Err(ex.Err).Msg(ex.Message)
			response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		h.log.Error().Err(err).Msg("failed to get transaction by id")
		response.Error(ctx, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Transaction found", result)
}

// Donation godoc
// @Summary      Create donation transaction
// @Description  Create a new donation transaction
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Param        request  body      payloadV1.TransactionRequest  true  "Donation Request"
// @Success      201  {object} map[string]string
// @Failure      400  {object} response.ErrorResponse
// @Failure      401  {object} response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /transactions/donation [post]
// @Security     BearerAuth
func (h *TransactionHandlerV1) Donation(ctx *gin.Context) {
	var req payloadV1.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Warn().Err(err).Msg("request payload")
		msg := validation.ExtractValidationError(err)
		response.Error(ctx, http.StatusBadRequest, msg, err.Error())
		return
	}

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

	txID, paymentUrl, err := h.transactionUseCase.Create(ctx.Request.Context(), req.Amount, req.CampaignID, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			h.log.Error().Err(ex.Err).Msg(ex.Message)
			response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		h.log.Error().Err(err).Msg("failed to create donation")
		response.Error(ctx, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, "Donation created", gin.H{
		"transaction_id": txID,
		"payment_url":    paymentUrl,
	})
}
