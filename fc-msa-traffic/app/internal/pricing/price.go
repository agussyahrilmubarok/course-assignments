package pricing

import "time"

type Pricing struct {
	ID         string    `json:"id" gorm:"type:char(36);primaryKey"`
	ProductID  string    `json:"product_id" gorm:"type:char(36);not null;index"`
	Markup     float64   `json:"markup" gorm:"type:decimal(5,2);not null;default:0"`
	Discount   float64   `json:"discount" gorm:"type:decimal(5,2);not null;default:0"`
	FinalPrice float64   `json:"final_price" gorm:"type:decimal(12,2);not null;default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PricingRule struct {
	ID                string    `json:"id" gorm:"type:char(36);primaryKey"`
	ProductID         string    `json:"product_id" gorm:"type:char(36);not null;index"`
	ProductPrice      float64   `json:"product_price" gorm:"type:decimal(12,2);not null;default:0"`
	DefaultMarkup     float64   `json:"default_markup" gorm:"type:decimal(5,2);not null;default:0"`
	DefaultDiscount   float64   `json:"default_discount" gorm:"type:decimal(5,2);not null;default:0"`
	StockThreshold    int       `json:"stock_threshold" gorm:"not null;default:0"`
	MarkupIncrease    float64   `json:"markup_increase" gorm:"type:decimal(5,2);not null;default:0"`
	DiscountReduction float64   `json:"discount_reduction" gorm:"type:decimal(5,2);not null;default:0"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
