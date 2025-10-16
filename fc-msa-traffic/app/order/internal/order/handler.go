package order

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Handler struct {
	store IStore
	log   zerolog.Logger
}

func NewHandler(
	store IStore,
	log zerolog.Logger,
) *Handler {
	return &Handler{
		store: store,
		log:   log,
	}
}

func (h *Handler) CreateOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) CancelOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
