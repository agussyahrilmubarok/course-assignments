package order

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type CreateOrderRequest struct {
	OrderItems []OrderItemRequest `json:"order_items"`
	UserID     string             `json:"user_id"`
}

type CancelOrderRequest struct {
	OrderID string `json:"order_id" validate:"required"`
	UserID  string `json:"user_id"`
}

type PricingResponse struct {
	ProductID  string  `json:"product_id"`
	Markup     float64 `json:"markup"`
	Discount   float64 `json:"discount"`
	FinalPrice float64 `json:"final_price"`
}
