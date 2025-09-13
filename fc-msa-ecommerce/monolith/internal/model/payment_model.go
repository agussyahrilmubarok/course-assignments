package model

type CheckoutRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Method  string `json:"method" binding:"required"` // e.g. credit_card, transfer
}

type VerifyPaymentRequest struct {
	ReferenceID string `json:"reference_id" binding:"required"`
}
