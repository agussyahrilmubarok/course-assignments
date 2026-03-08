package handlerV1

import (
	"net/http"

	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/logger"
	"example.com.backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	usecaseV1 "example.com.backend/internal/rest/v1/usecase"
)

type PaymentHandlerV1 struct {
	transactionUseCase usecaseV1.ITransactionUseCaseV1
}

func NewPaymentHanderV1(
	transactionUseCase usecaseV1.ITransactionUseCaseV1,
) *PaymentHandlerV1 {
	return &PaymentHandlerV1{
		transactionUseCase: transactionUseCase,
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
func (h *PaymentHandlerV1) MidtransPaymentNotification(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)

	var input model.MidtransCallback
	var raw map[string]interface{}

	// Bind raw JSON payload
	if err := c.ShouldBindJSON(&raw); err != nil {
		log.Warn("invalid midtrans callback payload", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Decode known fields into struct
	if err := mapstructure.Decode(raw, &input); err != nil {
		log.Warn("failed to decode midtrans payload", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "Failed to decode known fields", err.Error())
		return
	}

	// Save extra Midtrans unknown fields
	input.ExtraFields = make(map[string]interface{})
	for key, val := range raw {
		if !input.IsKnownMidtransField(key) {
			input.ExtraFields[key] = val
		}
	}
	if len(input.ExtraFields) > 0 {
		log.Info("midtrans callback extra fields",
			zap.Any("extra_fields", input.ExtraFields))
	}

	// Fraud check logic
	if input.PaymentType == "credit_card" {
		switch input.FraudStatus {
		case model.MidtransFraudStatusAccept:
			log.Info("fraud check passed")

		case model.MidtransFraudStatusDeny:
			log.Warn("fraud check denied", zap.String("fraud_status", input.FraudStatus))
			response.Error(c, http.StatusForbidden, "Transaction marked as fraud by FDS", "fraud deny")
			return

		case "challenge":
			log.Warn("fraud status challenge — manual review needed",
				zap.String("fraud_status", input.FraudStatus))

		default:
			log.Warn("unknown fraud status", zap.String("fraud_status", input.FraudStatus))
			response.Error(c, http.StatusForbidden, "Unknown fraud status", input.FraudStatus)
			return
		}
	}

	// Map midtrans → internal status
	status := h.mapMidtransStatusToTxEntity(input.TransactionStatus)

	// Process payment
	if err := h.transactionUseCase.ProcessMidtransPayment(ctx, &input, status); err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error(ex.Message, zap.Error(ex.Err))
			response.Error(c, ex.Code, ex.Message, ex.Err.Error())
			return
		}

		log.Error("failed to process midtrans payment", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Failed process midtrans payment", err.Error())
		return
	}

	log.Info("midtrans callback processed",
		zap.String("transaction_status", input.TransactionStatus),
		zap.String("order_id", input.OrderID),
		zap.String("payment_type", input.PaymentType),
	)
	response.Success(c, http.StatusOK, "Callback received", gin.H{"message": "ok"})
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
