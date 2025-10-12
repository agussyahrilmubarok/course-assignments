package handler

import (
	"net/http"

	"example.com/internal/account/api/v1/payload"
	"github.com/labstack/echo/v4"
)

type authHandler struct {
}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}

func (h *authHandler) SignUp(c echo.Context) error {
	var req payload.SignUpRequest
	if err := c.Bind(&req); err != nil {
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "invalid request body",
			})
		}

	}

	return c.JSON(http.StatusCreated, nil)
}
