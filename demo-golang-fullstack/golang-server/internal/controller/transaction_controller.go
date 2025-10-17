package controller

import (
	"context"
	"fmt"
	"net/http"

	"example.com/backend/internal/domain"
	"example.com/backend/internal/model"
	"example.com/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type transactionController struct {
	baseController
	transactionService service.ITransactionService
	userService        service.IUserService
	campaignService    service.ICampaignService
	log                zerolog.Logger
}

func NewTransactionController(
	transactionService service.ITransactionService,
	userService service.IUserService,
	campaignService service.ICampaignService,
	log zerolog.Logger,
) *transactionController {
	return &transactionController{
		transactionService: transactionService,
		userService:        userService,
		campaignService:    campaignService,
		log:                log,
	}
}

func (h *transactionController) Index(c *gin.Context) {
	data := gin.H{
		"title": "Transactions",
	}

	ctx := c.Request.Context()
	h.showAllTransactions(c, ctx, data)
}

func (h *transactionController) Add(c *gin.Context) {
	data := gin.H{
		"title": "New Transaction",
	}
	ctx := c.Request.Context()
	h.showNewTransactionWithData(c, ctx, data)
}

func (h *transactionController) Create(c *gin.Context) {
	data := gin.H{
		"title": "New Transaction",
	}

	ctx := c.Request.Context()
	var input struct {
		Amount     int    `form:"amount" binding:"required,gt=0"`
		UserID     string `form:"user_id" binding:"required,gt=0"`
		CampaignID string `form:"campaign_id" binding:"required,gt=0"`
	}
	if err := c.ShouldBind(&input); err != nil {
		h.log.Warn().Msgf("invalid bind add transaction")
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
	err := h.transactionService.Create(ctx, transactionDto)
	if err != nil {
		h.log.Error().Err(err).Msgf("failed create transaction")
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
	data := gin.H{
		"title": "Transaction",
	}

	ctx := c.Request.Context()

	idStr := c.Param("id")
	transactionDto, err := h.transactionService.FindByID(ctx, idStr)
	if transactionDto == nil || err != nil {
		h.log.Error().Err(err).Msgf("transaction is not found")
		data["error"] = "Transaction not found."
		h.showAllTransactions(c, ctx, data)
		return
	}

	data["transaction"] = transactionDto

	h.renderHTML(c, http.StatusOK, "transaction_show.html", data)
}

func (h *transactionController) Edit(c *gin.Context) {
	data := gin.H{
		"title": "Edit Transaction",
	}

	ctx := c.Request.Context()

	idStr := c.Param("id")
	transaction, err := h.transactionService.FindByID(ctx, idStr)
	if transaction == nil || err != nil {
		h.log.Error().Err(err).Msgf("transaction is not found")
		data["error"] = "Transaction not found."
		h.showAllTransactions(c, ctx, data)
		return
	}

	data["transaction"] = transaction
	h.renderHTML(c, http.StatusOK, "transaction_edit.html", data)
}

func (h *transactionController) Update(c *gin.Context) {
	data := gin.H{
		"title": "Edit Transaction",
	}

	ctx := c.Request.Context()

	idStr := c.Param("id")
	transactionDto, err := h.transactionService.FindByID(ctx, idStr)
	if transactionDto == nil || err != nil {
		h.log.Error().Err(err).Msgf("transaction is not found")
		data["error"] = "Transaction not found."
		h.showAllTransactions(c, ctx, data)
		return
	}

	var input struct {
		Amount int    `form:"amount" binding:"required"`
		Status string `form:"status" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		h.log.Warn().Msgf("invalid bind edit transaction")
		data["form"] = input
		data["error"] = "Invalid input."
		h.showNewTransactionWithData(c, ctx, data)
		return
	}

	transactionDto.Amount = float64(input.Amount)
	transactionDto.Status = domain.TransactionStatus(input.Status)

	err = h.transactionService.Update(ctx, *transactionDto)
	if err != nil {
		h.log.Error().Err(err).Msgf("failed to update transaction")
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
	data := gin.H{
		"title": "Delete Transaction",
	}

	idStr := c.Param("id")
	err := h.transactionService.DeleteByID(ctx, idStr)
	if err != nil {
		h.log.Error().Err(err).Msgf("failed to delete transaction")
		data["error"] = "Failed to delete a transaction."
		h.showAllTransactions(c, ctx, data)
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
