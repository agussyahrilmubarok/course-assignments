package account

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service IService
	logger  *logrus.Logger
}

func NewHandler(service IService, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) SignUp(c *gin.Context) {

}

func (h *Handler) SignIn(c *gin.Context) {

}

func (h *Handler) FindByID(c *gin.Context) {

}
