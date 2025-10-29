package handler

import (
	"net/http"

	"example.com/coupon/internal/coupon"
	v1 "example.com/coupon/internal/coupon/feature/v1"
	"example.com/coupon/internal/middleware"
	"example.com/coupon/pkg/exception"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type couponHandler struct {
	couponFeature v1.ICouponFeature
}

func NewCouponHandler(couponFeature v1.ICouponFeature) *couponHandler {
	return &couponHandler{
		couponFeature: couponFeature,
	}
}

// IssueCoupon godoc
// @Summary Issue a coupon for the current user
// @Description Issues a coupon based on the given coupon policy code
// @Tags Coupons
// @Accept json
// @Produce json
// @Param payload body coupon.IssueCouponRequest true "Coupon issue request"
// @Success 200 {object} coupon.Coupon
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coupons/issue [post]
// @Security ApiKeyAuth
func (h *couponHandler) IssueCoupon(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	var payload coupon.IssueCouponRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		log.Warn().Err(err).Msg("User ID not found in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	issuedCoupon, err := h.couponFeature.IssueCoupon(c.Request.Context(), payload.CouponPolicyCode, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error().
				Err(ex.Err).
				Str("coupon_policy_code", payload.CouponPolicyCode).
				Str("user_id", userID).
				Msg("Failed to issue coupon")
			c.JSON(ex.Code, gin.H{"error": ex.Message})
			return
		}

		log.Error().
			Err(err).
			Str("coupon_policy_code", payload.CouponPolicyCode).
			Str("user_id", userID).
			Msg("Unexpected error while issuing coupon")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to issue coupon"})
		return
	}

	log.Info().
		Str("coupon_code", issuedCoupon.Code).
		Str("coupon_policy_code", payload.CouponPolicyCode).
		Str("user_id", userID).
		Msg("Coupon issued successfully")
	c.JSON(http.StatusOK, issuedCoupon)
}

// UseCoupon godoc
// @Summary Use a coupon for an order
// @Description Marks the coupon as used for a given order
// @Tags Coupons
// @Accept json
// @Produce json
// @Param code path string true "Coupon Code"
// @Param payload body coupon.UseCouponRequest true "Use coupon request"
// @Success 200 {object} coupon.Coupon
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coupons/{code}/use [post]
// @Security ApiKeyAuth
func (h *couponHandler) UseCoupon(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	couponCode := c.Param("code")
	var payload coupon.UseCouponRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn().Err(err).Msg("Invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		log.Warn().Err(err).Msg("User ID not found in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	usedCoupon, err := h.couponFeature.UseCoupon(c.Request.Context(), couponCode, userID, payload.OrderID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error().
				Err(ex.Err).
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Str("order_id", payload.OrderID).
				Msg("Failed to use coupon")
			c.JSON(ex.Code, gin.H{"error": ex.Message})
			return
		}

		log.Error().
			Err(err).
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Str("order_id", payload.OrderID).
			Msg("Unexpected error while using coupon")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to use coupon"})
		return
	}

	log.Info().
		Str("coupon_code", usedCoupon.Code).
		Str("user_id", userID).
		Str("order_id", payload.OrderID).
		Msg("Coupon used successfully")
	c.JSON(http.StatusOK, usedCoupon)
}

// CancelCoupon godoc
// @Summary Cancel a coupon
// @Description Cancels a previously issued coupon
// @Tags Coupons
// @Produce json
// @Param code path string true "Coupon Code"
// @Success 200 {object} coupon.Coupon
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coupons/{code}/cancel [post]
// @Security ApiKeyAuth
func (h *couponHandler) CancelCoupon(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	couponCode := c.Param("code")
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		log.Warn().Err(err).Msg("User ID not found in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	canceledCoupon, err := h.couponFeature.CancelCoupon(c.Request.Context(), couponCode, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error().
				Err(ex.Err).
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Msg("Failed to cancel coupon")
			c.JSON(ex.Code, gin.H{"error": ex.Message})
			return
		}

		log.Error().
			Err(err).
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Msg("Unexpected error while canceling coupon")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel coupon"})
		return
	}

	log.Info().
		Str("coupon_code", canceledCoupon.Code).
		Str("user_id", userID).
		Msg("Coupon canceled successfully")
	c.JSON(http.StatusOK, canceledCoupon)
}

// FindCouponByCode godoc
// @Summary Get coupon by code
// @Description Retrieves coupon details by code for the current user
// @Tags Coupons
// @Produce json
// @Param code path string true "Coupon Code"
// @Success 200 {object} coupon.Coupon
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coupons/{code} [get]
// @Security ApiKeyAuth
func (h *couponHandler) FindCouponByCode(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	couponCode := c.Param("code")
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		log.Warn().Err(err).Msg("User ID not found in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	couponData, err := h.couponFeature.FindCouponByCode(c.Request.Context(), couponCode, userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error().
				Err(ex.Err).
				Str("coupon_code", couponCode).
				Str("user_id", userID).
				Msg("Failed to find coupon by code")
			c.JSON(ex.Code, gin.H{"error": ex.Message})
			return
		}

		log.Error().
			Err(err).
			Str("coupon_code", couponCode).
			Str("user_id", userID).
			Msg("Unexpected error while finding coupon by code")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find coupon"})
		return
	}

	log.Info().
		Str("coupon_code", couponCode).
		Str("user_id", userID).
		Msg("Coupon found successfully")
	c.JSON(http.StatusOK, couponData)
}

// FindCouponsByUserID godoc
// @Summary Get all coupons of current user
// @Description Retrieves all coupons owned by the current user
// @Tags Coupons
// @Produce json
// @Success 200 {array} coupon.Coupon
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coupons/user [get]
// @Security ApiKeyAuth
func (h *couponHandler) FindCouponsByUserID(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		log.Warn().Err(err).Msg("User ID not found in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	coupons, err := h.couponFeature.FindCouponsByUserID(c.Request.Context(), userID)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error().
				Err(ex.Err).
				Str("user_id", userID).
				Msg("Failed to find coupons by user ID")
			c.JSON(ex.Code, gin.H{"error": ex.Message})
			return
		}

		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("Unexpected error while finding coupons by user ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find coupons"})
		return
	}

	log.Info().
		Str("user_id", userID).
		Int("count", len(coupons)).
		Msg("Coupons retrieved successfully")
	c.JSON(http.StatusOK, coupons)
}

// FindCouponsByCouponPolicyCode godoc
// @Summary Get coupons by policy code
// @Description Retrieves all coupons issued for a specific coupon policy code
// @Tags Coupons
// @Produce json
// @Param policyCode path string true "Coupon Policy Code"
// @Success 200 {array} coupon.Coupon
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /coupons/policy/{policyCode} [get]
// @Security ApiKeyAuth
func (h *couponHandler) FindCouponsByCouponPolicyCode(c *gin.Context) {
	log := zerolog.Ctx(c.Request.Context())

	couponPolicyCode := c.Param("policyCode")
	coupons, err := h.couponFeature.FindCouponsByCouponPolicyCode(c.Request.Context(), couponPolicyCode)
	if err != nil {
		if ex, ok := err.(*exception.Http); ok {
			log.Error().
				Err(ex.Err).
				Str("coupon_policy_code", couponPolicyCode).
				Msg("Failed to find coupons by policy code")
			c.JSON(ex.Code, gin.H{"error": ex.Message})
			return
		}

		log.Error().
			Err(err).
			Str("coupon_policy_code", couponPolicyCode).
			Msg("Unexpected error while finding coupons by policy code")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find coupons"})
		return
	}

	log.Info().
		Str("coupon_policy_code", couponPolicyCode).
		Int("count", len(coupons)).
		Msg("Coupons retrieved successfully by policy code")
	c.JSON(http.StatusOK, coupons)
}
