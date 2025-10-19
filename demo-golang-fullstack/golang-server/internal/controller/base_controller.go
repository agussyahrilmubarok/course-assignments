package controller

import (
	"encoding/json"

	"example.com/backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type baseController struct {
	log zerolog.Logger
}

func NewBaseController(
	log zerolog.Logger,
) *baseController {
	return &baseController{
		log: log,
	}
}

func (h *baseController) renderHTML(c *gin.Context, status int, template string, data gin.H) {
	profile := h.getUserSession(c)
	data["profile"] = profile
	c.HTML(status, template, data)
}

func (h *baseController) saveUserSession(c *gin.Context, profile model.UserDTO) {
	session := sessions.Default(c)
	pJSON, err := json.Marshal(profile)
	if err != nil {
		h.log.Warn().Msgf("failed to save profile session, err %v", err)
	}
	session.Set("profile", pJSON)
	if err := session.Save(); err != nil {
		h.log.Warn().Msgf("failed to save profile session, err %v", err)
	}
}

func (h *baseController) getUserSession(c *gin.Context) model.UserDTO {
	session := sessions.Default(c)
	pJSON, ok := session.Get("profile").([]byte)
	if !ok {
		h.log.Warn().Msgf("failed to get profile session")
	}
	var profile model.UserDTO
	if err := json.Unmarshal(pJSON, &profile); err != nil {
		h.log.Warn().Msgf("failed to get profile session, err %v", err)
	}
	return profile
}

func (h *baseController) deleteUserSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("profile")
	session.Save()
}
