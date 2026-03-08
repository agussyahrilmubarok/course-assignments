package pricing

type PricingRuleRequest struct {
	ProductID         string  `json:"product_id" validate:"required"`
	DefaultDiscount   float64 `json:"default_discount" validate:"gte=0,lte=1"`
	DefaultMarkup     float64 `json:"default_markup" validate:"gte=0,lte=1"`
	StockThreshold    int     `json:"stock_threshold" validate:"gte=0"`
	DiscountReduction float64 `json:"discount_reduction" validate:"gte=0,lte=1"`
	MarkupIncrease    float64 `json:"markup_increase" validate:"gte=0,lte=1"`
}

func (r *PricingRuleRequest) ToPricingRule(productPrice float64) *PricingRule {
	return &PricingRule{
		ProductID:         r.ProductID,
		ProductPrice:      productPrice,
		DefaultDiscount:   r.DefaultDiscount,
		DefaultMarkup:     r.DefaultMarkup,
		StockThreshold:    r.StockThreshold,
		DiscountReduction: r.DiscountReduction,
		MarkupIncrease:    r.MarkupIncrease,
	}
}
