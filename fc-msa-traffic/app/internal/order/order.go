package order

import "time"

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusProcessed Status = "PROCESSED"
	StatusShipped   Status = "SHIPPED"
	StatusCompleted Status = "COMPLETED"
	StatusCancelled Status = "CANCELLED"
)

type Order struct {
	ID            string      `json:"id" gorm:"type:char(36);primaryKey"`
	UserID        string      `json:"user_id" gorm:"type:char(36);not null"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity      int         `json:"quantity" gorm:"not null;default:0"`
	TotalPrice    float64     `json:"total_price" gorm:"type:decimal(12,2);not null;default:0"`
	TotalMarkup   float64     `json:"total_markup" gorm:"type:decimal(12,2);not null;default:0"`
	TotalDiscount float64     `json:"total_discount" gorm:"type:decimal(12,2);not null;default:0"`
	Status        Status      `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	IdempotentKey string      `json:"idempotent_key" gorm:"type:varchar(100);uniqueIndex"`
	CreatedAt     time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderItem struct {
	ID         string    `json:"id" gorm:"type:char(36);primaryKey"`
	OrderID    string    `json:"order_id" gorm:"type:char(36);not null;index"`
	ProductID  string    `json:"product_id" gorm:"type:char(36);not null"`
	BasePrice  float64   `json:"base_price" gorm:"type:decimal(12,2);not null"`
	Quantity   int       `json:"quantity" gorm:"not null;default:1"`
	TotalPrice float64   `json:"total_price" gorm:"type:decimal(12,2);not null;default:0"`
	Markup     float64   `json:"markup" gorm:"type:decimal(12,2);not null;default:0"`
	Discount   float64   `json:"discount" gorm:"type:decimal(12,2);not null;default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
