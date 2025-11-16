package store

import "time"

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	StatusCompleted PaymentStatus = "completed"
	StatusFailed    PaymentStatus = "failed"
	StatusPending   PaymentStatus = "pending"
)

// Payment represents a payment document in MongoDB.
type Payment struct {
	ID         string        `json:"id" bson:"_id,omitempty"`
	RefCode    string        `json:"ref_code" bson:"ref_code"`
	UserID     string        `json:"user_id" bson:"user_id"`
	OrderID    string        `json:"order_id" bson:"order_id"`
	Amount     float64       `json:"amount" bson:"amount"`
	Status     PaymentStatus `json:"status" bson:"status"`
	InvoiceURL string        `json:"invoice_url" bson:"invoice_url"`
	PaidAt     time.Time     `json:"paid_at" bson:"paid_at"`
	ExpiredAt  time.Time     `json:"expired_at" bson:"expired_at"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" bson:"updated_at"`
}
