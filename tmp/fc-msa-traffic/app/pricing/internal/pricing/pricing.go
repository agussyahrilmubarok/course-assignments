package pricing

// Pricing represents calculated pricing info to be stored in Redis.
type Pricing struct {
	ProductID  string  `json:"product_id"`  // ID of the product
	Markup     float64 `json:"markup"`      // Markup percentage (e.g., 0.15 for 15%)
	Discount   float64 `json:"discount"`    // Discount percentage (e.g., 0.10 for 10%)
	FinalPrice float64 `json:"final_price"` // Final price = (base_price * (1 + markup)) - discount
}

type PricingRule struct {
	ProductID         string  `json:"product_id" validate:"required"`            // ID of the product
	ProductPrice      float64 `json:"product_price"`                             // Base price must be >= 0
	DefaultDiscount   float64 `json:"default_discount" validate:"gte=0,lte=1"`   // Must be between 0 and 1
	DefaultMarkup     float64 `json:"default_markup" validate:"gte=0,lte=1"`     // Must be between 0 and 1
	StockThreshold    int     `json:"stock_threshold" validate:"gte=0"`          // Must be >= 0
	DiscountReduction float64 `json:"discount_reduction" validate:"gte=0,lte=1"` // Must be between 0 and 1
	MarkupIncrease    float64 `json:"markup_increase" validate:"gte=0,lte=1"`    // Must be between 0 and 1
}
