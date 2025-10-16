package order

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type CreateOrderRequest struct {
	OrderItems []OrderItemRequest `json:"order_items"`
}

type CancelOrderRequest struct {
	OrderID string `json:"order_id" validate:"required"`
}
