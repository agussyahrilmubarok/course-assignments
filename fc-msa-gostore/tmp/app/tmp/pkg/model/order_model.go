package model

import "time"

type OrderModel struct {
	ID            string           `json:"id"`
	UserID        string           `json:"user_id"`
	Items         []OrderItemModel `json:"items"`
	TotalAmount   float64          `json:"total_amount"`
	TotalQuantity int              `json:"total_quantity"`
	Status        string           `json:"status" `
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type OrderItemModel struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type CheckoutOrderRequest struct {
	UserID string                     `json:"user_id" validate:"required,uuid4"`
	Items  []CheckoutOrderItemRequest `json:"items" validate:"required,dive,required"`
}

type CheckoutOrderItemRequest struct {
	ID       string `json:"id" validate:"required,uuid4"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
}

type OrderCancelRequest struct {
	OrderID string `json:"order_id" validate:"required,uuid4"`
}
