package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type dashboardController struct {
	baseController
}

func NewDashboardController() *dashboardController {
	return &dashboardController{}
}

func (h *dashboardController) Index(c *gin.Context) {
	data := gin.H{
		"title": "Dashboard",
	}

	h.renderHTML(c, http.StatusOK, "dashboard_index.html", data)
}

func (h *dashboardController) Logout(c *gin.Context) {
	h.deleteUserSession(c)

	c.Redirect(http.StatusFound, "/login")
}
