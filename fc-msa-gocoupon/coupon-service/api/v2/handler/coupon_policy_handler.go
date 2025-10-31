package handler

import (
	"net/http"
	"time"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/instrument"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type couponPolicyHandler struct {
	db     *gorm.DB
	log    zerolog.Logger
	tracer trace.Tracer
}

func NewCouponPolicyHandler(db *gorm.DB, log zerolog.Logger, tracer trace.Tracer) *couponPolicyHandler {
	return &couponPolicyHandler{
		db:     db,
		log:    log,
		tracer: tracer,
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
	ctx, span := h.tracer.Start(ctx, "handler.SearchCouponPolicy")
	defer span.End()

	log := instrument.GetLogger(ctx, h.log)

	id := c.Query("id")
	code := c.Query("code")
	name := c.Query("name")

	span.SetAttributes(
		attribute.String("query.id", id),
		attribute.String("query.code", code),
		attribute.String("query.name", name),
	)

	var policies []coupon.CouponPolicy
	query := h.db.WithContext(ctx).Model(&coupon.CouponPolicy{})

	if id == "" && code == "" && name == "" {
		_, dbSpan := h.tracer.Start(ctx, "db.FindAllCouponPolicies")
		if err := query.Find(&policies).Error; err != nil {
			dbSpan.RecordError(err)
			dbSpan.End()
			log.Error().Err(err).Msg("Failed to fetch all coupon policies")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		dbSpan.End()

		span.SetAttributes(attribute.Int("result.count", len(policies)))
		log.Info().Int("count", len(policies)).Msg("Fetched all coupon policies successfully")
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

	_, dbSpan := h.tracer.Start(ctx, "db.SearchCouponPolicies")
	if err := query.Find(&policies).Error; err != nil {
		dbSpan.RecordError(err)
		dbSpan.End()
		log.Error().Err(err).Msg("Failed to search coupon policies")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	dbSpan.End()

	span.SetAttributes(attribute.Int("result.count", len(policies)))
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
	ctx, span := h.tracer.Start(ctx, "handler.CreateCouponPolicyDummy")
	defer span.End()

	log := instrument.GetLogger(ctx, h.log)

	_, cleanupSpan := h.tracer.Start(ctx, "db.CleanupTables")
	if err := h.db.Exec("DELETE FROM coupons").Error; err != nil {
		cleanupSpan.RecordError(err)
		cleanupSpan.End()
		log.Error().Err(err).Msg("Failed to delete coupons")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupons"})
		return
	}

	if err := h.db.Exec("DELETE FROM coupon_policies").Error; err != nil {
		cleanupSpan.RecordError(err)
		cleanupSpan.End()
		log.Error().Err(err).Msg("Failed to delete coupon policies")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon policies"})
		return
	}
	cleanupSpan.End()

	policies := []coupon.CouponPolicy{
		// Valid Coupon Policy Data for `issue success case`
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

		// Not Started Coupon Policy Data for `issue invalid period case`
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-200",
			Name:                  "Future 15K Off Voucher",
			Description:           "15,000 discount for orders above 100,000 starting tomorrow",
			TotalQuantity:         50,
			StartTime:             time.Now().Add(48 * time.Hour),
			EndTime:               time.Now().Add(6 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypeFixedAmount,
			DiscountValue:         15000,
			MinimumOrderAmount:    100000,
			MaximumDiscountAmount: 15000,
		},

		// Ended Time Coupon Policy Data for `issue invalid period case`
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-300",
			Name:                  "Expired 20% Off Voucher",
			Description:           "10% discount for orders above 50,000 (expired)",
			TotalQuantity:         50,
			StartTime:             time.Now().Add(-10 * 24 * time.Hour),
			EndTime:               time.Now().Add(-3 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         10,
			MinimumOrderAmount:    50000,
			MaximumDiscountAmount: 20000,
		},

		// Empty Coupon Policy Data for `issue quota exceeded case`
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-400",
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

		// Valid Coupon Policy Data for `issue success load testing`
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-500",
			Name:                  "20% Off Voucher",
			Description:           "20% discount up to 50,000 for orders above 200,000",
			TotalQuantity:         500000,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(10 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         20,
			MinimumOrderAmount:    200000,
			MaximumDiscountAmount: 50000,
		},

		// Valid Coupon Policy Data for `issue success stress testing (~20–30% ekstra user 1.000*1.3)`
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-600",
			Name:                  "30% Off Voucher",
			Description:           "30% discount up to 50,000 for orders above 300,000",
			TotalQuantity:         13000,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(10 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         30,
			MinimumOrderAmount:    300000,
			MaximumDiscountAmount: 50000,
		},

		// Valid Coupon Policy Data for `issue success spike testing (1.5× total request)`
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-700",
			Name:                  "40% Off Voucher",
			Description:           "40% discount up to 50,000 for orders above 400,000",
			TotalQuantity:         4500,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(10 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         40,
			MinimumOrderAmount:    400000,
			MaximumDiscountAmount: 50000,
		},

		// Valid Coupon Policy Data for `issue success spike testing (total user × total request)`
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-800",
			Name:                  "50% Off Voucher",
			Description:           "50% discount up to 50,000 for orders above 400,000",
			TotalQuantity:         2160000,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(10 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         40,
			MinimumOrderAmount:    400000,
			MaximumDiscountAmount: 50000,
		},

		// Percentage Coupon Policy Case for concurrency test 100 request
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-900",
			Name:                  "15% Off Voucher",
			Description:           "Get 15% off for orders above 100,000, max discount 25,000",
			TotalQuantity:         10,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(7 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         15,
			MinimumOrderAmount:    100000,
			MaximumDiscountAmount: 25000,
		},

		// Percentage Coupon Policy Case for concurrency test 1.000 request
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-1000",
			Name:                  "20% Off Voucher",
			Description:           "Get 20% off for orders above 100,000, max discount 25,000",
			TotalQuantity:         100,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(7 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         20,
			MinimumOrderAmount:    100000,
			MaximumDiscountAmount: 25000,
		},

		// Percentage Coupon Policy Case for concurrency test 1.000.000 request
		{
			ID:                    uuid.New().String(),
			Code:                  "COUPON-1100",
			Name:                  "30% Off Voucher",
			Description:           "Get 30% off for orders above 100,000, max discount 25,000",
			TotalQuantity:         100,
			StartTime:             time.Now(),
			EndTime:               time.Now().Add(7 * 24 * time.Hour),
			DiscountType:          coupon.DiscountTypePercentage,
			DiscountValue:         30,
			MinimumOrderAmount:    100000,
			MaximumDiscountAmount: 25000,
		},
	}

	_, insertSpan := h.tracer.Start(ctx, "db.InsertDummyPolicies")
	for _, p := range policies {
		if err := h.db.Create(&p).Error; err != nil {
			insertSpan.RecordError(err)
			insertSpan.End()
			log.Error().Err(err).Msg("Failed to create dummy coupon policy")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create dummy policies"})
			return
		}
	}
	insertSpan.End()

	span.SetAttributes(attribute.Int("dummy.created_count", len(policies)))
	log.Info().Int("count", len(policies)).Msg("Dummy coupon policies created successfully")
	c.JSON(http.StatusOK, policies)
}
