package controller

import (
	"encoding/json"

	"example.com.backend/internal/model"
	"example.com.backend/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type baseController struct {
}

func NewBaseController() *baseController {
	return &baseController{}
}

func (h *baseController) renderHTML(c *gin.Context, status int, template string, data gin.H) {
	profile := h.getUserSession(c)

	data["profile"] = profile

	c.HTML(status, template, data)
}

func (h *baseController) saveUserSession(c *gin.Context, profile model.UserDTO) {
	log := logger.GetLoggerFromContext(c)
	session := sessions.Default(c)

	pJSON, err := json.Marshal(profile)
	if err != nil {
		log.Warn("failed to marshal profile session",
			zap.Error(err),
		)
		return
	}

	session.Set("profile", pJSON)
	if err := session.Save(); err != nil {
		log.Warn("failed to save profile session",
			zap.Error(err),
		)
	}
}

func (h *baseController) getUserSession(c *gin.Context) model.UserDTO {
	log := logger.GetLoggerFromContext(c)
	session := sessions.Default(c)

	pJSON, ok := session.Get("profile").([]byte)
	if !ok {
		log.Warn("profile session not found")
		return model.UserDTO{}
	}

	var profile model.UserDTO
	if err := json.Unmarshal(pJSON, &profile); err != nil {
		log.Warn("failed to unmarshal profile session",
			zap.Error(err),
		)
		return model.UserDTO{}
	}

	return profile
}

func (h *baseController) deleteUserSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("profile")
	session.Save()
}
