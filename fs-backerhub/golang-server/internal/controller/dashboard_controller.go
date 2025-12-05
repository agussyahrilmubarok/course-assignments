package controller

import (
	"net/http"

	"example.com.backend/internal/repos"
	"github.com/gin-gonic/gin"
)

type dashboardController struct {
	baseController
	userRepo        repos.IUserRepository
	campaignRepo    repos.ICampaignRepository
	transactionRepo repos.ITransactionRepository
}

func NewDashboardController(
	userRepo repos.IUserRepository,
	campaignRepo repos.ICampaignRepository,
	transactionRepo repos.ITransactionRepository,
) *dashboardController {
	return &dashboardController{
		userRepo:        userRepo,
		campaignRepo:    campaignRepo,
		transactionRepo: transactionRepo,
	}
}

func (h *dashboardController) Index(c *gin.Context) {
	ctx := c.Request.Context()
	data := gin.H{"title": "Dashboard"}

	totalMemberActive, _ := h.userRepo.CountMemberActive(ctx)
	totalCampaignActive, _ := h.campaignRepo.CountActive(ctx)
	totalTransactionPending, _ := h.transactionRepo.CountPending(ctx)

	data["totalMemberActive"] = totalMemberActive
	data["totalCampaignActive"] = totalCampaignActive
	data["totalTransactionPending"] = totalTransactionPending

	h.renderHTML(c, http.StatusOK, "dashboard_index.html", data)
}

func (h *dashboardController) Logout(c *gin.Context) {
	h.deleteUserSession(c)

	c.Redirect(http.StatusFound, "/login")
}
