package handler

import (
	"net/http"
	"time"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/instrument"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type couponPolicyHandler struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewCouponPolicyHandler(db *gorm.DB, log zerolog.Logger) *couponPolicyHandler {
	return &couponPolicyHandler{
		db:  db,
		log: log,
	}
}

// SearchCouponPolicy godoc
// @Summary Search coupon policies (by ID, Code, Name, or all)
// @Description Get coupon policies by ID, code, or name. If no parameters are provided, return all coupon policies.
// @Tags Coupon Policy
// @Produce json
// @Param id query string false "Coupon Policy ID"
// @Param code query string false "Coupon Code"
// @Param name query string false "Coupon Name"
// @Success 200 {array} coupon.CouponPolicy
// @Failure 500 {object} map[string]string
// @Router /couponPolicies [get]
func (h *couponPolicyHandler) SearchCouponPolicy(c *gin.Context) {
	ctx := c.Request.Context()
	log := instrument.GetLogger(ctx, h.log)

	id := c.Query("id")
	code := c.Query("code")
	name := c.Query("name")

	var policies []coupon.CouponPolicy
	query := h.db.WithContext(ctx).Model(&coupon.CouponPolicy{})

	if id == "" && code == "" && name == "" {
		if err := query.Find(&policies).Error; err != nil {
			log.Error().Err(err).Msg("Failed to fetch all coupon policies")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, policies)
		return
	}

	if id != "" {
		query = query.Where("id = ?", id)
	}
	if code != "" {
		query = query.Or("code ILIKE ?", "%"+code+"%")
	}
	if name != "" {
		query = query.Or("name ILIKE ?", "%"+name+"%")
	}

	if err := query.Find(&policies).Error; err != nil {
		log.Error().Err(err).Msg("Failed to search coupon policies")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Info().Int("count", len(policies)).Msg("Search coupon policy completed")

	c.JSON(http.StatusOK, policies)
}

// CreateCouponPolicyDummy godoc
// @Summary Create dummy coupon policies
// @Tags Coupon Policy
// @Produce json
// @Success 200 {array} coupon.CouponPolicy
// @Router /couponPolicies/dummy [post]
func (h *couponPolicyHandler) CreateCouponPolicyDummy(c *gin.Context) {
	ctx := c.Request.Context()
	log := instrument.GetLogger(ctx, h.log)

	if err := h.db.Exec("DELETE FROM coupons").Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete coupons")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupons"})
		return
	}

	if err := h.db.Exec("DELETE FROM coupon_policies").Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete coupon policies")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon policies"})
		return
	}

	policies := []coupon.CouponPolicy{
		// Valid Coupon Policy Case
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-100",
			Name:                  "10K Off Voucher",
			Description:           "Discount 10,000 for orders above 50,000",
			TotalQuantity:         100,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(7 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypeFixedAmount,
			DiscountValue:         10000,
			MinimumOrderAmount:    50000,
			MaximumDiscountAmount: 10000,
		},

		// Epic Coupon Policy Case
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-200",
			Name:                  "20% Off Voucher",
			Description:           "20% discount up to 50,000 for orders above 200,000",
			TotalQuantity:         50,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(10 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         20,
			MinimumOrderAmount:    200000,
			MaximumDiscountAmount: 50000,
		},

		// Empty Coupon Policy Case
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-300",
			Name:                  "50K Off Voucher",
			Description:           "Discount 50,000 for orders above 500,000",
			TotalQuantity:         0,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(5 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypeFixedAmount,
			DiscountValue:         50000,
			MinimumOrderAmount:    500000,
			MaximumDiscountAmount: 50000,
		},

		// Not Started Coupon Policy Case
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-400",
			Name:                  "Future 15K Off Voucher",
			Description:           "15,000 discount for orders above 100,000 starting tomorrow",
			TotalQuantity:         50,
			StartTime:             time.Now().Add(24 * time.Hour),
			EndTime:               time.Now().Add(6 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypeFixedAmount,
			DiscountValue:         15000,
			MinimumOrderAmount:    100000,
			MaximumDiscountAmount: 15000,
		},

		// Ended Time Coupon Policy Case
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-500",
			Name:                  "Expired 20% Off Voucher",
			Description:           "20% discount for orders above 50,000 (expired)",
			TotalQuantity:         50,
			StartTime:             time.Now().Add(-10 * 24 * time.Hour),
			EndTime:               time.Now().Add(-3 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         20,
			MinimumOrderAmount:    50000,
			MaximumDiscountAmount: 20000,
		},

		// Percentage Coupon Policy Case
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-600",
			Name:                  "15% Off Voucher",
			Description:           "Get 15% off for orders above 100,000, max discount 25,000",
			TotalQuantity:         100,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(7 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         15,
			MinimumOrderAmount:    100000,
			MaximumDiscountAmount: 25000,
		},
	}

	for _, p := range policies {
		if err := h.db.Create(&p).Error; err != nil {
			log.Error().Err(err).Msg("Failed to create dummy coupon policy")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create dummy policies"})
			return
		}
	}

	c.JSON(http.StatusOK, policies)
}
