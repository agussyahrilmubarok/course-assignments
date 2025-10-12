package account

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SignUp(c echo.Context) {
	
}

func (h *Handler) SignIn(c echo.Context) {

}

func (h *Handler) GetMe(c echo.Context) {

}
