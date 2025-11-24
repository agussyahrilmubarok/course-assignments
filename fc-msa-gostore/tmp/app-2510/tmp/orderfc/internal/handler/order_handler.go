package handler

import (
	"net/http"

	"example.com/orderfc/internal/service"
	"example.com/pkg/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type orderHandler struct {
	orderService service.IOrderService
	log          *zap.Logger
}

func NewOrderHandler(orderService service.IOrderService, log *zap.Logger) *orderHandler {
	return &orderHandler{
		orderService: orderService,
		log:          log,
	}
}

func (h *orderHandler) CheckoutOrder(c *gin.Context) {
	var req model.CheckoutOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid checkout order request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	response, err := h.orderService.Checkout(c.Request.Context(), &req)
	if err != nil || response == nil {
		h.log.Error("failed to checkout order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to checkout order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *orderHandler) CancelOrder(c *gin.Context) {
	// TODO: validate true owner
	var req model.OrderCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid cancel order request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	response, err := h.orderService.Cancel(c.Request.Context(), req.OrderID)
	if err != nil || response == nil {
		h.log.Error("failed to cancel order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *orderHandler) GerOrderHistory(c *gin.Context) {
	userID := c.GetHeader("X-USER-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-USER-ID header"})
		return
	}

	response, err := h.orderService.GetHistory(c.Request.Context(), userID)
	if err != nil || response == nil {
		h.log.Error("failed to get order history", zap.String("user_id", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order history", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
