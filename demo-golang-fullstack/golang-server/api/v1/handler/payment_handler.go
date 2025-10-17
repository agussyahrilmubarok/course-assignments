package handlerV1

import (
	"net/http"

	usecaseV1 "example.com/backend/api/v1/usecase"
	"example.com/backend/internal/domain"
	"example.com/backend/internal/exception"
	"example.com/backend/internal/model"
	"example.com/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog"
)

type PaymentHandlerV1 struct {
	transactionUseCase usecaseV1.ITransactionUseCaseV1
	log                zerolog.Logger
}

func NewPaymentHanderV1(
	transactionUseCase usecaseV1.ITransactionUseCaseV1,
	log zerolog.Logger,
) *PaymentHandlerV1 {
	return &PaymentHandlerV1{
		transactionUseCase: transactionUseCase,
		log:                log,
	}
}

// MidtransPaymentNotification godoc
// @Summary      Midtrans payment notification callback
// @Description  Handle notification callback from Midtrans
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Success      200  {object} map[string]string
// @Failure      400  {object} response.ErrorResponse
// @Failure      401  {object} response.ErrorResponse
// @Failure      403  {object} response.ErrorResponse
// @Failure      500  {object} response.ErrorResponse
// @Router       /payments/midtrans/notification [post]
func (h *PaymentHandlerV1) MidtransPaymentNotification(ctx *gin.Context) {
	var input model.MidtransCallback
	var raw map[string]interface{}

	// Bind raw JSON payload
	if err := ctx.ShouldBindJSON(&raw); err != nil {
		h.log.Warn().Err(err).Msg("Invalid midtrans callback payload")
		response.Error(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Decode known fields into MidtransCallback struct
	if err := mapstructure.Decode(raw, &input); err != nil {
		h.log.Warn().Err(err).Msg("Failed to decode midtrans payload")
		response.Error(ctx, http.StatusBadRequest, "Failed to decode known fields", err.Error())
		return
	}

	// Capture extra unknown fields
	input.ExtraFields = make(map[string]interface{})
	for key, val := range raw {
		if !input.IsKnownMidtransField(key) {
			input.ExtraFields[key] = val
		}
	}
	if len(input.ExtraFields) > 0 {
		h.log.Info().Interface("extra_fields", input.ExtraFields).Msg("Midtrans extra fields")
	}

	// Fraud check logic (credit card only)
	if input.PaymentType == "credit_card" {
		switch input.FraudStatus {
		case model.MidtransFraudStatusAccept:
			h.log.Info().Msg("Fraud check passed")
		case model.MidtransFraudStatusDeny:
			h.log.Warn().Str("fraud_status", input.FraudStatus).Msg("Fraud check failed — rejecting")
			response.Error(ctx, http.StatusForbidden, "Transaction marked as fraud by FDS", "fraud deny")
			return
		case "challenge":
			h.log.Warn().Msg("Fraud status challenge — requires manual review")
			// Optional: continue or reject, depending on business rules
		default:
			h.log.Warn().Str("fraud_status", input.FraudStatus).Msg("Unknown fraud status")
			response.Error(ctx, http.StatusForbidden, "Unknown fraud status", input.FraudStatus)
			return
		}
	}

	// Map Midtrans status → internal transaction status
	status := h.mapMidtransStatusToTxEntity(input.TransactionStatus)

	// Pass to use case layer
	if err := h.transactionUseCase.ProcessMidtransPayment(ctx.Request.Context(), &input, status); err != nil {
		if ex, ok := err.(*exception.Http); ok {
			h.log.Error().Err(ex.Err).Msg(ex.Message)
			response.Error(ctx, ex.Code, ex.Message, ex.Err.Error())
			return
		}
		h.log.Error().Err(err).Msg("Failed to process midtrans payment")
		response.Error(ctx, http.StatusInternalServerError, "Failed process midtrans payment", err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Callback received", gin.H{"message": "ok"})
}

func (h *PaymentHandlerV1) mapMidtransStatusToTxEntity(midtransStatus string) domain.TransactionStatus {
	switch midtransStatus {
	case model.MidtransTransactionStatusSettlement, model.MidtransTransactionStatusCapture:
		return domain.StatusPaid
	case model.MidtransTransactionStatusPending, model.MidtransTransactionStatusAuthorize:
		return domain.StatusPending
	case model.MidtransTransactionStatusDeny, model.MidtransTransactionStatusFailure:
		return domain.StatusFailed
	case model.MidtransTransactionStatusCancel:
		return domain.StatusCanceled
	case model.MidtransTransactionStatusExpire:
		return domain.StatusExpired
	case model.MidtransTransactionStatusRefund, model.MidtransTransactionStatusPartialRefund:
		// Optional: treat refund as failed, or create a separate status if needed
		return domain.StatusFailed
	default:
		// Fallback: unknown status, consider it failed or log error
		return domain.StatusFailed
	}
}
