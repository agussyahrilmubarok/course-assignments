package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type homeController struct{}

func NewHomeController() *homeController {
	return &homeController{}
}

func (h *homeController) Index(c *gin.Context) {
	data := gin.H{
		"title": "Home",
	}

	c.HTML(http.StatusOK, "home_index.html", data)
}
