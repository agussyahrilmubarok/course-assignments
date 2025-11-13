package store

import "time"

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderCreated   OrderStatus = "created"
	OrderCancelled OrderStatus = "cancelled"
	OrderPaid      OrderStatus = "paid"
	OrderProcessed OrderStatus = "processed"
)

// Order represents an order document in MongoDB.
type Order struct {
	ID            string      `json:"id" bson:"_id,omitempty"`
	UserID        string      `json:"user_id" bson:"user_id"`
	Items         []OrderItem `json:"details" bson:"details"`
	TotalAmount   float64     `json:"total_amount" bson:"total_amount"`
	TotalQuantity int         `json:"total_quantity" bson:"total_quantity"`
	Status        OrderStatus `json:"status" bson:"status"`
	CreatedAt     time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" bson:"updated_at"`
}

// OrderDetail represents a single product item in an order.
type OrderItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Name      string  `json:"name" bson:"name"`
	Price     float64 `json:"price" bson:"price"`
	Quantity  int     `json:"quantity" bson:"quantity"`
}
