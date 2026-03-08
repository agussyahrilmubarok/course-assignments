package order

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusProcessed OrderStatus = "PROCESSED"
	StatusShipped   OrderStatus = "SHIPPED"
	StatusCompleted OrderStatus = "COMPLETED"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusFailed    OrderStatus = "FAILED"
)

type Order struct {
	ID          string      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID      string      `json:"user_id" gorm:"type:char(36);not null"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	FinalAmount float64     `json:"final_amount" gorm:"type:decimal(12,2);not null"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
	OrderItems  []OrderItem `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (o *Order) CalculateFinalAmount() float64 {
	total := 0.0
	for _, item := range o.OrderItems {
		total += item.TotalPrice
	}
	return total
}

type OrderItem struct {
	ID           string  `json:"id" gorm:"type:char(36);primaryKey"`
	OrderID      string  `json:"order_id" gorm:"type:char(36);not null"`
	ProductID    string  `json:"product_id" gorm:"type:char(36);not null"`
	ProductPrice float64 `json:"product_price" gorm:"type:decimal(12,2);not null"` // base_price/unit
	Quantity     int     `json:"quantity" gorm:"not null;default:1"`               // total unit
	Discount     float64 `json:"discount" gorm:"type:decimal(12,2);not null"`      // price reduction given to the buyer
	MarkUp       float64 `json:"markup" gorm:"type:decimal(12,2);not null"`        // the amount added to the cost of a product to determine its selling price
	TotalPrice   float64 `json:"total_price" gorm:"type:decimal(12,2);not null"`   // (ProductPrice - Discount + MarkUp) * Quantity

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (item *OrderItem) CalculateTotalPrice() float64 {
	unitPrice := item.ProductPrice - item.Discount + item.MarkUp
	return unitPrice * float64(item.Quantity)
}
