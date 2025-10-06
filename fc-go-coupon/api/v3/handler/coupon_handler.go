package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type couponHandler struct {
}

func NewCouponHandler() *couponHandler {
	return &couponHandler{}
}

// GetCoupons godoc
// @Summary List users
// @Tags Coupons
// @Produce json
// @Success 200 {array} string
// @Router /coupons [get]
func (h *couponHandler) GetCoupons(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"Coupon", "V1"})

}
