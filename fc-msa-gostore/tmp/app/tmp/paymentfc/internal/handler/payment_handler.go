package handler

import (
	"net/http"

	"example.com/paymentfc/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type paymentHandler struct {
	xenditWebhookToken string
	paymentService     service.IPaymentService
	log                *zap.Logger
}

func NewPaymentHandler(xenditWebhookToken string, paymentService service.IPaymentService, log *zap.Logger) *paymentHandler {
	return &paymentHandler{
		xenditWebhookToken: xenditWebhookToken,
		paymentService:     paymentService,
		log:                log,
	}
}

func (h *paymentHandler) XenditWebhook(c *gin.Context) {
	var req struct {
		ExternalID string  `json:"external_id"`
		Status     string  `json:"status"`
		Amount     float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid xendit webhook request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	headerWebhookToken := c.GetHeader("x-callback-token")
	if h.xenditWebhookToken != headerWebhookToken {
		h.log.Warn("unauthorized webhook callback", zap.String("received_token", headerWebhookToken))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized webhook callback"})
		return
	}

	if err := h.paymentService.ProcessXenditWebhook(c.Request.Context(), req.ExternalID, req.Status, req.Amount); err != nil {
		h.log.Error("failed to process Xendit webhook", zap.Error(err), zap.String("external_id", req.ExternalID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process webhook", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "webhook received"})
}

func (h *paymentHandler) DownloadInvoice(c *gin.Context) {
	orderIDStr := c.Param("order_id")

	filePath, err := h.paymentService.DownloadInvoiceInPDF(c.Request.Context(), orderIDStr)
	if err != nil {
		h.log.Error("failed to download invoice", zap.Error(err), zap.String("order_id", orderIDStr))
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to download invoice", "details": err.Error()})
		return
	}

	c.FileAttachment(filePath, filePath)
}
