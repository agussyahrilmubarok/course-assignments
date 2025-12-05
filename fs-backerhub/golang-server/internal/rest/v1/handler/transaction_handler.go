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

type TransactionHandlerV1 struct {
	transactionUseCase usecaseV1.ITransactionUseCaseV1
}

func NewTransactionHanderV1(
	transactionUseCase usecaseV1.ITransactionUseCaseV1,
) *TransactionHandlerV1 {
	return &TransactionHandlerV1{
		transactionUseCase: transactionUseCase,
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
func (h *TransactionHandlerV1) FindAllByUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	result, err := h.transactionUseCase.FindAllByUser(ctx, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error(ex.Message, zap.Error(ex.Err))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		log.Error("failed to get transactions by user", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	log.Info("transactions retrieved", zap.String("user_id", userID))
	response.Success(c, http.StatusOK, "Transactions found", result)
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
func (h *TransactionHandlerV1) FindAllByCampaign(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	campaignID := c.Param("id")

	result, err := h.transactionUseCase.FindAllByCampaign(ctx, campaignID, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error(ex.Message, zap.Error(ex.Err))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		log.Error("failed to get transactions by campaign", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	log.Info("transactions campaign retrieved", zap.String("campaign_id", campaignID))
	response.Success(c, http.StatusOK, "Transactions found", result)
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
func (h *TransactionHandlerV1) FindByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	txID := c.Param("id")

	result, err := h.transactionUseCase.FindByID(ctx, txID, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error(ex.Message, zap.Error(ex.Err))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		log.Error("failed to get transaction by id", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	log.Info("transaction retrieved", zap.String("transaction_id", txID))
	response.Success(c, http.StatusOK, "Transaction found", result)
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
func (h *TransactionHandlerV1) Donation(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	var req payloadV1.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid transaction request", zap.Error(err))
		msg := validation.ExtractValidationError(err)
		response.Error(c, http.StatusBadRequest, msg, err.Error())
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("missing userID in context"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		ex := exception.NewUnauthorized("Unauthorized", errors.New("invalid userID type"))
		log.Error("unauthorized", zap.Error(ex.Err))
		response.Error(c, ex.Code, ex.Message, ex.Err.Error())
		return
	}

	txID, paymentUrl, err := h.transactionUseCase.Create(ctx, req.Amount, req.CampaignID, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error(ex.Message, zap.Error(ex.Err))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		log.Error("failed to create donation", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Unexpected error", err.Error())
		return
	}

	log.Info("donation created",
		zap.String("transaction_id", txID),
		zap.String("user_id", userID),
		zap.String("campaign_id", req.CampaignID),
	)
	response.Success(c, http.StatusCreated, "Donation created", gin.H{
		"transaction_id": txID,
		"payment_url":    paymentUrl,
	})
}
