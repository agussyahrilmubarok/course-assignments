package controller

import (
	"context"
	"fmt"
	"net/http"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	"example.com.backend/internal/service"
	"example.com.backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type transactionController struct {
	baseController
	transactionService service.ITransactionService
	userService        service.IUserService
	campaignService    service.ICampaignService
}

func NewTransactionController(
	transactionService service.ITransactionService,
	userService service.IUserService,
	campaignService service.ICampaignService,
) *transactionController {
	return &transactionController{
		transactionService: transactionService,
		userService:        userService,
		campaignService:    campaignService,
	}
}

func (h *transactionController) Index(c *gin.Context) {
	data := gin.H{"title": "Transactions"}

	ctx := c.Request.Context()

	h.showAllTransactions(c, ctx, data)
}

func (h *transactionController) Add(c *gin.Context) {
	data := gin.H{"title": "New Transaction"}

	ctx := c.Request.Context()

	h.showNewTransactionWithData(c, ctx, data)
}
func (h *transactionController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "New Transaction"}

	var input struct {
		Amount     int    `form:"amount" binding:"required,gt=0"`
		UserID     string `form:"user_id" binding:"required"`
		CampaignID string `form:"campaign_id" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		log.Warn("invalid bind add transaction", zap.Error(err))
		data["form"] = input
		data["error"] = "Invalid input."
		h.showNewTransactionWithData(c, ctx, data)
		return
	}

	transactionDto := model.TransactionDTO{
		Amount:     float64(input.Amount),
		UserID:     input.UserID,
		CampaignID: input.CampaignID,
	}
	if err := h.transactionService.Create(ctx, transactionDto); err != nil {
		log.Error("failed to create transaction", zap.Error(err))
		data["form"] = input
		data["error"] = "An error occurred while processing the request."
		h.showNewTransactionWithData(c, ctx, data)
		return
	}

	data["title"] = "Transaction"
	data["form"] = nil
	data["success"] = fmt.Sprintf("New transaction created successfully: %v", input.Amount)

	h.showAllTransactions(c, ctx, data)
}

func (h *transactionController) Show(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Transaction"}

	idStr := c.Param("id")
	transactionDto, err := h.transactionService.FindByID(ctx, idStr)
	if transactionDto == nil || err != nil {
		log.Error("transaction not found", zap.Error(err), zap.String("id", idStr))
		data["error"] = "Transaction not found."
		h.showAllTransactions(c, ctx, data)
		return
	}

	data["transaction"] = transactionDto

	h.renderHTML(c, http.StatusOK, "transaction_show.html", data)
}

func (h *transactionController) Edit(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Edit Transaction"}

	idStr := c.Param("id")
	transaction, err := h.transactionService.FindByID(ctx, idStr)
	if transaction == nil || err != nil {
		log.Error("transaction not found", zap.Error(err), zap.String("id", idStr))
		data["error"] = "Transaction not found."
		h.showAllTransactions(c, ctx, data)
		return
	}

	data["transaction"] = transaction

	h.renderHTML(c, http.StatusOK, "transaction_edit.html", data)
}

func (h *transactionController) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Edit Transaction"}

	idStr := c.Param("id")
	transactionDto, err := h.transactionService.FindByID(ctx, idStr)
	if transactionDto == nil || err != nil {
		log.Error("transaction not found", zap.Error(err), zap.String("id", idStr))
		data["error"] = "Transaction not found."
		h.showAllTransactions(c, ctx, data)
		return
	}

	var input struct {
		Amount int    `form:"amount" binding:"required"`
		Status string `form:"status" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		log.Warn("invalid bind edit transaction", zap.Error(err))
		data["form"] = input
		data["error"] = "Invalid input."
		h.showNewTransactionWithData(c, ctx, data)
		return
	}

	transactionDto.Amount = float64(input.Amount)
	transactionDto.Status = domain.TransactionStatus(input.Status)

	if err := h.transactionService.Update(ctx, *transactionDto); err != nil {
		log.Error("failed to update transaction", zap.Error(err), zap.String("id", idStr))
		data["transaction"] = input
		data["error"] = "Failed to update transaction."
		h.renderHTML(c, http.StatusInternalServerError, "transaction_edit.html", data)
		return
	}

	data["title"] = "Transactions"
	data["success"] = fmt.Sprintf("Transaction updated successfully: %v", input.Amount)

	h.showAllTransactions(c, ctx, data)
}

func (h *transactionController) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(ctx)
	data := gin.H{"title": "Delete Transaction"}

	idStr := c.Param("id")
	if err := h.transactionService.DeleteByID(ctx, idStr); err != nil {
		log.Error("failed to delete transaction", zap.Error(err), zap.String("id", idStr))
		data["error"] = "Failed to delete a transaction."
		h.showAllTransactions(c, ctx, data)
		return
	}

	data["title"] = "Transactions"
	data["success"] = "Transaction has been deleted successfully"

	h.showAllTransactions(c, ctx, data)
}

func (h *transactionController) showAllTransactions(c *gin.Context, ctx context.Context, data gin.H) {
	transactionDtos, err := h.transactionService.FindAll(ctx)
	if err != nil {
		data["transactions"] = []model.TransactionDTO{}
		h.renderHTML(c, http.StatusBadRequest, "transaction_index.html", data)
		return
	}

	data["transactions"] = transactionDtos

	h.renderHTML(c, http.StatusOK, "transaction_index.html", data)
}

func (h *transactionController) showNewTransactionWithData(c *gin.Context, ctx context.Context, data gin.H) {
	userDtos, err := h.userService.FindAll(ctx)
	if err != nil {
		data["users"] = []model.UserDTO{}
		data["campaigns"] = []model.CampaignDetailDTO{}
		h.renderHTML(c, http.StatusBadRequest, "user_index.html", data)
		return
	}

	campaignDtos, err := h.campaignService.FindAll(ctx)
	if err != nil {
		data["users"] = []model.UserDTO{}
		data["campaigns"] = []model.CampaignDetailDTO{}
		h.renderHTML(c, http.StatusOK, "transaction_add.html", data)
		return
	}

	var onlyUsers []model.UserDTO
	for _, userDto := range userDtos {
		if userDto.Role == domain.RoleUser {
			onlyUsers = append(onlyUsers, userDto)
		}
	}

	data["users"] = onlyUsers
	data["campaigns"] = campaignDtos

	h.renderHTML(c, http.StatusOK, "transaction_add.html", data)
}
