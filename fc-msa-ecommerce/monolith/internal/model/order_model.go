package model

type CreateOrderRequest struct {
	AddressID uint `json:"address_id" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"` // e.g., pending, paid, shipped, cancelled
}
