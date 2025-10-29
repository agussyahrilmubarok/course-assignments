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
// @Summary List coupons
// @Tags Coupons
// @Produce json
// @Success 200 {array} string
// @Router /coupons [get]
func (h *couponHandler) IssueCoupon(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"Coupon", "V1"})

}

func (h *couponHandler) UseCoupon(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"Coupon", "V1"})
}

func (h *couponHandler) CancelCoupon(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"Coupon", "V1"})
}

func (h *couponHandler) FindCoupon(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"Coupon", "V1"})
}

func (h *couponHandler) FindMyCoupons(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"Coupon", "V1"})
}
