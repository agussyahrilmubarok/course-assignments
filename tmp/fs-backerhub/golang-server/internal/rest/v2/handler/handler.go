package handlerV2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// API V2 godoc
// @Summary      API V2
// @Description  Returns a message
// @Tags         Public
// @Produce      json
// @Success      200  {object} map[string]string
// @Router       / [get]
func (h *Handler) V2(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, Welcome to Backerhub API V2!",
	})
}
